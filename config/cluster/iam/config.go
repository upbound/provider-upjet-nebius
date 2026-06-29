// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package iam

import (
	"github.com/crossplane/crossplane-runtime/v2/pkg/fieldpath"
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the iam group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("nebius_iam_v1_group_membership", func(r *config.Resource) {
		// parent_id of a group membership represents the Group it belongs to.
		r.References["parent_id"] = config.Reference{
			TerraformName: "nebius_iam_v1_group",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`, //nolint: goconst // Upjet extractor name
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

	p.AddResourceConfigurator("nebius_iam_v2_access_key", func(r *config.Resource) {
		// account.service_account.id references the service account that owns the key.
		r.References["account.service_account.id"] = config.Reference{
			TerraformName: "nebius_iam_v1_service_account",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}

		// The secret access key already flows to connection details as "secret"
		// (it is a sensitive attribute). aws_access_key_id is not sensitive, so it
		// must be added explicitly to be usable alongside the secret.
		r.Sensitive.AdditionalConnectionDetailsFn = func(attr map[string]any) (map[string][]byte, error) {
			awsAccessKeyID, err := fieldpath.Pave(attr).GetString("status.aws_access_key_id")
			if err != nil {
				return nil, err
			}

			return map[string][]byte{
				"access_key_id": []byte(awsAccessKeyID),
			}, nil
		}
	})
}
