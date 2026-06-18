// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package storagev1

import (
	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Configure configures the storagev1 group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("nebius_storage_v1_transfer", func(r *config.Resource) {
		// source.nebius.bucket_name and destination.nebius.bucket_name reference
		// Buckets managed by this provider, resolved by the bucket name.
		r.References["source.nebius.bucket_name"] = config.Reference{
			TerraformName: "nebius_storage_v1_bucket",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("name",true)`,
		}
		r.References["destination.nebius.bucket_name"] = config.Reference{
			TerraformName: "nebius_storage_v1_bucket",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("name",true)`,
		}

		// access_key_id is not marked as "sensitive" in the Terraform schema output.
		// It pairs with the (already sensitive) secret_access_key, so configure it
		// explicitly to be a secret reference in every access_key block.
		r.TerraformResource.Schema["source"].Elem.(*schema.Resource).
			Schema["nebius"].Elem.(*schema.Resource).
			Schema["access_key"].Elem.(*schema.Resource).
			Schema["access_key_id"].Sensitive = true
		r.TerraformResource.Schema["source"].Elem.(*schema.Resource).
			Schema["s3_compatible"].Elem.(*schema.Resource).
			Schema["access_key"].Elem.(*schema.Resource).
			Schema["access_key_id"].Sensitive = true
		r.TerraformResource.Schema["destination"].Elem.(*schema.Resource).
			Schema["nebius"].Elem.(*schema.Resource).
			Schema["access_key"].Elem.(*schema.Resource).
			Schema["access_key_id"].Sensitive = true
		r.TerraformResource.Schema["destination"].Elem.(*schema.Resource).
			Schema["s3_compatible"].Elem.(*schema.Resource).
			Schema["access_key"].Elem.(*schema.Resource).
			Schema["access_key_id"].Sensitive = true
	})
}
