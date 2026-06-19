// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package common

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// clusterAttr builds a Terraform state attribute map for a nebius_mk8s_v1_cluster
// observation. Empty endpoint arguments are omitted to model an unset endpoint.
func clusterAttr(ca, public, private string) map[string]any {
	endpoints := map[string]any{}
	if public != "" {
		endpoints["public_endpoint"] = public
	}
	if private != "" {
		endpoints["private_endpoint"] = private
	}
	return map[string]any{
		"status": map[string]any{
			"control_plane": map[string]any{
				"auth": map[string]any{
					"cluster_ca_certificate": ca,
				},
				"endpoints": endpoints,
			},
		},
	}
}

const testCABase64 = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM2akNDQWRLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJMk1EWXhPREUyTkRJeE5Gb1hEVE0yTURZeE5URTJORGN4TkZvd0ZURVRNQkVHQTFVRQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="

// wantKubeconfig renders the kubeconfig YAML expected for a given server. The
// single context/cluster/user is always named "context".
func wantKubeconfig(server string) []byte {
	return []byte(`apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ` + testCABase64 + `
    server: ` + server + `
  name: context
contexts:
- context:
    cluster: context
    user: context
  name: context
current-context: context
kind: Config
users:
- name: context
  user: {}
`)
}

func TestBuildKubeconfig(t *testing.T) {
	testCACertificate := `-----BEGIN CERTIFICATE-----
MIIC6jCCAdKgAwIBAgIBADANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwprdWJl
cm5ldGVzMB4XDTI2MDYxODE2NDIxNFoXDTM2MDYxNTE2NDcxNFowFTETMBEGA1UE
-----END CERTIFICATE-----
`
	testPublicEndpoint := "https://pu.mk8scluster-u00k7eaa4yg335q1sm.mk8s.us-central1.nebius.cloud:443"
	testPrivateEndpoint := "https://pr.mk8scluster-u00k7eaa4yg335q1sm.mk8s.us-central1.nebius.cloud:443"

	type args struct {
		ca      string
		public  string
		private string
	}
	cases := map[string]struct {
		args args
		want map[string][]byte
	}{
		"PublicClusterExposesBothEndpoints": {
			args: args{
				ca:      testCACertificate,
				public:  testPublicEndpoint,
				private: testPrivateEndpoint,
			},
			want: map[string][]byte{
				"kubeconfig":         wantKubeconfig(testPublicEndpoint),
				"kubeconfig.public":  wantKubeconfig(testPublicEndpoint),
				"kubeconfig.private": wantKubeconfig(testPrivateEndpoint),
			},
		},
		"PrivateClusterExposesPrivateEndpointOnly": {
			args: args{
				ca:      testCACertificate,
				private: testPrivateEndpoint,
			},
			want: map[string][]byte{
				"kubeconfig":         wantKubeconfig(testPrivateEndpoint),
				"kubeconfig.private": wantKubeconfig(testPrivateEndpoint),
			},
		},
		"ReturnsNilWhenNoPrivateEndpoint": {
			args: args{
				ca:     testCACertificate,
				public: testPublicEndpoint,
			},
			want: nil,
		},
		"ReturnsNilWhenNoCACertificate": {
			args: args{
				public:  testPublicEndpoint,
				private: testPrivateEndpoint,
			},
			want: nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := MK8SClusterBuildKubeconfig(clusterAttr(tc.args.ca, tc.args.public, tc.args.private))
			if err != nil {
				t.Fatalf("BuildKubeconfig returned error: %v", err)
			}
			// Render []byte values as strings so the diff shows the kubeconfig YAML.
			asString := cmp.Transformer("string", func(b []byte) string { return string(b) })
			if diff := cmp.Diff(tc.want, got, asString); diff != "" {
				t.Errorf("BuildKubeconfig() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
