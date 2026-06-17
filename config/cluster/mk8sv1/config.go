// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package mk8sv1

import (
	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		r.TerraformResource.Schema["status"].Elem.(*schema.Resource).
			Schema["control_plane"].Elem.(*schema.Resource).
			Schema["auth"].Elem.(*schema.Resource).
			Schema["cluster_ca_certificate"].Sensitive = true
		r.TerraformResource.Schema["status"].Elem.(*schema.Resource).
			Schema["control_plane"].Elem.(*schema.Resource).
			Schema["endpoints"].Elem.(*schema.Resource).
			Schema["public_endpoint"].Sensitive = true
		r.TerraformResource.Schema["status"].Elem.(*schema.Resource).
			Schema["control_plane"].Elem.(*schema.Resource).
			Schema["endpoints"].Elem.(*schema.Resource).
			Schema["private_endpoint"].Sensitive = true
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
