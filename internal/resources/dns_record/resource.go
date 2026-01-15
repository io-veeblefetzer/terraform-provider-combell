// Copyright (c) Veeblefetzer
// SPDX-License-Identifier: MPL-2.0

package dns_record

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/io-veeblefetzer/terraform-provider-combell/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DNSRecordResource{}
var _ resource.ResourceWithImportState = &DNSRecordResource{}

// NewDNSRecordResource creates a new DNS record resource.
func NewDNSRecordResource() resource.Resource {
	return &DNSRecordResource{}
}

// DNSRecordResource defines the resource implementation.
type DNSRecordResource struct {
	client *client.Client
}

// DNSRecordResourceModel describes the resource data model.
type DNSRecordResourceModel struct {
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

// Metadata returns the resource type name.
func (r *DNSRecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

// Schema defines the schema for the resource.
func (r *DNSRecordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages a DNS record in a Combell domain zone.",
		MarkdownDescription: "Manages a DNS record in a Combell domain zone.\n\nSupported record types: A, AAAA, CAA, CNAME, MX, TXT, SRV, ALIAS, TLSA.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The unique identifier of the DNS record.",
				MarkdownDescription: "The unique identifier of the DNS record.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_name": schema.StringAttribute{
				Description:         "The domain name the record belongs to.",
				MarkdownDescription: "The domain name the record belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description:         "The type of DNS record (A, AAAA, CAA, CNAME, MX, TXT, SRV, ALIAS, TLSA).",
				MarkdownDescription: "The type of DNS record (`A`, `AAAA`, `CAA`, `CNAME`, `MX`, `TXT`, `SRV`, `ALIAS`, `TLSA`).",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("A", "AAAA", "CAA", "CNAME", "MX", "TXT", "SRV", "ALIAS", "TLSA"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"record_name": schema.StringAttribute{
				Description:         "The name of the record (subdomain). Use '@' or empty string for the root domain.",
				MarkdownDescription: "The name of the record (subdomain). Use `@` or empty string for the root domain.",
				Optional:            true,
				Computed:            true,
			},
			"ttl": schema.Int64Attribute{
				Description:         "Time to live in seconds (60-86400). Default: 3600.",
				MarkdownDescription: "Time to live in seconds (60-86400). Default: `3600`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(3600),
			},
			"content": schema.StringAttribute{
				Description:         "The content/value of the record. For A records: IPv4 address. For AAAA: IPv6 address. For CNAME: target hostname. For MX: mail server hostname. For TXT: text content.",
				MarkdownDescription: "The content/value of the record.\n\n- **A**: IPv4 address\n- **AAAA**: IPv6 address\n- **CNAME**: target hostname\n- **MX**: mail server hostname\n- **TXT**: text content\n- **CAA**: `{flag} {tag} {ca}` format\n- **ALIAS**: target hostname\n- **TLSA**: `{usage} {selector} {matching_type} {data}` format",
				Optional:            true,
			},
			"priority": schema.Int64Attribute{
				Description:         "Priority for MX and SRV records. Lower values have higher priority.",
				MarkdownDescription: "Priority for MX and SRV records. Lower values have higher priority.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(10),
			},
			"service": schema.StringAttribute{
				Description:         "The service name for SRV records (e.g., _sip, _http).",
				MarkdownDescription: "The service name for SRV records (e.g., `_sip`, `_http`).",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"weight": schema.Int64Attribute{
				Description:         "Weight for SRV records with the same priority. Higher values are preferred.",
				MarkdownDescription: "Weight for SRV records with the same priority. Higher values are preferred.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
			"target": schema.StringAttribute{
				Description:         "The target hostname for SRV records.",
				MarkdownDescription: "The target hostname for SRV records.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				Description:         "The protocol for SRV records (TCP, UDP).",
				MarkdownDescription: "The protocol for SRV records (`TCP`, `UDP`).",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port": schema.Int64Attribute{
				Description:         "The port number for SRV records.",
				MarkdownDescription: "The port number for SRV records.",
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *DNSRecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

// Create creates a new DNS record.
func (r *DNSRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DNSRecordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the API request
	record := r.buildAPIRecord(data)

	tflog.Debug(ctx, "Creating DNS record", map[string]interface{}{
		"domain":      data.DomainName.ValueString(),
		"record_type": data.Type.ValueString(),
		"record_name": data.RecordName.ValueString(),
	})

	// Make API call
	apiResp, body, err := r.client.Post(ctx, fmt.Sprintf("/dns/%s/records", data.DomainName.ValueString()), record)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating DNS Record",
			fmt.Sprintf("Could not create DNS record: %s", err.Error()),
		)
		return
	}

	// Extract record ID from Location header
	location := client.GetLocationHeader(apiResp)
	if location == "" {
		resp.Diagnostics.AddError(
			"Error Creating DNS Record",
			"API did not return a Location header with the new record ID",
		)
		return
	}

	// Parse record ID from location (e.g., /v2/dns/example.com/records/123)
	parts := strings.Split(location, "/")
	if len(parts) > 0 {
		data.ID = types.StringValue(parts[len(parts)-1])
	}

	tflog.Debug(ctx, "Created DNS record", map[string]interface{}{
		"id":       data.ID.ValueString(),
		"response": string(body),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read retrieves the DNS record state.
func (r *DNSRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DNSRecordResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading DNS record", map[string]interface{}{
		"id":     data.ID.ValueString(),
		"domain": data.DomainName.ValueString(),
	})

	// Make API call
	_, body, err := r.client.Get(ctx, fmt.Sprintf("/dns/%s/records/%s", data.DomainName.ValueString(), data.ID.ValueString()))
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			// Resource no longer exists
			resp.State.RemoveResource(ctx)
			return
		}
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

	// Update state from API response
	r.updateModelFromAPI(&data, &apiRecord)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update modifies an existing DNS record.
func (r *DNSRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DNSRecordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the API request
	record := r.buildAPIRecord(data)

	tflog.Debug(ctx, "Updating DNS record", map[string]interface{}{
		"id":     data.ID.ValueString(),
		"domain": data.DomainName.ValueString(),
	})

	// Make API call
	_, _, err := r.client.Put(ctx, fmt.Sprintf("/dns/%s/records/%s", data.DomainName.ValueString(), data.ID.ValueString()), record)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating DNS Record",
			fmt.Sprintf("Could not update DNS record: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete removes the DNS record.
func (r *DNSRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DNSRecordResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting DNS record", map[string]interface{}{
		"id":     data.ID.ValueString(),
		"domain": data.DomainName.ValueString(),
	})

	// Make API call
	_, _, err := r.client.Delete(ctx, fmt.Sprintf("/dns/%s/records/%s", data.DomainName.ValueString(), data.ID.ValueString()))
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			// Resource already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting DNS Record",
			fmt.Sprintf("Could not delete DNS record: %s", err.Error()),
		)
		return
	}
}

// ImportState imports an existing DNS record.
func (r *DNSRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: domain_name/record_id
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be in format 'domain_name/record_id', got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_name"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

// buildAPIRecord converts the Terraform model to an API request.
func (r *DNSRecordResource) buildAPIRecord(data DNSRecordResourceModel) *client.DnsRecord {
	record := &client.DnsRecord{
		Type:       data.Type.ValueString(),
		RecordName: data.RecordName.ValueString(),
		TTL:        int(data.TTL.ValueInt64()),
		Content:    data.Content.ValueString(),
	}

	if !data.Priority.IsNull() && !data.Priority.IsUnknown() {
		priority := int(data.Priority.ValueInt64())
		record.Priority = &priority
	}

	if !data.Service.IsNull() && !data.Service.IsUnknown() {
		record.Service = data.Service.ValueString()
	}

	if !data.Weight.IsNull() && !data.Weight.IsUnknown() {
		weight := int(data.Weight.ValueInt64())
		record.Weight = &weight
	}

	if !data.Target.IsNull() && !data.Target.IsUnknown() {
		record.Target = data.Target.ValueString()
	}

	if !data.Protocol.IsNull() && !data.Protocol.IsUnknown() {
		record.Protocol = data.Protocol.ValueString()
	}

	if !data.Port.IsNull() && !data.Port.IsUnknown() {
		port := int(data.Port.ValueInt64())
		record.Port = &port
	}

	return record
}

// updateModelFromAPI updates the Terraform model from an API response.
func (r *DNSRecordResource) updateModelFromAPI(data *DNSRecordResourceModel, record *client.DnsRecord) {
	data.ID = types.StringValue(record.ID)
	data.Type = types.StringValue(record.Type)
	data.RecordName = types.StringValue(record.RecordName)
	data.TTL = types.Int64Value(int64(record.TTL))
	data.Content = types.StringValue(record.Content)

	if record.Priority != nil {
		data.Priority = types.Int64Value(int64(*record.Priority))
	}

	if record.Service != "" {
		data.Service = types.StringValue(record.Service)
	}

	if record.Weight != nil {
		data.Weight = types.Int64Value(int64(*record.Weight))
	}

	if record.Target != "" {
		data.Target = types.StringValue(record.Target)
	}

	if record.Protocol != "" {
		data.Protocol = types.StringValue(record.Protocol)
	}

	if record.Port != nil {
		data.Port = types.Int64Value(int64(*record.Port))
	}
}
