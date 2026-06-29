// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package mysterybox

import (
	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Configure configures the mysterybox group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("nebius_mysterybox_v1_secret", func(r *config.Resource) {
		config.MoveToStatus(r.TerraformResource, "primary_version_id") // primary_version_id is managed by SecretVersion.spec.forProvider.setPrimary
		config.MoveToStatus(r.TerraformResource, "secret_version")     // secret_version is managed by SecretVersion
	})
	p.AddResourceConfigurator("nebius_mysterybox_v1_secret_version", func(r *config.Resource) {
		// parent_id of a secret version represents the secret it belongs to.
		r.References["parent_id"] = config.Reference{
			TerraformName: "nebius_mysterybox_v1_secret",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("id",true)`,
		}
		r.TerraformResource.Schema["payload"].Elem.(*schema.Resource).
			Schema["key"].Sensitive = false
	})
}
