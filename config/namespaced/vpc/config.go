// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package vpc

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the vpc group
func Configure(p *config.Provider) { //nolint:gocyclo
	p.AddResourceConfigurator("nebius_vpc_v1_network", func(r *config.Resource) {
		r.References["ipv4_private_pools.pools.id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_pool",                                                      //nolint: goconst // Terraform resource name
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`, //nolint: goconst // Upjet extractor name
		}
		r.References["ipv4_public_pools.pools.id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_pool",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
	})

	p.AddResourceConfigurator("nebius_vpc_v1_pool", func(r *config.Resource) {
		r.References["source_pool_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_pool",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
	})

	p.AddResourceConfigurator("nebius_vpc_v1_subnet", func(r *config.Resource) {
		r.References["network_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_network",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["route_table_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_route_table",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
	})

	p.AddResourceConfigurator("nebius_vpc_v1_route_table", func(r *config.Resource) {
		r.References["network_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_network",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
	})

	p.AddResourceConfigurator("nebius_vpc_v1_allocation", func(r *config.Resource) {
		r.References["ipv4_private.subnet_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_subnet",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["ipv4_private.pool_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_pool",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["ipv4_public.subnet_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_subnet",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["ipv4_public.pool_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_pool",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
	})

	p.AddResourceConfigurator("nebius_vpc_v1_route", func(r *config.Resource) {
		// parent_id of a route represents the RouteTable it belongs to.
		r.References["parent_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_route_table",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["next_hop.allocation.id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_allocation",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
	})

	p.AddResourceConfigurator("nebius_vpc_v1_security_group", func(r *config.Resource) {
		r.References["network_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_network",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
	})

	p.AddResourceConfigurator("nebius_vpc_v1_security_rule", func(r *config.Resource) {
		// parent_id of a security rule represents the SecurityGroup it belongs to.
		r.References["parent_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_security_group",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["ingress.source_security_group_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_security_group",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.References["egress.destination_security_group_id"] = config.Reference{
			TerraformName: "nebius_vpc_v1_security_group",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
	})
}
