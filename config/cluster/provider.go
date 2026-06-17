// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package cluster

import (
	"github.com/upbound/provider-nebius/config/cluster/mk8sv1"
	"github.com/upbound/provider-nebius/config/cluster/vpcv1"
)

func init() {
	ProviderConfiguration.AddConfig(vpcv1.Configure)
	ProviderConfiguration.AddConfig(mk8sv1.Configure)
}
