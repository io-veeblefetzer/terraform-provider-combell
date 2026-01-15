// Copyright (c) Veeblefetzer
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/io-veeblefetzer/terraform-provider-combell/internal/client"
	"github.com/io-veeblefetzer/terraform-provider-combell/internal/resources/dns_record"
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
	APIKey    types.String `tfsdk:"api_key"`
	APISecret types.String `tfsdk:"api_secret"`
	BaseURL   types.String `tfsdk:"base_url"`
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
		MarkdownDescription: `
The Combell provider enables Terraform to manage resources in your Combell hosting account via the Combell API v2.

## Authentication

The provider uses HMAC authentication. You need to provide your API key and secret, which can be obtained from the Combell control panel.

Credentials can be provided in the provider configuration or via environment variables:

- ` + "`COMBELL_API_KEY`" + ` - API key for authentication
- ` + "`COMBELL_API_SECRET`" + ` - API secret for HMAC signature
- ` + "`COMBELL_BASE_URL`" + ` - Optional API base URL override
`,
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description:         "API key from Combell control panel. Can also be set via COMBELL_API_KEY environment variable.",
				MarkdownDescription: "API key from Combell control panel. Can also be set via `COMBELL_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"api_secret": schema.StringAttribute{
				Description:         "API secret from Combell control panel. Can also be set via COMBELL_API_SECRET environment variable.",
				MarkdownDescription: "API secret from Combell control panel. Can also be set via `COMBELL_API_SECRET` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"base_url": schema.StringAttribute{
				Description:         "API base URL. Defaults to https://api.combell.com/v2. Can also be set via COMBELL_BASE_URL environment variable.",
				MarkdownDescription: "API base URL. Defaults to `https://api.combell.com/v2`. Can also be set via `COMBELL_BASE_URL` environment variable.",
				Optional:            true,
			},
		},
	}
}

// Configure prepares a Combell API client for data sources and resources.
func (p *CombellProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config CombellProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get API key from config or environment variable
	apiKey := os.Getenv("COMBELL_API_KEY")
	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}

	// Get API secret from config or environment variable
	apiSecret := os.Getenv("COMBELL_API_SECRET")
	if !config.APISecret.IsNull() {
		apiSecret = config.APISecret.ValueString()
	}

	// Get base URL from config or environment variable
	baseURL := os.Getenv("COMBELL_BASE_URL")
	if !config.BaseURL.IsNull() {
		baseURL = config.BaseURL.ValueString()
	}

	// Validate required configuration
	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"The Combell API key is required. "+
				"Set it in the provider configuration or via the COMBELL_API_KEY environment variable.",
		)
	}

	if apiSecret == "" {
		resp.Diagnostics.AddError(
			"Missing API Secret",
			"The Combell API secret is required. "+
				"Set it in the provider configuration or via the COMBELL_API_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the Combell API client
	combellClient := client.NewClient(apiKey, apiSecret, baseURL)

	// Make the client available during DataSource and Resource type Configure methods.
	resp.DataSourceData = combellClient
	resp.ResourceData = combellClient
}

// DataSources defines the data sources implemented in the provider.
func (p *CombellProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Will be populated in Phase 2 and 3
	}
}

// Resources defines the resources implemented in the provider.
func (p *CombellProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		dns_record.NewDNSRecordResource,
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
