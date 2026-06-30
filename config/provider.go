// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	"context"
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	nebiusimpl "github.com/nebius/terraform-provider-nebius/provider/impl"

	"github.com/upbound/provider-nebius/config/cluster"
	"github.com/upbound/provider-nebius/config/common"
	"github.com/upbound/provider-nebius/config/templates"
)

const (
	resourcePrefix = "nebius"
	modulePath     = "github.com/upbound/provider-nebius"
	versionV1Beta1 = "v1beta1"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider(_ context.Context) (*ujconfig.Provider, error) {
	fwProvider := nebiusimpl.New()()

	defaultResourceOptions := []ujconfig.ResourceOption{
		GroupKindOverrides(),
		ExternalNameConfigurations(),
	}

	pc := ujconfig.NewProvider(
		[]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("nebius.upbound.io"),
		ujconfig.WithIncludeList([]string{}),
		ujconfig.WithControllerTemplate(templates.ControllerTemplate),
		ujconfig.WithTerraformPluginFrameworkIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithTerraformPluginFrameworkProvider(fwProvider),
		ujconfig.WithSchemaTraversers(&ujconfig.SingletonListEmbedder{}),
		ujconfig.WithDefaultResourceOptions(defaultResourceOptions...),
	)

	// add custom config functions
	for _, configure := range cluster.ProviderConfiguration {
		configure(pc)
	}

	// Default parent_id from ProviderConfig.spec.projectID for project-parented
	// resources. Must run before ConfigureResources so the configurators apply.
	for _, name := range common.ProjectParentedResources {
		pc.AddResourceConfigurator(name, common.ConfigureProjectParent)
	}

	pc.ConfigureResources()

	registerTFConversions(pc)

	return pc, nil
}
