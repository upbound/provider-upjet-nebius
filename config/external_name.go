package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	// The placeholder is a valid Nebius ID (NID) of the form <type>-<routingCode><weakID>.
	// The 3-char segment after the type prefix is the routing code, which names a region
	// (e.g. e00 eu-north1, e01 eu-west1, u00 us-central1). The SDK's legacy default e0t matches
	// no region, and some instances reject it with "This instance expects u00, but received e0t",
	// so placeholders use a real region code.
	// vpc_v1 resources can have any valid routing code prefix e.g. u00 in ComputedIdentifier independently of the project they're deployed in
	"nebius_vpc_v1_network":        config.FrameworkResourceWithComputedIdentifier("id", "vpcnetwork-u00000000000000000"),
	"nebius_vpc_v1_pool":           config.FrameworkResourceWithComputedIdentifier("id", "vpcpool-u00000000000000000"),
	"nebius_vpc_v1_subnet":         config.FrameworkResourceWithComputedIdentifier("id", "vpcsubnet-u00000000000000000"),
	"nebius_vpc_v1_route_table":    config.FrameworkResourceWithComputedIdentifier("id", "vpcroutetable-u00000000000000000"),
	"nebius_vpc_v1_allocation":     config.FrameworkResourceWithComputedIdentifier("id", "vpcallocation-u00000000000000000"),
	"nebius_vpc_v1_route":          config.FrameworkResourceWithComputedIdentifier("id", "vpcroute-u00000000000000000"),
	"nebius_vpc_v1_security_group": config.FrameworkResourceWithComputedIdentifier("id", "vpcsecuritygroup-u00000000000000000"),
	"nebius_vpc_v1_security_rule":  config.FrameworkResourceWithComputedIdentifier("id", "vpcsecurityrule-u00000000000000000"),
	// iam_v(1|2) resources can have any valid routing code prefix e.g. u00 in ComputedIdentifier independently of the project they're deployed in
	"nebius_iam_v1_service_account":        config.FrameworkResourceWithComputedIdentifier("id", "serviceaccount-u00000000000000000"),
	"nebius_iam_v1_group":                  config.FrameworkResourceWithComputedIdentifier("id", "group-u00000000000000000"),
	"nebius_iam_v1_group_membership":       config.FrameworkResourceWithComputedIdentifier("id", "groupmembership-u00000000000000000"),
	"nebius_iam_v2_access_key":             config.FrameworkResourceWithComputedIdentifier("id", "accesskey-u00000000000000000"),
	"nebius_iam_v1_access_permit":          config.FrameworkResourceWithComputedIdentifier("id", "accesspermit-u00000000000000000"),
	"nebius_iam_v1_auth_public_key":        config.FrameworkResourceWithComputedIdentifier("id", "publickey-u00000000000000000"),
	"nebius_iam_v1_federated_credentials":  config.FrameworkResourceWithComputedIdentifier("id", "federatedcredentials-u00000000000000000"),
	"nebius_iam_v1_federation":             config.FrameworkResourceWithComputedIdentifier("id", "federation-u00000000000000000"),
	"nebius_iam_v1_federation_certificate": config.FrameworkResourceWithComputedIdentifier("id", "federationcertificate-u00000000000000000"),
	"nebius_iam_v1_invitation":             config.FrameworkResourceWithComputedIdentifier("id", "invitation-u00000000000000000"),
	"nebius_iam_v2_project":                config.FrameworkResourceWithComputedIdentifier("id", "project-u00000000000000000"),
	// compute_v1 resources reject routing codes that are not a real region, e.g. e0t: "region of disk id ... is not supported"
	"nebius_compute_v1_gpu_cluster": config.FrameworkResourceWithComputedIdentifier("id", "computegpucluster-u00000000000000000"),
	"nebius_compute_v1_filesystem":  config.FrameworkResourceWithComputedIdentifier("id", "computefilesystem-u00000000000000000"),
	"nebius_compute_v1_disk":        config.FrameworkResourceWithComputedIdentifier("id", "computedisk-u00000000000000000"),
	// computeinstance is grep-confirmed in the gosdk.
	"nebius_compute_v1_instance": config.FrameworkResourceWithComputedIdentifier("id", "computeinstance-u00000000000000000"),
	// computenvlinstancegroup follows the compute<message> convention (gosdk message NVLInstanceGroup); confirm via E2E observe.
	"nebius_compute_v1_nvl_instance_group": config.FrameworkResourceWithComputedIdentifier("id", "computenvlinstancegroup-u00000000000000000"),
	// mk8s_v1 resources can have any valid routing code prefix e.g. u00 in ComputedIdentifier independently of the project they're deployed in
	"nebius_mk8s_v1_cluster":    config.FrameworkResourceWithComputedIdentifier("id", "mk8scluster-u00000000000000000"),
	"nebius_mk8s_v1_node_group": config.FrameworkResourceWithComputedIdentifier("id", "mk8snodegroup-u00000000000000000"),
	// dns_v1 resources can have any valid routing code prefix e.g. u00 in ComputedIdentifier independently of the project they're deployed in
	"nebius_dns_v1_zone":   config.FrameworkResourceWithComputedIdentifier("id", "dnszone-u00000000000000000"),
	"nebius_dns_v1_record": config.FrameworkResourceWithComputedIdentifier("id", "dnsrecord-u00000000000000000"),
	// mysterybox_v1 resources currently accept any routing code in ComputedIdentifier
	"nebius_mysterybox_v1_secret":         config.FrameworkResourceWithComputedIdentifier("id", "mbsec-u00000000000000000"),
	"nebius_mysterybox_v1_secret_version": config.FrameworkResourceWithComputedIdentifier("id", "mbsecver-u00000000000000000"),
	// storage_v1 resources reject routing codes that are not a real region, e.g. e0t
	"nebius_storage_v1_bucket": config.FrameworkResourceWithComputedIdentifier("id", "storagebucket-u00000000000000000"),
	// Transfer IDs use the "u00" routing code followed by a UUID weak ID, e.g.
	// storagetransfer-u00ee51697f-d12d-4831-8e4a-6c793142da94.
	"nebius_storage_v1_transfer": config.FrameworkResourceWithComputedIdentifier("id", "storagetransfer-u0000000000-0000-0000-0000-000000000000"),
	// quotas_v1 / registry_v1 / kms_v1 type prefixes and the u00 routing code are taken from
	// real resource IDs observed in the test project (e.g. registry-u00fjr764yng9b6w19,
	// kmssymkey-u00wr5je7pmjae87ep, kmsasymkey-u00yqjhwefn8vwye1r) via the Nebius CLI.
	"nebius_quotas_v1_quota_allowance": config.FrameworkResourceWithComputedIdentifier("id", "quotaallowance-u00000000000000000"),
	"nebius_registry_v1_registry":      config.FrameworkResourceWithComputedIdentifier("id", "registry-u00000000000000000"),
	"nebius_kms_v1_asymmetric_key":     config.FrameworkResourceWithComputedIdentifier("id", "kmsasymkey-u00000000000000000"),
	"nebius_kms_v1_symmetric_key":      config.FrameworkResourceWithComputedIdentifier("id", "kmssymkey-u00000000000000000"),
	// tunnel_v1: the SDK type is "applicationtunnel" (not "tunnel") with routing "u00", taken from
	// a real ID the API assigned in E2E (applicationtunnel-u00v6sb9b47bktqbhw); tunnel is not in the
	// Nebius CLI, so this was confirmed via the resource's external-name after create, not the CLI.
	"nebius_tunnel_v1_tunnel":               config.FrameworkResourceWithComputedIdentifier("id", "applicationtunnel-u00000000000000000"),
	"nebius_capacity_v1_capacity_allowance": config.FrameworkResourceWithComputedIdentifier("id", "capacityallowance-u00000000000000000"),
	// compute_v1_disk_snapshot follows the compute_v1 routing code rules above
	"nebius_compute_v1_disk_snapshot": config.FrameworkResourceWithComputedIdentifier("id", "computedisksnapshot-u00000000000000000"),
	// Like storage_v1_transfer, inventory IDs use a UUID weak ID, e.g.
	// storagebucketinventory-u00a5aaf6f6-4836-4177-9ada-8375f669668e
	"nebius_storage_v1_inventory": config.FrameworkResourceWithComputedIdentifier("id", "storagebucketinventory-u0000000000-0000-0000-0000-000000000000"),
	// billing_v1: the SDK type is "pricingpolicy" (not "billingpricingpolicy"), taken from the
	// NID annotation on GetPricingPolicyRequest.id in gosdk.
	"nebius_billing_v1_pricing_policy": config.FrameworkResourceWithComputedIdentifier("id", "pricingpolicy-u00000000000000000"),
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
