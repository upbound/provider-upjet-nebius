// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package common

import (
	"context"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/fieldpath"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/config"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ProjectParentedResources is the single source of truth for the Terraform
// resources whose parent_id is the Nebius project. For these, parentId becomes
// optional and is defaulted from ProviderConfig.spec.projectID.
//
// Only project-parented kinds belong here. Resource-parented kinds (parent is
// another managed resource) and tenant-parented kinds (parent is the tenant,
// not the project) MUST NOT be added: defaulting their parent_id to a project
// id would be wrong.
var ProjectParentedResources = []string{
	// vpc_v1
	"nebius_vpc_v1_network",
	"nebius_vpc_v1_pool",
	"nebius_vpc_v1_subnet",
	"nebius_vpc_v1_route_table",
	"nebius_vpc_v1_allocation",
	"nebius_vpc_v1_security_group",
	// iam_(v1|v2)
	"nebius_iam_v1_service_account",
	"nebius_iam_v1_auth_public_key",
	"nebius_iam_v1_federated_credentials",
	"nebius_iam_v2_access_key",
	// compute_v1
	"nebius_compute_v1_gpu_cluster",
	"nebius_compute_v1_filesystem",
	"nebius_compute_v1_disk",
	"nebius_compute_v1_instance",
	// dns_v1
	"nebius_dns_v1_zone",
	// mysterybox_v1
	"nebius_mysterybox_v1_secret",
	// mk8s_v1
	"nebius_mk8s_v1_cluster",
	// storage_v1
	"nebius_storage_v1_bucket",
	"nebius_storage_v1_transfer",
	// kms_v1
	"nebius_kms_v1_asymmetric_key",
	"nebius_kms_v1_symmetric_key",
	// registry_v1
	"nebius_registry_v1_registry",
	// quotas_v1
	"nebius_quotas_v1_quota_allowance",
}

const (
	parentIDSchemaKey           = "parent_id"
	forProviderParentIDPath     = "spec.forProvider.parentId"
	initProviderParentIDPath    = "spec.initProvider.parentId"
	providerConfigProjectIDPath = "spec.projectID"

	// ProviderConfig GroupVersionKind coordinates. These mirror the root groups
	// configured in config.GetProvider / GetProviderNamespaced and are kept as
	// literals here on purpose: importing the generated apis types would pull the
	// apis packages into the code generator's build graph, and apis/generate.go
	// deletes their DeepCopyObject before the generator compiles, breaking
	// `make generate`.
	clusterProviderConfigGroup    = "nebius.upbound.io"
	namespacedProviderConfigGroup = "nebius.m.upbound.io"
	providerConfigVersion         = "v1beta1"
	clusterProviderConfigKind     = "ProviderConfig"
)

// ConfigureProjectParent applies both levers for a project-parented resource:
// it makes parent_id optional so codegen drops the required-parameter CEL rule
// (Lever 1), and appends an initializer that defaults parentId from the
// ProviderConfig at reconcile time (Lever 2).
func ConfigureProjectParent(r *config.Resource) {
	if s, ok := r.TerraformResource.Schema[parentIDSchemaKey]; ok {
		s.Required = false
		s.Optional = true
	}
	r.InitializerFns = append(r.InitializerFns, NewParentIDInitializer)
}

type parentIDInitializer struct {
	kube client.Client
}

// NewParentIDInitializer returns an initializer that defaults
// spec.forProvider.parentId from the referenced ProviderConfig.spec.projectID.
// Its signature matches config.NewInitializerFn.
func NewParentIDInitializer(kube client.Client) managed.Initializer {
	return &parentIDInitializer{kube: kube}
}

// Initialize fills spec.forProvider.parentId from the ProviderConfig when it is
// not already set on the resource. parent_id is immutable, so pinning the
// resolved value into the spec is correct and mirrors existing external-name
// defaulting.
func (i *parentIDInitializer) Initialize(ctx context.Context, mg xpresource.Managed) error {
	// Do not mutate spec.forProvider for observe-only resources.
	if sets.New(mg.GetManagementPolicies()...).
		Equal(sets.New(xpv1.ManagementActionObserve)) {
		return nil
	}

	paved, err := fieldpath.PaveObject(mg)
	if err != nil {
		return err
	}

	forProvider, _ := paved.GetString(forProviderParentIDPath)
	if forProvider != "" {
		return nil
	}

	initProvider, _ := paved.GetString(initProviderParentIDPath)
	if initProvider != "" {
		return nil
	}

	projectID, err := i.providerConfigProjectID(ctx, mg)
	if err != nil {
		return errors.Wrap(err, "cannot default to provider config's project id")
	}
	if projectID == "" {
		return errors.Errorf("parentId is not set on %q and the referenced ProviderConfig has no spec.projectID", mg.GetName())
	}

	// Pin the resolved value onto the in-memory resource so the current
	// reconcile pass observes and creates with it, then persist it via
	// server-side apply so it survives subsequent reconciles.
	if err := paved.SetString(forProviderParentIDPath, projectID); err != nil {
		return err
	}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(paved.UnstructuredContent(), mg); err != nil {
		return errors.Wrap(err, "cannot write defaulted parentId back to the resource")
	}

	patch := &unstructured.Unstructured{}
	patch.SetGroupVersionKind(mg.GetObjectKind().GroupVersionKind())
	patch.SetName(mg.GetName())
	patch.SetNamespace(mg.GetNamespace())

	if err := unstructured.SetNestedField(patch.Object, projectID, "spec", "forProvider", "parentId"); err != nil {
		return err
	}
	return i.kube.Apply(ctx, client.ApplyConfigurationFromUnstructured(patch), client.FieldOwner("provider"), client.ForceOwnership)
}

// providerConfigProjectID reads spec.projectID from the ProviderConfig
// referenced by mg. It resolves the ProviderConfig kind/scope from the
// managed-resource type and reads it as an unstructured object so this package
// does not depend on the generated apis types.
func (i *parentIDInitializer) providerConfigProjectID(ctx context.Context, mg xpresource.Managed) (string, error) {
	pc := &unstructured.Unstructured{}
	var nn types.NamespacedName

	switch m := mg.(type) {
	case xpresource.ModernManaged:
		ref := m.GetProviderConfigReference()
		if ref == nil {
			return "", errors.New("no providerConfigRef provided")
		}
		// Namespaced provider: the ProviderConfig lives in the namespaced group;
		// ref.Kind distinguishes ProviderConfig (namespaced) from
		// ClusterProviderConfig (cluster-scoped). The namespace is ignored by the
		// client for the cluster-scoped kind.
		pc.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   namespacedProviderConfigGroup,
			Version: providerConfigVersion,
			Kind:    ref.Kind,
		})
		nn = types.NamespacedName{Name: ref.Name, Namespace: m.GetNamespace()}
	case xpresource.LegacyManaged: //nolint:staticcheck // still handling cluster-scoped behavior
		ref := m.GetProviderConfigReference()
		if ref == nil {
			return "", errors.New("no providerConfigRef provided")
		}
		// Cluster-scoped (legacy) provider: the ProviderConfig is the
		// cluster-scoped kind in the cluster group.
		pc.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   clusterProviderConfigGroup,
			Version: providerConfigVersion,
			Kind:    clusterProviderConfigKind,
		})
		nn = types.NamespacedName{Name: ref.Name}
	default:
		return "", errors.New("resource is not a managed resource")
	}

	if err := i.kube.Get(ctx, nn, pc); err != nil {
		return "", errors.Wrap(err, "cannot get referenced ProviderConfig")
	}

	return fieldpath.Pave(pc.Object).GetString(providerConfigProjectIDPath)
}
