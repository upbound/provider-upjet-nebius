// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	secret "github.com/upbound/provider-nebius/internal/controller/namespaced/mysteryboxv1/secret"
	secretversion "github.com/upbound/provider-nebius/internal/controller/namespaced/mysteryboxv1/secretversion"
)

// Setup_mysteryboxv1 creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_mysteryboxv1(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		secret.Setup,
		secretversion.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_mysteryboxv1 creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_mysteryboxv1(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		secret.SetupGated,
		secretversion.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
