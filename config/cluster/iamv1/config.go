// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package iamv1

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the iamv1 group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("nebius_iam_v1_group_membership", func(r *config.Resource) {
		// parent_id of a group membership represents the Group it belongs to.
		r.References["parent_id"] = config.Reference{
			TerraformName: "nebius_iam_v1_group",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		// member_id can be a service account (managed here) or a tenant user
		// account; the reference lets a service account be selected, a literal
		// user account id still works.
		r.References["member_id"] = config.Reference{
			TerraformName: "nebius_iam_v1_service_account",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
	})
}
