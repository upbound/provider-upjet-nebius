// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	providerconfig "github.com/upbound/provider-nebius/internal/controller/cluster/providerconfig"
	allocation "github.com/upbound/provider-nebius/internal/controller/cluster/vpcv1/allocation"
	network "github.com/upbound/provider-nebius/internal/controller/cluster/vpcv1/network"
	pool "github.com/upbound/provider-nebius/internal/controller/cluster/vpcv1/pool"
	route "github.com/upbound/provider-nebius/internal/controller/cluster/vpcv1/route"
	routetable "github.com/upbound/provider-nebius/internal/controller/cluster/vpcv1/routetable"
	securitygroup "github.com/upbound/provider-nebius/internal/controller/cluster/vpcv1/securitygroup"
	securityrule "github.com/upbound/provider-nebius/internal/controller/cluster/vpcv1/securityrule"
	subnet "github.com/upbound/provider-nebius/internal/controller/cluster/vpcv1/subnet"
)

// Setup_monolith creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_monolith(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		providerconfig.Setup,
		allocation.Setup,
		network.Setup,
		pool.Setup,
		route.Setup,
		routetable.Setup,
		securitygroup.Setup,
		securityrule.Setup,
		subnet.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_monolith creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_monolith(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		providerconfig.SetupGated,
		allocation.SetupGated,
		network.SetupGated,
		pool.SetupGated,
		route.SetupGated,
		routetable.SetupGated,
		securitygroup.SetupGated,
		securityrule.SetupGated,
		subnet.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
