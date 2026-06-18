// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	accesskey "github.com/upbound/provider-nebius/internal/controller/cluster/iamv2/accesskey"
	project "github.com/upbound/provider-nebius/internal/controller/cluster/iamv2/project"
)

// Setup_iamv2 creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_iamv2(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		accesskey.Setup,
		project.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_iamv2 creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_iamv2(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		accesskey.SetupGated,
		project.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
