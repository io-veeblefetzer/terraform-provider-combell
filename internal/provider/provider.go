// Copyright (c) Veeblefetzer
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure CombellProvider satisfies various provider interfaces.
var _ provider.Provider = &CombellProvider{}

// CombellProvider defines the provider implementation.
type CombellProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// CombellProviderModel describes the provider data model.
type CombellProviderModel struct {
	// Will be populated in Issue #3
}

// Metadata returns the provider type name.
func (p *CombellProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "combell"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *CombellProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provider for managing Combell resources via the Combell API v2.",
		Attributes:  map[string]schema.Attribute{
			// Will be populated in Issue #3
		},
	}
}

// Configure prepares a Combell API client for data sources and resources.
func (p *CombellProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Will be implemented in Issue #3
}

// DataSources defines the data sources implemented in the provider.
func (p *CombellProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Will be populated in later issues
	}
}

// Resources defines the resources implemented in the provider.
func (p *CombellProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Will be populated in later issues
	}
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CombellProvider{
			version: version,
		}
	}
}
