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
	record "github.com/upbound/provider-nebius/internal/controller/namespaced/dns/record"
	zone "github.com/upbound/provider-nebius/internal/controller/namespaced/dns/zone"
	accesskey "github.com/upbound/provider-nebius/internal/controller/namespaced/iam/accesskey"
	accesspermit "github.com/upbound/provider-nebius/internal/controller/namespaced/iam/accesspermit"
	authpublickey "github.com/upbound/provider-nebius/internal/controller/namespaced/iam/authpublickey"
	federatedcredentials "github.com/upbound/provider-nebius/internal/controller/namespaced/iam/federatedcredentials"
	federation "github.com/upbound/provider-nebius/internal/controller/namespaced/iam/federation"
	federationcertificate "github.com/upbound/provider-nebius/internal/controller/namespaced/iam/federationcertificate"
	group "github.com/upbound/provider-nebius/internal/controller/namespaced/iam/group"
	groupmembership "github.com/upbound/provider-nebius/internal/controller/namespaced/iam/groupmembership"
	invitation "github.com/upbound/provider-nebius/internal/controller/namespaced/iam/invitation"
	project "github.com/upbound/provider-nebius/internal/controller/namespaced/iam/project"
	serviceaccount "github.com/upbound/provider-nebius/internal/controller/namespaced/iam/serviceaccount"
	asymmetrickey "github.com/upbound/provider-nebius/internal/controller/namespaced/kms/asymmetrickey"
	symmetrickey "github.com/upbound/provider-nebius/internal/controller/namespaced/kms/symmetrickey"
	cluster "github.com/upbound/provider-nebius/internal/controller/namespaced/mk8s/cluster"
	nodegroup "github.com/upbound/provider-nebius/internal/controller/namespaced/mk8s/nodegroup"
	secret "github.com/upbound/provider-nebius/internal/controller/namespaced/mysterybox/secret"
	secretversion "github.com/upbound/provider-nebius/internal/controller/namespaced/mysterybox/secretversion"
	providerconfig "github.com/upbound/provider-nebius/internal/controller/namespaced/providerconfig"
	quotaallowance "github.com/upbound/provider-nebius/internal/controller/namespaced/quotas/quotaallowance"
	registry "github.com/upbound/provider-nebius/internal/controller/namespaced/registry/registry"
	bucket "github.com/upbound/provider-nebius/internal/controller/namespaced/storage/bucket"
	transfer "github.com/upbound/provider-nebius/internal/controller/namespaced/storage/transfer"
	allocation "github.com/upbound/provider-nebius/internal/controller/namespaced/vpc/allocation"
	network "github.com/upbound/provider-nebius/internal/controller/namespaced/vpc/network"
	pool "github.com/upbound/provider-nebius/internal/controller/namespaced/vpc/pool"
	route "github.com/upbound/provider-nebius/internal/controller/namespaced/vpc/route"
	routetable "github.com/upbound/provider-nebius/internal/controller/namespaced/vpc/routetable"
	securitygroup "github.com/upbound/provider-nebius/internal/controller/namespaced/vpc/securitygroup"
	securityrule "github.com/upbound/provider-nebius/internal/controller/namespaced/vpc/securityrule"
	subnet "github.com/upbound/provider-nebius/internal/controller/namespaced/vpc/subnet"
)

// Setup_monolith creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_monolith(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		disk.Setup,
		filesystem.Setup,
		gpucluster.Setup,
		instance.Setup,
		record.Setup,
		zone.Setup,
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
		asymmetrickey.Setup,
		symmetrickey.Setup,
		cluster.Setup,
		nodegroup.Setup,
		secret.Setup,
		secretversion.Setup,
		providerconfig.Setup,
		quotaallowance.Setup,
		registry.Setup,
		bucket.Setup,
		transfer.Setup,
		allocation.Setup,
		network.Setup,
		pool.Setup,
		route.Setup,
		routetable.Setup,
		securitygroup.Setup,
		securityrule.Setup,
		subnet.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_monolith creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_monolith(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		disk.SetupGated,
		filesystem.SetupGated,
		gpucluster.SetupGated,
		instance.SetupGated,
		record.SetupGated,
		zone.SetupGated,
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
		asymmetrickey.SetupGated,
		symmetrickey.SetupGated,
		cluster.SetupGated,
		nodegroup.SetupGated,
		secret.SetupGated,
		secretversion.SetupGated,
		providerconfig.SetupGated,
		quotaallowance.SetupGated,
		registry.SetupGated,
		bucket.SetupGated,
		transfer.SetupGated,
		allocation.SetupGated,
		network.SetupGated,
		pool.SetupGated,
		route.SetupGated,
		routetable.SetupGated,
		securitygroup.SetupGated,
		securityrule.SetupGated,
		subnet.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
