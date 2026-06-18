// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package config

import (
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/types/name"
)

// GroupKindOverrides overrides the group and kind of the resource if it matches
// any entry in the GroupMap.
func GroupKindOverrides() config.ResourceOption {
	return func(r *config.Resource) {
		if f, ok := GroupMap[r.Name]; ok {
			r.ShortGroup, r.Kind = f(r.Name)
		}
	}
}

// GroupKindCalculator returns the correct group and kind name for given TF
// resource.
type GroupKindCalculator func(resource string) (string, string)

// ReplaceGroupWords uses given group as the group of the resource and removes
// a number of words in resource name before calculating the kind of the resource.
func ReplaceGroupWords(group string, count int) GroupKindCalculator {
	return func(resource string) (string, string) {
		// "nebius_iam_v1_access_permit": "iamv1" -> (iamv1, AccessPermit)
		words := strings.Split(strings.TrimPrefix(resource, "nebius_"), "_")
		snakeKind := strings.Join(words[count:], "_")
		return group, name.NewFromSnake(snakeKind).Camel
	}
}

// GroupMap contains all overrides we'd like to make to the default group search.
// It's written with data from TF Provider documentation.
// We want to remove version from occurying in Kind of the resource and instead end up in Group
var GroupMap = map[string]GroupKindCalculator{
	"nebius_vpc_v1_network":                ReplaceGroupWords("vpcv1", 2),
	"nebius_vpc_v1_pool":                   ReplaceGroupWords("vpcv1", 2),
	"nebius_vpc_v1_subnet":                 ReplaceGroupWords("vpcv1", 2),
	"nebius_vpc_v1_route_table":            ReplaceGroupWords("vpcv1", 2),
	"nebius_vpc_v1_allocation":             ReplaceGroupWords("vpcv1", 2),
	"nebius_vpc_v1_route":                  ReplaceGroupWords("vpcv1", 2),
	"nebius_vpc_v1_security_group":         ReplaceGroupWords("vpcv1", 2),
	"nebius_vpc_v1_security_rule":          ReplaceGroupWords("vpcv1", 2),
	"nebius_mk8s_v1_cluster":               ReplaceGroupWords("mk8sv1", 2),
	"nebius_mk8s_v1_node_group":            ReplaceGroupWords("mk8sv1", 2),
	"nebius_compute_v1_gpu_cluster":        ReplaceGroupWords("computev1", 2),
	"nebius_compute_v1_filesystem":         ReplaceGroupWords("computev1", 2),
	"nebius_compute_v1_disk":               ReplaceGroupWords("computev1", 2),
	"nebius_iam_v1_service_account":        ReplaceGroupWords("iamv1", 2),
	"nebius_iam_v1_group":                  ReplaceGroupWords("iamv1", 2),
	"nebius_iam_v1_group_membership":       ReplaceGroupWords("iamv1", 2),
	"nebius_iam_v2_access_key":             ReplaceGroupWords("iamv2", 2),
	"nebius_iam_v1_access_permit":          ReplaceGroupWords("iamv1", 2),
	"nebius_iam_v1_auth_public_key":        ReplaceGroupWords("iamv1", 2),
	"nebius_iam_v1_federated_credentials":  ReplaceGroupWords("iamv1", 2),
	"nebius_iam_v1_federation":             ReplaceGroupWords("iamv1", 2),
	"nebius_iam_v1_federation_certificate": ReplaceGroupWords("iamv1", 2),
	"nebius_iam_v1_invitation":             ReplaceGroupWords("iamv1", 2),
	"nebius_iam_v2_project":                ReplaceGroupWords("iamv2", 2),
	"nebius_dns_v1_zone":                   ReplaceGroupWords("dnsv1", 2),
	"nebius_dns_v1_record":                 ReplaceGroupWords("dnsv1", 2),
	"nebius_mysterybox_v1_secret":          ReplaceGroupWords("mysteryboxv1", 2),
	"nebius_mysterybox_v1_secret_version":  ReplaceGroupWords("mysteryboxv1", 2),
	"nebius_storage_v1_bucket":             ReplaceGroupWords("storagev1", 2),
	"nebius_storage_v1_transfer":           ReplaceGroupWords("storagev1", 2),
}
