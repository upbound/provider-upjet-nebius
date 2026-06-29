// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package mk8s

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/upbound/provider-nebius/config/common"
)

// Configure configures the mk8s group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("nebius_mk8s_v1_cluster", func(r *config.Resource) {
		r.References["control_plane.subnet_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_subnet",                                                    //nolint: goconst // Terraform resource name
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`, //nolint: goconst // Upjet extractor name
		}
		// Note(jonasz-lasut): Following fields are not marked as "sensitive" in Terraform cli schema output.
		// We need to configure them explicitly to store in connectionDetails secret.
		config.GetSchema(r.TerraformResource, "status.control_plane.auth.cluster_ca_certificate").Sensitive = true
		config.GetSchema(r.TerraformResource, "status.control_plane.endpoints.public_endpoint").Sensitive = true
		config.GetSchema(r.TerraformResource, "status.control_plane.endpoints.private_endpoint").Sensitive = true

		// Assemble the sensitive endpoint + CA values into a standard kubeconfig
		// so the cluster can be consumed directly by provider-kubernetes.
		r.Sensitive.AdditionalConnectionDetailsFn = common.MK8SClusterBuildKubeconfig
	})

	p.AddResourceConfigurator("nebius_mk8s_v1_node_group", func(r *config.Resource) {
		// parent_id of a node group represents the Cluster it belongs to.
		r.References["parent_id"] = config.Reference{
			TerraformName: "nebius_mk8s_v1_cluster",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["template.network_interfaces.subnet_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_subnet",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["template.gpu_cluster.id"] = config.Reference{
			TerraformName: "nebius_compute_v1_gpu_cluster",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["template.filesystems.existing_filesystem.id"] = config.Reference{
			TerraformName: "nebius_compute_v1_filesystem",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["template.service_account_id"] = config.Reference{
			TerraformName: "nebius_iam_v1_service_account",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		// Do not late initialize the node count as it may conflict with autoscaling/karpenter
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"fixed_node_count"},
		}
	})
}
