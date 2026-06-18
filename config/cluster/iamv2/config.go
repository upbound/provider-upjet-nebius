// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package iamv2

import (
	"github.com/crossplane/crossplane-runtime/v2/pkg/fieldpath"
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the iamv2 group
func Configure(p *config.Provider) {
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
