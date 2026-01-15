// Copyright (c) Veeblefetzer
// SPDX-License-Identifier: MPL-2.0

package dns_record

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/io-veeblefetzer/terraform-provider-combell/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DNSRecordDataSource{}

// NewDNSRecordDataSource creates a new DNS record data source.
func NewDNSRecordDataSource() datasource.DataSource {
	return &DNSRecordDataSource{}
}

// DNSRecordDataSource defines the data source implementation.
type DNSRecordDataSource struct {
	client *client.Client
}

// DNSRecordDataSourceModel describes the data source data model.
type DNSRecordDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	DomainName types.String `tfsdk:"domain_name"`
	Type       types.String `tfsdk:"type"`
	RecordName types.String `tfsdk:"record_name"`
	TTL        types.Int64  `tfsdk:"ttl"`
	Content    types.String `tfsdk:"content"`
	Priority   types.Int64  `tfsdk:"priority"`
	Service    types.String `tfsdk:"service"`
	Weight     types.Int64  `tfsdk:"weight"`
	Target     types.String `tfsdk:"target"`
	Protocol   types.String `tfsdk:"protocol"`
	Port       types.Int64  `tfsdk:"port"`
}

// Metadata returns the data source type name.
func (d *DNSRecordDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

// Schema defines the schema for the data source.
func (d *DNSRecordDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a DNS record from a Combell domain zone.",
		MarkdownDescription: "Retrieves a DNS record from a Combell domain zone by its ID.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The unique identifier of the DNS record.",
				MarkdownDescription: "The unique identifier of the DNS record.",
				Required:            true,
			},
			"domain_name": schema.StringAttribute{
				Description:         "The domain name the record belongs to.",
				MarkdownDescription: "The domain name the record belongs to.",
				Required:            true,
			},
			"type": schema.StringAttribute{
				Description:         "The type of DNS record (A, AAAA, CAA, CNAME, MX, TXT, SRV, ALIAS, TLSA).",
				MarkdownDescription: "The type of DNS record (`A`, `AAAA`, `CAA`, `CNAME`, `MX`, `TXT`, `SRV`, `ALIAS`, `TLSA`).",
				Computed:            true,
			},
			"record_name": schema.StringAttribute{
				Description:         "The name of the record (subdomain).",
				MarkdownDescription: "The name of the record (subdomain).",
				Computed:            true,
			},
			"ttl": schema.Int64Attribute{
				Description:         "Time to live in seconds.",
				MarkdownDescription: "Time to live in seconds.",
				Computed:            true,
			},
			"content": schema.StringAttribute{
				Description:         "The content/value of the record.",
				MarkdownDescription: "The content/value of the record.",
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				Description:         "Priority for MX and SRV records.",
				MarkdownDescription: "Priority for MX and SRV records.",
				Computed:            true,
			},
			"service": schema.StringAttribute{
				Description:         "The service name for SRV records.",
				MarkdownDescription: "The service name for SRV records.",
				Computed:            true,
			},
			"weight": schema.Int64Attribute{
				Description:         "Weight for SRV records.",
				MarkdownDescription: "Weight for SRV records.",
				Computed:            true,
			},
			"target": schema.StringAttribute{
				Description:         "The target hostname for SRV records.",
				MarkdownDescription: "The target hostname for SRV records.",
				Computed:            true,
			},
			"protocol": schema.StringAttribute{
				Description:         "The protocol for SRV records.",
				MarkdownDescription: "The protocol for SRV records.",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				Description:         "The port number for SRV records.",
				MarkdownDescription: "The port number for SRV records.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *DNSRecordDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = c
}

// Read retrieves the DNS record data.
func (d *DNSRecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DNSRecordDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading DNS record", map[string]interface{}{
		"id":     data.ID.ValueString(),
		"domain": data.DomainName.ValueString(),
	})

	// Make API call
	_, body, err := d.client.Get(ctx, fmt.Sprintf("/dns/%s/records/%s", data.DomainName.ValueString(), data.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading DNS Record",
			fmt.Sprintf("Could not read DNS record: %s", err.Error()),
		)
		return
	}

	// Parse response
	var apiRecord client.DnsRecord
	if err := json.Unmarshal(body, &apiRecord); err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing DNS Record",
			fmt.Sprintf("Could not parse DNS record response: %s", err.Error()),
		)
		return
	}

	// Update model from API response
	data.Type = types.StringValue(apiRecord.Type)
	data.RecordName = types.StringValue(apiRecord.RecordName)
	data.TTL = types.Int64Value(int64(apiRecord.TTL))
	data.Content = types.StringValue(apiRecord.Content)

	if apiRecord.Priority != nil {
		data.Priority = types.Int64Value(int64(*apiRecord.Priority))
	} else {
		data.Priority = types.Int64Null()
	}

	if apiRecord.Service != "" {
		data.Service = types.StringValue(apiRecord.Service)
	} else {
		data.Service = types.StringNull()
	}

	if apiRecord.Weight != nil {
		data.Weight = types.Int64Value(int64(*apiRecord.Weight))
	} else {
		data.Weight = types.Int64Null()
	}

	if apiRecord.Target != "" {
		data.Target = types.StringValue(apiRecord.Target)
	} else {
		data.Target = types.StringNull()
	}

	if apiRecord.Protocol != "" {
		data.Protocol = types.StringValue(apiRecord.Protocol)
	} else {
		data.Protocol = types.StringNull()
	}

	if apiRecord.Port != nil {
		data.Port = types.Int64Value(int64(*apiRecord.Port))
	} else {
		data.Port = types.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
