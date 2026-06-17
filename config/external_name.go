package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	// The placeholder is a valid Nebius ID (NID) of the form <type>-<routingCode><weakID>.
	// The 3-char segment after the type prefix is the routing code; e0t is the default SDK routing code.
	// vpc_v1 resources can have any valid routing code prefix e.g. e0t in ComputedIdentifier independently of the project they're deployed in
	"nebius_vpc_v1_network":         config.FrameworkResourceWithComputedIdentifier("id", "vpcnetwork-e0t000000000000000"),
	"nebius_vpc_v1_pool":            config.FrameworkResourceWithComputedIdentifier("id", "vpcpool-e0t000000000000000"),
	"nebius_vpc_v1_subnet":          config.FrameworkResourceWithComputedIdentifier("id", "vpcsubnet-e0t000000000000000"),
	"nebius_vpc_v1_route_table":     config.FrameworkResourceWithComputedIdentifier("id", "vpcroutetable-e0t000000000000000"),
	"nebius_vpc_v1_allocation":      config.FrameworkResourceWithComputedIdentifier("id", "vpcallocation-e0t000000000000000"),
	"nebius_vpc_v1_route":           config.FrameworkResourceWithComputedIdentifier("id", "vpcroute-e0t000000000000000"),
	"nebius_vpc_v1_security_group":  config.FrameworkResourceWithComputedIdentifier("id", "vpcsecuritygroup-e0t000000000000000"),
	"nebius_vpc_v1_security_rule":   config.FrameworkResourceWithComputedIdentifier("id", "vpcsecurityrule-e0t000000000000000"),
	"nebius_iam_v1_service_account": config.FrameworkResourceWithComputedIdentifier("id", "serviceaccount-e0t000000000000000"),
	// compute_v1 resources need to have a valid project prefix in ComputedIdentifier e.g. e00, e01
	"nebius_compute_v1_gpu_cluster": config.FrameworkResourceWithComputedIdentifier("id", "computegpucluster-e01000000000000000"),
	"nebius_compute_v1_filesystem":  config.FrameworkResourceWithComputedIdentifier("id", "computefilesystem-e01000000000000000"),
	"nebius_compute_v1_disk":        config.FrameworkResourceWithComputedIdentifier("id", "computedisk-e01000000000000000"),
	// mk8s_v1 resources can have any valid routing code prefix e.g. e0t in ComputedIdentifier independently of the project they're deployed in
	"nebius_mk8s_v1_cluster":    config.FrameworkResourceWithComputedIdentifier("id", "mk8scluster-e0t000000000000000"),
	"nebius_mk8s_v1_node_group": config.FrameworkResourceWithComputedIdentifier("id", "mk8snodegroup-e0t000000000000000"),
	// dns_v1 resources can have any valid routing code prefix e.g. e0t in ComputedIdentifier independently of the project they're deployed in
	"nebius_dns_v1_zone":   config.FrameworkResourceWithComputedIdentifier("id", "dnszone-e0t000000000000000"),
	"nebius_dns_v1_record": config.FrameworkResourceWithComputedIdentifier("id", "dnsrecord-e0t000000000000000"),
	// mysterybox_v1 resources need to have a valid project prefix in ComputedIdentifier e.g. e00, e01
	"nebius_mysterybox_v1_secret":         config.FrameworkResourceWithComputedIdentifier("id", "mbsec-e00000000000000000"),
	"nebius_mysterybox_v1_secret_version": config.FrameworkResourceWithComputedIdentifier("id", "mbsecver-e00000000000000000"),
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		e, configured := ExternalNameConfigs[r.Name]
		if !configured {
			return
		}
		r.ExternalName = e
		r.Version = versionV1Beta1
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}
