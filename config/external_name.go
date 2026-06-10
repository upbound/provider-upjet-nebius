package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	// The placeholder is a valid Nebius ID (NID) of the form <type>-<routingCode><weakID>.
	// The 3-char segment after the type prefix is the routing code; e0t is the default SDK routing code.
	"nebius_vpc_v1_network":        config.FrameworkResourceWithComputedIdentifier("id", "vpcnetwork-e0t000000000000000"),
	"nebius_vpc_v1_pool":           config.FrameworkResourceWithComputedIdentifier("id", "vpcpool-e0t000000000000000"),
	"nebius_vpc_v1_subnet":         config.FrameworkResourceWithComputedIdentifier("id", "vpcsubnet-e0t000000000000000"),
	"nebius_vpc_v1_route_table":    config.FrameworkResourceWithComputedIdentifier("id", "vpcroutetable-e0t000000000000000"),
	"nebius_vpc_v1_allocation":     config.FrameworkResourceWithComputedIdentifier("id", "vpcallocation-e0t000000000000000"),
	"nebius_vpc_v1_route":          config.FrameworkResourceWithComputedIdentifier("id", "vpcroute-e0t000000000000000"),
	"nebius_vpc_v1_security_group": config.FrameworkResourceWithComputedIdentifier("id", "vpcsecuritygroup-e0t000000000000000"),
	"nebius_vpc_v1_security_rule":  config.FrameworkResourceWithComputedIdentifier("id", "vpcsecurityrule-e0t000000000000000"),
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
