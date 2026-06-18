// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	record "github.com/upbound/provider-nebius/internal/controller/namespaced/dnsv1/record"
	zone "github.com/upbound/provider-nebius/internal/controller/namespaced/dnsv1/zone"
)

// Setup_dnsv1 creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_dnsv1(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		record.Setup,
		zone.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_dnsv1 creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_dnsv1(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		record.SetupGated,
		zone.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
