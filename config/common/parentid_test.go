// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package common

import (
	"context"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/test"
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clustervpc "github.com/upbound/provider-nebius/apis/cluster/vpc/v1beta1"
)

func TestParentIDInitializer(t *testing.T) {
	ctx := context.Background()

	type args struct {
		forProvider  *string // spec.forProvider.parentId on the resource
		initProvider *string // spec.initProvider.parentId on the resource
		observeOnly  bool
		pcProjectID  *string // ProviderConfig.spec.projectID; nil means unset
	}
	type want struct {
		parentID  *string // resulting spec.forProvider.parentId
		getCalled bool
		patched   bool
		err       bool
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"DefaultsFromProviderConfigWhenEmpty": {
			args: args{pcProjectID: new("project-e00example")},
			want: want{parentID: new("project-e00example"), getCalled: true, patched: true},
		},
		"ForProviderSetIsNoOp": {
			args: args{forProvider: new("project-explicit"), pcProjectID: new("project-e00example")},
			want: want{parentID: new("project-explicit")},
		},
		"InitProviderSetIsNoOp": {
			args: args{initProvider: new("project-init"), pcProjectID: new("project-e00example")},
			want: want{},
		},
		"ObserveOnlyIsNoOp": {
			args: args{observeOnly: true, pcProjectID: new("project-e00example")},
			want: want{},
		},
		"ProviderConfigWithoutProjectIDErrors": {
			args: args{pcProjectID: nil},
			want: want{getCalled: true, err: true},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// A cluster-scoped (LegacyManaged) Network routes through the cluster
			// ProviderConfig path in providerConfigProjectID.
			n := &clustervpc.Network{}
			n.SetName("example")
			n.SetProviderConfigReference(&xpv1.Reference{Name: "default"})
			n.Spec.ForProvider.ParentID = tc.args.forProvider
			n.Spec.InitProvider.ParentID = tc.args.initProvider
			if tc.args.observeOnly {
				n.SetManagementPolicies(xpv1.ManagementPolicies{xpv1.ManagementActionObserve})
			}

			got := want{}
			mc := &test.MockClient{
				// Fill the unstructured ProviderConfig with the given
				// spec.projectID (omitted when nil) and record the call.
				MockGet: func(_ context.Context, _ client.ObjectKey, obj client.Object) error {
					got.getCalled = true
					u, ok := obj.(*unstructured.Unstructured)
					if !ok {
						return errors.Errorf("unexpected Get target %T", obj)
					}
					content := map[string]any{}
					if tc.args.pcProjectID != nil {
						content["spec"] = map[string]any{"projectID": *tc.args.pcProjectID}
					}
					u.Object = content
					return nil
				},
				MockApply: test.NewMockApplyFn(nil, func(_ runtime.ApplyConfiguration) error { got.patched = true; return nil }),
			}

			err := (&parentIDInitializer{kube: mc}).Initialize(ctx, n)
			got.err = err != nil
			got.parentID = n.Spec.ForProvider.ParentID

			if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(want{})); diff != "" {
				t.Errorf("Initialize() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
