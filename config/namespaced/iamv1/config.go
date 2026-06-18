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

	p.AddResourceConfigurator("nebius_iam_v1_access_permit", func(r *config.Resource) {
		// parent_id of an access permit represents the group the role is granted to.
		r.References["parent_id"] = config.Reference{
			TerraformName: "nebius_iam_v1_group",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
	})

	p.AddResourceConfigurator("nebius_iam_v1_auth_public_key", func(r *config.Resource) {
		// account.service_account.id references the service account the key authenticates.
		r.References["account.service_account.id"] = config.Reference{
			TerraformName: "nebius_iam_v1_service_account",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		// data carries a public key; the Terraform schema does not mark it
		// sensitive, so configure it as a secret reference explicitly.
		r.TerraformResource.Schema["data"].Sensitive = true
	})

	p.AddResourceConfigurator("nebius_iam_v1_federated_credentials", func(r *config.Resource) {
		// subject_id is the IAM service account the federated subject impersonates.
		r.References["subject_id"] = config.Reference{
			TerraformName: "nebius_iam_v1_service_account",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
	})

	p.AddResourceConfigurator("nebius_iam_v1_federation_certificate", func(r *config.Resource) {
		// parent_id of a federation certificate represents the Federation it belongs to.
		r.References["parent_id"] = config.Reference{
			TerraformName: "nebius_iam_v1_federation",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		// data carries an X.509 certificate; the Terraform schema does not mark it
		// sensitive, so configure it as a secret reference explicitly.
		r.TerraformResource.Schema["data"].Sensitive = true
	})
}
