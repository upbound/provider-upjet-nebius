// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package config

import (
	"context"
	"testing"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	"github.com/google/go-cmp/cmp"

	"github.com/upbound/provider-nebius/config/common"
)

func TestProjectParentDefaultingApplied(t *testing.T) {
	// marker captures the configuration every project-parented resource must
	// carry: parent_id present and optional (so codegen drops the required-
	// parameter CEL rule) and at least one InitializerFn (the defaulting
	// initializer). It is compared against want per resource via go-cmp.
	type marker struct {
		Present        bool
		Optional       bool
		Required       bool
		HasInitializer bool
	}
	want := marker{Present: true, Optional: true, Required: false, HasInitializer: true}

	cases := map[string]struct {
		provider func(context.Context) (*ujconfig.Provider, error)
	}{
		"cluster":    {provider: GetProvider},
		"namespaced": {provider: GetProviderNamespaced},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			pc, err := tc.provider(context.Background())
			if err != nil {
				t.Fatalf("provider: %v", err)
			}
			for _, res := range common.ProjectParentedResources {
				got := marker{}
				if r, ok := pc.Resources[res]; ok {
					got.HasInitializer = len(r.InitializerFns) > 0
					if s, ok := r.TerraformResource.Schema["parent_id"]; ok {
						got.Present = true
						got.Optional = s.Optional
						got.Required = s.Required
					}
				}
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("%s: parent_id defaulting not applied (-want +got):\n%s", res, diff)
				}
			}
		})
	}
}
