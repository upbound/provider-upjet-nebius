// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	asymmetrickey "github.com/upbound/provider-nebius/internal/controller/namespaced/kmsv1/asymmetrickey"
	symmetrickey "github.com/upbound/provider-nebius/internal/controller/namespaced/kmsv1/symmetrickey"
)

// Setup_kmsv1 creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_kmsv1(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		asymmetrickey.Setup,
		symmetrickey.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_kmsv1 creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_kmsv1(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		asymmetrickey.SetupGated,
		symmetrickey.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
