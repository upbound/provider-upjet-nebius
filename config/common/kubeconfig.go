// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

// Package kubeconfig builds a standard kubeconfig from the observed Terraform
// state of a nebius_mk8s_v1_cluster so it can be consumed by provider-kubernetes.
package common

import (
	"fmt"

	"github.com/crossplane/crossplane-runtime/v2/pkg/fieldpath"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"k8s.io/client-go/tools/clientcmd"
)

const (
	// KubeconfigConnectionDetailsKey is the connection-details key under which the
	// default kubeconfig is stored. It targets the public endpoint for public
	// clusters and the private endpoint otherwise, mirroring the Nebius CLI. It
	// matches the default key expected by provider-kubernetes ProviderConfig
	// credential secretRefs.
	KubeconfigConnectionDetailsKey = "kubeconfig"
	// KubeconfigPublicConnectionDetailsKey is the connection-details key under
	// which the kubeconfig targeting the cluster's public endpoint is stored. It
	// is only present when the cluster exposes a public endpoint.
	KubeconfigPublicConnectionDetailsKey = "kubeconfig.public"
	// KubeconfigPrivateConnectionDetailsKey is the connection-details key under
	// which the kubeconfig targeting the cluster's private endpoint is stored.
	KubeconfigPrivateConnectionDetailsKey = "kubeconfig.private"

	// kubeContextName names the single context/cluster/user in each rendered
	// kubeconfig. It is an arbitrary fixed label: the kubeconfig holds one context,
	// so the name carries no external contract.
	kubeContextName = "context"
)

// MK8SClusterBuildKubeconfig renders kubeconfigs for a nebius_mk8s_v1_cluster
// observation and returns them as connection details. Every cluster has a private
// endpoint, so "kubeconfig.private" is always set; "kubeconfig.public" is set only
// when the cluster exposes a public endpoint. The default "kubeconfig" key targets
// the public endpoint when available and the private endpoint otherwise, mirroring
// the Nebius CLI. It returns nil when the private endpoint or CA certificate are
// not yet available, e.g. while the cluster is still being provisioned.
//
// Its signature matches config.AdditionalConnectionDetailsFn so it can be
// assigned directly to r.Sensitive.AdditionalConnectionDetailsFn.
func MK8SClusterBuildKubeconfig(attr map[string]any) (map[string][]byte, error) {
	paved := fieldpath.Pave(attr)

	publicEndpoint, _ := paved.GetString("status.control_plane.endpoints.public_endpoint")
	privateEndpoint, _ := paved.GetString("status.control_plane.endpoints.private_endpoint")
	caCertificate, _ := paved.GetString("status.control_plane.auth.cluster_ca_certificate")

	// A Nebius mk8s cluster always has a private endpoint; its absence (with the
	// CA) means the control plane is not provisioned yet.
	if caCertificate == "" || privateEndpoint == "" {
		return nil, nil
	}

	privateKubeconfig, err := renderKubeconfig(privateEndpoint, caCertificate)
	if err != nil {
		return nil, err
	}
	conn := map[string][]byte{
		KubeconfigPrivateConnectionDetailsKey: privateKubeconfig,
		// Default to the private endpoint; overridden below for public clusters.
		KubeconfigConnectionDetailsKey: privateKubeconfig,
	}

	if publicEndpoint != "" {
		publicKubeconfig, err := renderKubeconfig(publicEndpoint, caCertificate)
		if err != nil {
			return nil, err
		}
		conn[KubeconfigPublicConnectionDetailsKey] = publicKubeconfig
		// Public clusters default to the public endpoint, mirroring the Nebius CLI.
		conn[KubeconfigConnectionDetailsKey] = publicKubeconfig
	}

	return conn, nil
}

// renderKubeconfig serializes a single-context kubeconfig targeting server with
// the given CA certificate.
func renderKubeconfig(server, caCertificate string) ([]byte, error) {
	cfg := clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			kubeContextName: {
				Server:                   server,
				CertificateAuthorityData: []byte(caCertificate),
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			kubeContextName: {
				Cluster:  kubeContextName,
				AuthInfo: kubeContextName,
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			kubeContextName: {},
		},
		CurrentContext: kubeContextName,
	}

	raw, err := clientcmd.Write(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot serialize kubeconfig: %w", err)
	}
	return raw, nil
}
