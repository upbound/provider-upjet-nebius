// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package config

import (
	"context"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	nebiusimpl "github.com/nebius/terraform-provider-nebius/provider/impl"

	"github.com/upbound/provider-nebius/config/namespaced"
	"github.com/upbound/provider-nebius/config/templates"
	"github.com/upbound/provider-nebius/hack"
)

// GetProviderNamespaced returns the namespaced provider configuration
func GetProviderNamespaced(_ context.Context) (*ujconfig.Provider, error) {
	fwProvider := nebiusimpl.New()()

	defaultResourceOptions := []ujconfig.ResourceOption{
		GroupKindOverrides(),
		ExternalNameConfigurations(),
	}

	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("nebius.m.upbound.io"),
		ujconfig.WithIncludeList([]string{}),
		ujconfig.WithControllerTemplate(templates.ControllerTemplate),
		ujconfig.WithTerraformPluginFrameworkIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithMainTemplate(hack.MainTemplate),
		ujconfig.WithTerraformPluginFrameworkProvider(fwProvider),
		ujconfig.WithSchemaTraversers(&ujconfig.SingletonListEmbedder{}),
		ujconfig.WithDefaultResourceOptions(defaultResourceOptions...),
	)

	// add custom config functions
	for _, configure := range namespaced.ProviderConfiguration {
		configure(pc)
	}

	pc.ConfigureResources()

	registerTFConversions(pc)
	return pc, nil
}

func registerTFConversions(pc *ujconfig.Provider) {
	for name, r := range pc.Resources {
		r := r
		// nothing to do if no singleton list has been converted to
		// an embedded object
		if len(r.CRDListConversionPaths()) == 0 {
			continue
		}

		// the controller will be reconciling on the CRD API version
		// with the converted API (with embedded objects in place of
		// singleton lists), so we need the appropriate Terraform
		// converter in this case.
		r.TerraformConversions = []ujconfig.TerraformConversion{
			ujconfig.NewTFSingletonConversion(),
		}

		pc.Resources[name] = r
	}
}
