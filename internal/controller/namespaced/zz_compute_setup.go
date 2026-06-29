// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	disk "github.com/upbound/provider-nebius/internal/controller/namespaced/compute/disk"
	filesystem "github.com/upbound/provider-nebius/internal/controller/namespaced/compute/filesystem"
	gpucluster "github.com/upbound/provider-nebius/internal/controller/namespaced/compute/gpucluster"
	instance "github.com/upbound/provider-nebius/internal/controller/namespaced/compute/instance"
)

// Setup_compute creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_compute(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		disk.Setup,
		filesystem.Setup,
		gpucluster.Setup,
		instance.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_compute creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_compute(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		disk.SetupGated,
		filesystem.SetupGated,
		gpucluster.SetupGated,
		instance.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
