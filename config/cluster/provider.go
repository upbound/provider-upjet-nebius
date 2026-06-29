// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package cluster

import (
	"github.com/upbound/provider-nebius/config/cluster/compute"
	"github.com/upbound/provider-nebius/config/cluster/dns"
	"github.com/upbound/provider-nebius/config/cluster/iam"
	"github.com/upbound/provider-nebius/config/cluster/mk8s"
	"github.com/upbound/provider-nebius/config/cluster/mysterybox"
	"github.com/upbound/provider-nebius/config/cluster/storage"
	"github.com/upbound/provider-nebius/config/cluster/vpc"
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
