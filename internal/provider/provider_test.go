// Copyright (c) Veeblefetzer
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
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

func TestCombellProvider_Schema(t *testing.T) {
	p := &CombellProvider{version: "1.0.0"}

	req := provider.SchemaRequest{}
	resp := &provider.SchemaResponse{}

	p.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics)
	}

	// Check that expected attributes exist
	schema := resp.Schema

	apiKeyAttr, ok := schema.Attributes["api_key"]
	if !ok {
		t.Error("Schema missing 'api_key' attribute")
	} else {
		if !apiKeyAttr.IsOptional() {
			t.Error("api_key should be optional")
		}
		if !apiKeyAttr.IsSensitive() {
			t.Error("api_key should be sensitive")
		}
	}

	apiSecretAttr, ok := schema.Attributes["api_secret"]
	if !ok {
		t.Error("Schema missing 'api_secret' attribute")
	} else {
		if !apiSecretAttr.IsOptional() {
			t.Error("api_secret should be optional")
		}
		if !apiSecretAttr.IsSensitive() {
			t.Error("api_secret should be sensitive")
		}
	}

	baseURLAttr, ok := schema.Attributes["base_url"]
	if !ok {
		t.Error("Schema missing 'base_url' attribute")
	} else {
		if !baseURLAttr.IsOptional() {
			t.Error("base_url should be optional")
		}
		if baseURLAttr.IsSensitive() {
			t.Error("base_url should not be sensitive")
		}
	}
}

func TestCombellProvider_DataSources(t *testing.T) {
	p := &CombellProvider{version: "test"}

	dataSources := p.DataSources(context.Background())

	// Currently no data sources implemented
	if dataSources == nil {
		t.Error("DataSources() returned nil, expected empty slice")
	}
}

func TestCombellProvider_Resources(t *testing.T) {
	p := &CombellProvider{version: "test"}

	resources := p.Resources(context.Background())

	// Currently no resources implemented
	if resources == nil {
		t.Error("Resources() returned nil, expected empty slice")
	}
}

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"combell": providerserver.NewProtocol6WithError(New("test")()),
}
