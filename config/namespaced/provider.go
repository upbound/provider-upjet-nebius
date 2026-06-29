// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package namespaced

import (
	"github.com/upbound/provider-nebius/config/namespaced/compute"
	"github.com/upbound/provider-nebius/config/namespaced/dns"
	"github.com/upbound/provider-nebius/config/namespaced/iam"
	"github.com/upbound/provider-nebius/config/namespaced/mk8s"
	"github.com/upbound/provider-nebius/config/namespaced/mysterybox"
	"github.com/upbound/provider-nebius/config/namespaced/storage"
	"github.com/upbound/provider-nebius/config/namespaced/vpc"
)

func init() {
	ProviderConfiguration.AddConfig(vpc.Configure)
	ProviderConfiguration.AddConfig(compute.Configure)
	ProviderConfiguration.AddConfig(mk8s.Configure)
	ProviderConfiguration.AddConfig(dns.Configure)
	ProviderConfiguration.AddConfig(mysterybox.Configure)
	ProviderConfiguration.AddConfig(iam.Configure)
	ProviderConfiguration.AddConfig(storage.Configure)
}
