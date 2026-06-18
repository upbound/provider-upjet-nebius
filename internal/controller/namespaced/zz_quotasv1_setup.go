// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	quotaallowance "github.com/upbound/provider-nebius/internal/controller/namespaced/quotasv1/quotaallowance"
)

// Setup_quotasv1 creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_quotasv1(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		quotaallowance.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_quotasv1 creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_quotasv1(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		quotaallowance.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
