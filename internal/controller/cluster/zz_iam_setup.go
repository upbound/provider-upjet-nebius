// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	accesskey "github.com/upbound/provider-nebius/internal/controller/cluster/iam/accesskey"
	accesspermit "github.com/upbound/provider-nebius/internal/controller/cluster/iam/accesspermit"
	authpublickey "github.com/upbound/provider-nebius/internal/controller/cluster/iam/authpublickey"
	federatedcredentials "github.com/upbound/provider-nebius/internal/controller/cluster/iam/federatedcredentials"
	federation "github.com/upbound/provider-nebius/internal/controller/cluster/iam/federation"
	federationcertificate "github.com/upbound/provider-nebius/internal/controller/cluster/iam/federationcertificate"
	group "github.com/upbound/provider-nebius/internal/controller/cluster/iam/group"
	groupmembership "github.com/upbound/provider-nebius/internal/controller/cluster/iam/groupmembership"
	invitation "github.com/upbound/provider-nebius/internal/controller/cluster/iam/invitation"
	project "github.com/upbound/provider-nebius/internal/controller/cluster/iam/project"
	serviceaccount "github.com/upbound/provider-nebius/internal/controller/cluster/iam/serviceaccount"
)

// Setup_iam creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_iam(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		accesskey.Setup,
		accesspermit.Setup,
		authpublickey.Setup,
		federatedcredentials.Setup,
		federation.Setup,
		federationcertificate.Setup,
		group.Setup,
		groupmembership.Setup,
		invitation.Setup,
		project.Setup,
		serviceaccount.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_iam creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_iam(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		accesskey.SetupGated,
		accesspermit.SetupGated,
		authpublickey.SetupGated,
		federatedcredentials.SetupGated,
		federation.SetupGated,
		federationcertificate.SetupGated,
		group.SetupGated,
		groupmembership.SetupGated,
		invitation.SetupGated,
		project.SetupGated,
		serviceaccount.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
