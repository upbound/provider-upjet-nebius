// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package computev1

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the compute group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("nebius_compute_v1_instance", func(r *config.Resource) {
		r.References["network_interfaces.subnet_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_subnet",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`, //nolint: goconst // Upjet extractor name
		}
		r.References["network_interfaces.security_groups.id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_security_group",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["boot_disk.existing_disk.id"] = config.Reference{
			TerraformName: "nebius_compute_v1_disk",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["secondary_disks.existing_disk.id"] = config.Reference{
			TerraformName: "nebius_compute_v1_disk",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["gpu_cluster.id"] = config.Reference{
			TerraformName: "nebius_compute_v1_gpu_cluster",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["filesystems.existing_filesystem.id"] = config.Reference{
			TerraformName: "nebius_compute_v1_filesystem",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["service_account_id"] = config.Reference{
			TerraformName: "nebius_iam_v1_service_account",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
	})
}
