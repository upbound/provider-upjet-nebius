// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	disk "github.com/upbound/provider-nebius/internal/controller/namespaced/computev1/disk"
	filesystem "github.com/upbound/provider-nebius/internal/controller/namespaced/computev1/filesystem"
	gpucluster "github.com/upbound/provider-nebius/internal/controller/namespaced/computev1/gpucluster"
	instance "github.com/upbound/provider-nebius/internal/controller/namespaced/computev1/instance"
	record "github.com/upbound/provider-nebius/internal/controller/namespaced/dnsv1/record"
	zone "github.com/upbound/provider-nebius/internal/controller/namespaced/dnsv1/zone"
	accesspermit "github.com/upbound/provider-nebius/internal/controller/namespaced/iamv1/accesspermit"
	authpublickey "github.com/upbound/provider-nebius/internal/controller/namespaced/iamv1/authpublickey"
	federatedcredentials "github.com/upbound/provider-nebius/internal/controller/namespaced/iamv1/federatedcredentials"
	federation "github.com/upbound/provider-nebius/internal/controller/namespaced/iamv1/federation"
	federationcertificate "github.com/upbound/provider-nebius/internal/controller/namespaced/iamv1/federationcertificate"
	group "github.com/upbound/provider-nebius/internal/controller/namespaced/iamv1/group"
	groupmembership "github.com/upbound/provider-nebius/internal/controller/namespaced/iamv1/groupmembership"
	invitation "github.com/upbound/provider-nebius/internal/controller/namespaced/iamv1/invitation"
	serviceaccount "github.com/upbound/provider-nebius/internal/controller/namespaced/iamv1/serviceaccount"
	accesskey "github.com/upbound/provider-nebius/internal/controller/namespaced/iamv2/accesskey"
	project "github.com/upbound/provider-nebius/internal/controller/namespaced/iamv2/project"
	asymmetrickey "github.com/upbound/provider-nebius/internal/controller/namespaced/kmsv1/asymmetrickey"
	symmetrickey "github.com/upbound/provider-nebius/internal/controller/namespaced/kmsv1/symmetrickey"
	cluster "github.com/upbound/provider-nebius/internal/controller/namespaced/mk8sv1/cluster"
	nodegroup "github.com/upbound/provider-nebius/internal/controller/namespaced/mk8sv1/nodegroup"
	secret "github.com/upbound/provider-nebius/internal/controller/namespaced/mysteryboxv1/secret"
	secretversion "github.com/upbound/provider-nebius/internal/controller/namespaced/mysteryboxv1/secretversion"
	providerconfig "github.com/upbound/provider-nebius/internal/controller/namespaced/providerconfig"
	quotaallowance "github.com/upbound/provider-nebius/internal/controller/namespaced/quotasv1/quotaallowance"
	registry "github.com/upbound/provider-nebius/internal/controller/namespaced/registryv1/registry"
	bucket "github.com/upbound/provider-nebius/internal/controller/namespaced/storagev1/bucket"
	transfer "github.com/upbound/provider-nebius/internal/controller/namespaced/storagev1/transfer"
	allocation "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/allocation"
	network "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/network"
	pool "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/pool"
	route "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/route"
	routetable "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/routetable"
	securitygroup "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/securitygroup"
	securityrule "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/securityrule"
	subnet "github.com/upbound/provider-nebius/internal/controller/namespaced/vpcv1/subnet"
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
		accesspermit.Setup,
		authpublickey.Setup,
		federatedcredentials.Setup,
		federation.Setup,
		federationcertificate.Setup,
		group.Setup,
		groupmembership.Setup,
		invitation.Setup,
		serviceaccount.Setup,
		accesskey.Setup,
		project.Setup,
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
		accesspermit.SetupGated,
		authpublickey.SetupGated,
		federatedcredentials.SetupGated,
		federation.SetupGated,
		federationcertificate.SetupGated,
		group.SetupGated,
		groupmembership.SetupGated,
		invitation.SetupGated,
		serviceaccount.SetupGated,
		accesskey.SetupGated,
		project.SetupGated,
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
