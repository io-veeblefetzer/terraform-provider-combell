// Copyright (c) Veeblefetzer
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestNew(t *testing.T) {
	version := "test"
	newFn := New(version)
	
	if newFn == nil {
		t.Fatal("New returned nil")
	}
	
	p := newFn()
	if p == nil {
		t.Fatal("New factory function returned nil provider")
	}
	
	// Verify it implements the Provider interface
	var _ provider.Provider = p
}

func TestCombellProvider_Metadata(t *testing.T) {
	p := &CombellProvider{version: "1.0.0"}
	
	req := provider.MetadataRequest{}
	resp := &provider.MetadataResponse{}
	
	p.Metadata(nil, req, resp)
	
	if resp.TypeName != "combell" {
		t.Errorf("expected TypeName to be 'combell', got '%s'", resp.TypeName)
	}
	
	if resp.Version != "1.0.0" {
		t.Errorf("expected Version to be '1.0.0', got '%s'", resp.Version)
	}
}
