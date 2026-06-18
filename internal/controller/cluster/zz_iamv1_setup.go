// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	accesspermit "github.com/upbound/provider-nebius/internal/controller/cluster/iamv1/accesspermit"
	authpublickey "github.com/upbound/provider-nebius/internal/controller/cluster/iamv1/authpublickey"
	federatedcredentials "github.com/upbound/provider-nebius/internal/controller/cluster/iamv1/federatedcredentials"
	federation "github.com/upbound/provider-nebius/internal/controller/cluster/iamv1/federation"
	federationcertificate "github.com/upbound/provider-nebius/internal/controller/cluster/iamv1/federationcertificate"
	group "github.com/upbound/provider-nebius/internal/controller/cluster/iamv1/group"
	groupmembership "github.com/upbound/provider-nebius/internal/controller/cluster/iamv1/groupmembership"
	invitation "github.com/upbound/provider-nebius/internal/controller/cluster/iamv1/invitation"
	serviceaccount "github.com/upbound/provider-nebius/internal/controller/cluster/iamv1/serviceaccount"
)

// Setup_iamv1 creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_iamv1(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		accesspermit.Setup,
		authpublickey.Setup,
		federatedcredentials.Setup,
		federation.Setup,
		federationcertificate.Setup,
		group.Setup,
		groupmembership.Setup,
		invitation.Setup,
		serviceaccount.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_iamv1 creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_iamv1(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		accesspermit.SetupGated,
		authpublickey.SetupGated,
		federatedcredentials.SetupGated,
		federation.SetupGated,
		federationcertificate.SetupGated,
		group.SetupGated,
		groupmembership.SetupGated,
		invitation.SetupGated,
		serviceaccount.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
