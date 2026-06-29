// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package dns

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the dns group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("nebius_dns_v1_zone", func(r *config.Resource) {
		r.References["vpc.primary_network_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_network",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
	})

	p.AddResourceConfigurator("nebius_dns_v1_record", func(r *config.Resource) {
		// parent_id of a record represents the DNS zone it belongs to.
		r.References["parent_id"] = config.Reference{
			TerraformName: "nebius_dns_v1_zone",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
	})
}
