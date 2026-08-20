// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kingpin/v2"
	authv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	xpcontroller "github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/feature"
	"github.com/crossplane/crossplane-runtime/v2/pkg/gate"
	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/crossplane/crossplane-runtime/v2/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/customresourcesgate"
	tjcontroller "github.com/crossplane/upjet/v2/pkg/controller"

	clusterapis "github.com/upbound/provider-nebius/apis/cluster"
	namespacedapis "github.com/upbound/provider-nebius/apis/namespaced"
	"github.com/upbound/provider-nebius/config"
	"github.com/upbound/provider-nebius/internal/clients"
	clustercontroller "github.com/upbound/provider-nebius/internal/controller/cluster"
	namespacedcontroller "github.com/upbound/provider-nebius/internal/controller/namespaced"
	"github.com/upbound/provider-nebius/internal/features"
)

const (
	webhookTLSCertDirEnvVar = "WEBHOOK_TLS_CERT_DIR"
	tlsServerCertDirEnvVar  = "TLS_SERVER_CERTS_DIR"
	certsDirEnvVar          = "CERTS_DIR"
	tlsServerCertDir        = "/tls/server"
)

//nolint:gocyclo // main wires up flags, scheme, manager, and controllers in one place; the crossplane-runtime v2.4.0 secret-cache flag pushed this over the default threshold.
func main() {
	var (
		app              = kingpin.New(filepath.Base(os.Args[0]), "Terraform based Crossplane provider for Nebius").DefaultEnvars()
		debug            = app.Flag("debug", "Run with debug logging.").Short('d').Bool()
		syncInterval     = app.Flag("sync", "Sync interval controls how often all resources will be double checked for drift.").Short('s').Default("1h").Duration()
		pollInterval     = app.Flag("poll", "Poll interval controls how often an individual resource should be checked for drift.").Default("10m").Duration()
		leaderElection   = app.Flag("leader-election", "Use leader election for the controller manager.").Short('l').Default("false").OverrideDefaultFromEnvar("LEADER_ELECTION").Bool()
		maxReconcileRate = app.Flag("max-reconcile-rate", "The global maximum rate per second at which resources may be checked for drift from the desired state.").Default("10").Int()

		webhookPort            = app.Flag("webhook-port", "The port the webhook listens on").Default("9443").Envar("WEBHOOK_PORT").Int()
		metricsBindAddress     = app.Flag("metrics-bind-address", "The address the metrics server listens on").Default(":8080").Envar("METRICS_BIND_ADDRESS").String()
		healthProbeBindAddress = app.Flag("health-probe-bind-addr", "The address the health/readiness probe server listens on").Default(":8081").Envar("HEALTH_PROBE_BIND_ADDRESS").String()

		enableManagementPolicies = app.Flag("enable-management-policies", "Enable support for Management Policies.").Default("true").Envar("ENABLE_MANAGEMENT_POLICIES").Bool()
		enableSecretCache        = app.Flag("enable-secret-cache", "Enable caching of Secret objects. When true, Secrets are served from the informer cache instead of direct API calls. This reduces API server load but increases memory usage.").Default("true").Envar("ENABLE_SECRET_CACHE").Bool()

		certsDirSet = false
		// we record whether the command-line option "--certs-dir" was supplied
		// in the registered PreAction for the flag.
		certsDir = app.Flag("certs-dir", "The directory that contains the server key and certificate.").Default(tlsServerCertDir).Envar(certsDirEnvVar).PreAction(func(_ *kingpin.ParseContext) error {
			certsDirSet = true
			return nil
		}).String()
	)

	kingpin.MustParse(app.Parse(os.Args[1:]))
	log.Default().SetOutput(io.Discard)
	ctrl.SetLogger(zap.New(zap.WriteTo(io.Discard)))

	zl := zap.New(zap.UseDevMode(*debug))
	log := logging.NewLogrLogger(zl.WithName("provider-nebius"))
	if *debug {
		// The controller-runtime runs with a no-op logger by default. It is
		// *very* verbose even at info level, so we only provide it a real
		// logger when we're running in debug mode.
		ctrl.SetLogger(zl)
	}

	// currently, we configure the jitter to be the 5% of the poll interval
	pollJitter := time.Duration(float64(*pollInterval) * 0.05)
	log.Debug("Starting", "sync-interval", syncInterval.String(),
		"poll-interval", pollInterval.String(), "poll-jitter", pollJitter, "max-reconcile-rate", *maxReconcileRate)

	cfg, err := ctrl.GetConfig()
	kingpin.FatalIfError(err, "Cannot get API server rest config")

	// Get the TLS certs directory from the environment variables set by
	// Crossplane if they're available. In older XP versions we used
	// WEBHOOK_TLS_CERT_DIR, in newer versions we use TLS_SERVER_CERTS_DIR. If an
	// explicit certs dir is not supplied via the command-line options, then these
	// environment variables are used instead.
	if !certsDirSet {
		xpCertsDir := os.Getenv(certsDirEnvVar)
		if xpCertsDir == "" {
			xpCertsDir = os.Getenv(tlsServerCertDirEnvVar)
		}
		if xpCertsDir == "" {
			xpCertsDir = os.Getenv(webhookTLSCertDirEnvVar)
		}
		if xpCertsDir != "" {
			*certsDir = xpCertsDir
		}
	}

	var clientOpts client.Options
	if !*enableSecretCache {
		clientOpts = client.Options{
			Cache: &client.CacheOptions{
				DisableFor: []client.Object{&corev1.Secret{}},
			},
		}
	}

	scheme := runtime.NewScheme()
	kingpin.FatalIfError(clientgoscheme.AddToScheme(scheme), "Cannot add client-go APIs to scheme")
	kingpin.FatalIfError(clusterapis.AddToScheme(scheme), "Cannot add cluster-scoped Nebius APIs to scheme")
	kingpin.FatalIfError(namespacedapis.AddToScheme(scheme), "Cannot add namespaced Nebius APIs to scheme")
	kingpin.FatalIfError(apiextensionsv1.AddToScheme(scheme), "Cannot register k8s apiextensions APIs to scheme")

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:           scheme,
		LeaderElection:   *leaderElection,
		LeaderElectionID: "crossplane-leader-election-provider-nebius",
		Cache: cache.Options{
			SyncPeriod: syncInterval,
			ByObject: map[client.Object]cache.ByObject{
				&apiextensionsv1.CustomResourceDefinition{}: {
					Transform: customresourcesgate.TransformStripCRDSchema,
				},
			},
		},
		Client: clientOpts,
		Metrics: metricsserver.Options{
			BindAddress: *metricsBindAddress,
		},
		WebhookServer: webhook.NewServer(
			webhook.Options{
				CertDir: *certsDir,
				Port:    *webhookPort,
			}),
		HealthProbeBindAddress:     *healthProbeBindAddress,
		LeaderElectionResourceLock: resourcelock.LeasesResourceLock,
		LeaseDuration:              func() *time.Duration { d := 60 * time.Second; return &d }(),
		RenewDeadline:              func() *time.Duration { d := 50 * time.Second; return &d }(),
	})
	kingpin.FatalIfError(err, "Cannot create controller manager")
	if len(*certsDir) > 0 {
		kingpin.FatalIfError(mgr.AddReadyzCheck("webhook", mgr.GetWebhookServer().StartedChecker()), "Cannot add webhook server readyz checker to controller manager")
	}

	ctx := context.Background()
	clusterProvider, err := config.GetProvider(context.Background())
	kingpin.FatalIfError(err, "Cannot initialize the cluster-scoped provider configuration")
	namespacedProvider, err := config.GetProviderNamespaced(context.Background())
	kingpin.FatalIfError(err, "Cannot initialize the namespaced provider configuration")

	clusterOpts := tjcontroller.Options{
		Options: xpcontroller.Options{
			Logger:                  log,
			GlobalRateLimiter:       ratelimiter.NewGlobal(*maxReconcileRate),
			PollInterval:            *pollInterval,
			MaxConcurrentReconciles: *maxReconcileRate,
			Features:                &feature.Flags{},
		},
		Provider:              clusterProvider,
		SetupFn:               clients.TerraformSetupBuilder(clusterProvider.TerraformPluginFrameworkProvider),
		PollJitter:            pollJitter,
		OperationTrackerStore: tjcontroller.NewOperationStore(log),
	}

	namespacedOpts := tjcontroller.Options{
		Options: xpcontroller.Options{
			Logger:                  log,
			GlobalRateLimiter:       ratelimiter.NewGlobal(*maxReconcileRate),
			PollInterval:            *pollInterval,
			MaxConcurrentReconciles: *maxReconcileRate,
			Features:                &feature.Flags{},
		},
		Provider:              namespacedProvider,
		SetupFn:               clients.TerraformSetupBuilder(namespacedProvider.TerraformPluginFrameworkProvider),
		PollJitter:            pollJitter,
		OperationTrackerStore: tjcontroller.NewOperationStore(log),
	}

	if *enableManagementPolicies {
		clusterOpts.Features.Enable(features.EnableBetaManagementPolicies)
		namespacedOpts.Features.Enable(features.EnableBetaManagementPolicies)
		log.Info("Beta feature enabled", "flag", features.EnableBetaManagementPolicies)
	}

	canSafeStart, err := canWatchCRD(ctx, mgr)
	kingpin.FatalIfError(err, "SafeStart precheck failed")
	if canSafeStart {
		crdGate := new(gate.Gate[schema.GroupVersionKind])
		clusterOpts.Gate = crdGate
		namespacedOpts.Gate = crdGate
		kingpin.FatalIfError(customresourcesgate.Setup(mgr, namespacedOpts.Options), "Cannot setup CRD gate")
		kingpin.FatalIfError(clustercontroller.SetupGated(mgr, clusterOpts), "Cannot setup cluster-scoped Nebius controllers")
		kingpin.FatalIfError(namespacedcontroller.SetupGated(mgr, namespacedOpts), "Cannot setup namespaced Nebius controllers")
	} else {
		log.Info("Provider has missing RBAC permissions for watching CRDs, controller SafeStart capability will be disabled")
		kingpin.FatalIfError(clustercontroller.Setup(mgr, clusterOpts), "Cannot setup cluster-scoped Nebius controllers")
		kingpin.FatalIfError(namespacedcontroller.Setup(mgr, namespacedOpts), "Cannot setup namespaced Nebius controllers")
	}

	// conversion webhooks are registered eagerly on every pod, not gated
	// behind leader election, whenever a certs dir is configured.
	if len(*certsDir) > 0 {
		kingpin.FatalIfError(clustercontroller.SetupWebhookWithManager(mgr), "Cannot setup cluster-scoped Nebius webhooks")
		kingpin.FatalIfError(namespacedcontroller.SetupWebhookWithManager(mgr), "Cannot setup namespaced Nebius webhooks")
	}
	kingpin.FatalIfError(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}

func canWatchCRD(ctx context.Context, mgr manager.Manager) (bool, error) {
	if err := authv1.AddToScheme(mgr.GetScheme()); err != nil {
		return false, err
	}
	verbs := []string{"get", "list", "watch"}
	for _, verb := range verbs {
		sar := &authv1.SelfSubjectAccessReview{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Group:    "apiextensions.k8s.io",
					Resource: "customresourcedefinitions",
					Verb:     verb,
				},
			},
		}
		if err := mgr.GetClient().Create(ctx, sar); err != nil {
			return false, errors.Wrapf(err, "unable to perform RBAC check for verb %s on CustomResourceDefinitions", verbs)
		}
		if !sar.Status.Allowed {
			return false, nil
		}
	}
	return true, nil
}
