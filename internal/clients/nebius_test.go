package clients

import (
	"strings"
	"testing"

	internalconfig "github.com/upbound/provider-nebius/internal/config"
)

func TestBuildConfiguration(t *testing.T) {
	str := func(s string) *string { return &s }
	tok := internalconfig.IdentityTypeToken
	sa := internalconfig.IdentityTypeServiceAccount

	tests := []struct {
		name            string
		creds           map[string]string
		identityType    internalconfig.IdentityType
		parentID        *string
		wantErr         bool
		wantErrContains string
		checkCfg        func(t *testing.T, cfg map[string]any)
	}{
		{
			name:         "token success",
			creds:        map[string]string{"token": "t1.abc"}, //nolint: goconst
			identityType: tok,
			checkCfg: func(t *testing.T, cfg map[string]any) {
				if cfg["token"] != "t1.abc" {
					t.Errorf("token: got %v", cfg["token"])
				}
				if _, ok := cfg["service_account"]; ok {
					t.Error("service_account should not be set")
				}
			},
		},
		{
			name:            "token empty string -> error",
			creds:           map[string]string{"token": ""},
			identityType:    tok,
			wantErr:         true,
			wantErrContains: `identityType is Token but credentials secret has no "token" key`,
		},
		{
			name:            "token missing key -> error",
			creds:           map[string]string{},
			identityType:    tok,
			wantErr:         true,
			wantErrContains: `identityType is Token but credentials secret has no "token" key`,
		},
		{
			name: "sa key all three fields",
			creds: map[string]string{
				"account_id":    "sa-123",
				"public_key_id": "pk-456",
				"private_key":   "-----BEGIN RSA PRIVATE KEY-----\n...",
			},
			identityType: sa,
			checkCfg: func(t *testing.T, cfg map[string]any) {
				saMap, ok := cfg["service_account"].(map[string]any)
				if !ok {
					t.Fatalf("service_account missing or wrong type: %T", cfg["service_account"])
				}
				if saMap["account_id"] != "sa-123" {
					t.Errorf("account_id: got %v", saMap["account_id"])
				}
				if saMap["public_key_id"] != "pk-456" {
					t.Errorf("public_key_id: got %v", saMap["public_key_id"])
				}
				if saMap["private_key"] != "-----BEGIN RSA PRIVATE KEY-----\n..." {
					t.Errorf("private_key: got %v", saMap["private_key"])
				}
				if _, ok := cfg["token"]; ok {
					t.Error("token should not be set")
				}
			},
		},
		{
			name:            "sa key all empty -> error",
			creds:           map[string]string{},
			identityType:    sa,
			wantErr:         true,
			wantErrContains: "identityType is ServiceAccount but credentials secret has none of account_id, public_key_id, private_key",
		},
		{
			name:         "sa partial - account_id only",
			creds:        map[string]string{"account_id": "sa-123"},
			identityType: sa,
			checkCfg: func(t *testing.T, cfg map[string]any) {
				saMap, ok := cfg["service_account"].(map[string]any)
				if !ok {
					t.Fatalf("service_account missing or wrong type")
				}
				if saMap["account_id"] != "sa-123" {
					t.Errorf("account_id: got %v", saMap["account_id"])
				}
				if saMap["public_key_id"] != "" {
					t.Errorf("public_key_id: expected empty, got %v", saMap["public_key_id"])
				}
				if saMap["private_key"] != "" {
					t.Errorf("private_key: expected empty, got %v", saMap["private_key"])
				}
			},
		},
		{
			name:         "sa partial - public_key_id only",
			creds:        map[string]string{"public_key_id": "pk-456"},
			identityType: sa,
			checkCfg: func(t *testing.T, cfg map[string]any) {
				saMap, ok := cfg["service_account"].(map[string]any)
				if !ok {
					t.Fatalf("service_account missing or wrong type")
				}
				if saMap["public_key_id"] != "pk-456" {
					t.Errorf("public_key_id: got %v", saMap["public_key_id"])
				}
				if saMap["account_id"] != "" {
					t.Errorf("account_id: expected empty, got %v", saMap["account_id"])
				}
				if saMap["private_key"] != "" {
					t.Errorf("private_key: expected empty, got %v", saMap["private_key"])
				}
			},
		},
		{
			name:         "sa partial - private_key only",
			creds:        map[string]string{"private_key": "-----BEGIN RSA PRIVATE KEY-----"},
			identityType: sa,
			checkCfg: func(t *testing.T, cfg map[string]any) {
				saMap, ok := cfg["service_account"].(map[string]any)
				if !ok {
					t.Fatalf("service_account missing or wrong type")
				}
				if saMap["private_key"] != "-----BEGIN RSA PRIVATE KEY-----" {
					t.Errorf("private_key: got %v", saMap["private_key"])
				}
				if saMap["account_id"] != "" {
					t.Errorf("account_id: expected empty, got %v", saMap["account_id"])
				}
				if saMap["public_key_id"] != "" {
					t.Errorf("public_key_id: expected empty, got %v", saMap["public_key_id"])
				}
			},
		},
		{
			name:            "unknown identity type -> error",
			creds:           map[string]string{"token": "t1.abc"},
			identityType:    internalconfig.IdentityType("Federated"),
			wantErr:         true,
			wantErrContains: `unknown identityType "Federated"`,
		},
		{
			name:         "token with parentID",
			creds:        map[string]string{"token": "t1.abc"},
			identityType: tok,
			parentID:     str("project-xyz"),
			checkCfg: func(t *testing.T, cfg map[string]any) {
				if cfg["parent_id"] != "project-xyz" {
					t.Errorf("parent_id: got %v", cfg["parent_id"])
				}
			},
		},
		{
			name:         "token without parentID -> no parent_id key",
			creds:        map[string]string{"token": "t1.abc"},
			identityType: tok,
			checkCfg: func(t *testing.T, cfg map[string]any) {
				if _, ok := cfg["parent_id"]; ok {
					t.Error("parent_id should not be set when parentID is nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := buildConfiguration(tt.creds, tt.identityType, tt.parentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr=%v got err=%v", tt.wantErr, err)
			}
			if err != nil && tt.wantErrContains != "" && !strings.Contains(err.Error(), tt.wantErrContains) {
				t.Errorf("error %q does not contain %q", err.Error(), tt.wantErrContains)
			}
			if err == nil && tt.checkCfg != nil {
				tt.checkCfg(t, cfg)
			}
		})
	}
}
