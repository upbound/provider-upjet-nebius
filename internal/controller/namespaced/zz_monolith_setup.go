// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	filesystem "github.com/upbound/provider-nebius/internal/controller/namespaced/computev1/filesystem"
	gpucluster "github.com/upbound/provider-nebius/internal/controller/namespaced/computev1/gpucluster"
	serviceaccount "github.com/upbound/provider-nebius/internal/controller/namespaced/iamv1/serviceaccount"
	cluster "github.com/upbound/provider-nebius/internal/controller/namespaced/mk8sv1/cluster"
	nodegroup "github.com/upbound/provider-nebius/internal/controller/namespaced/mk8sv1/nodegroup"
	providerconfig "github.com/upbound/provider-nebius/internal/controller/namespaced/providerconfig"
	allocation "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/allocation"
	network "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/network"
	pool "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/pool"
	route "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/route"
	routetable "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/routetable"
	securitygroup "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/securitygroup"
	securityrule "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/securityrule"
	subnet "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/subnet"
)

// Setup_monolith creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_monolith(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		filesystem.Setup,
		gpucluster.Setup,
		serviceaccount.Setup,
		cluster.Setup,
		nodegroup.Setup,
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
		filesystem.SetupGated,
		gpucluster.SetupGated,
		serviceaccount.SetupGated,
		cluster.SetupGated,
		nodegroup.SetupGated,
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
