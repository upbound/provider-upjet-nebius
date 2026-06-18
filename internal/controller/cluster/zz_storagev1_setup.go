// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	bucket "github.com/upbound/provider-nebius/internal/controller/cluster/storagev1/bucket"
	transfer "github.com/upbound/provider-nebius/internal/controller/cluster/storagev1/transfer"
)

// Setup_storagev1 creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_storagev1(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		bucket.Setup,
		transfer.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_storagev1 creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_storagev1(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		bucket.SetupGated,
		transfer.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
