// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package namespaced

import (
	"github.com/upbound/provider-nebius/config/namespaced/vpcv1"
)

func init() {
	ProviderConfiguration.AddConfig(vpcv1.Configure)
}
