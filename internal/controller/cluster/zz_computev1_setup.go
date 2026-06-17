// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	filesystem "github.com/upbound/provider-nebius/internal/controller/cluster/computev1/filesystem"
	gpucluster "github.com/upbound/provider-nebius/internal/controller/cluster/computev1/gpucluster"
)

// Setup_computev1 creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_computev1(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		filesystem.Setup,
		gpucluster.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_computev1 creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_computev1(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		filesystem.SetupGated,
		gpucluster.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
