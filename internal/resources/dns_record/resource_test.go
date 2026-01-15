// Copyright (c) Veeblefetzer
// SPDX-License-Identifier: MPL-2.0

package dns_record

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNewDNSRecordResource(t *testing.T) {
	r := NewDNSRecordResource()
	if r == nil {
		t.Fatal("NewDNSRecordResource() returned nil")
	}
}

func TestDNSRecordResource_Metadata(t *testing.T) {
	r := &DNSRecordResource{}

	req := resource.MetadataRequest{
		ProviderTypeName: "combell",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	expected := "combell_dns_record"
	if resp.TypeName != expected {
		t.Errorf("Metadata() TypeName = %s, want %s", resp.TypeName, expected)
	}
}

func TestDNSRecordResource_Schema(t *testing.T) {
	r := &DNSRecordResource{}

	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics)
	}

	schema := resp.Schema

	// Verify required attributes
	requiredAttrs := []string{"domain_name", "type"}
	for _, attr := range requiredAttrs {
		a, ok := schema.Attributes[attr]
		if !ok {
			t.Errorf("Schema missing required attribute: %s", attr)
			continue
		}
		if !a.IsRequired() {
			t.Errorf("Attribute %s should be required", attr)
		}
	}

	// Verify optional attributes
	optionalAttrs := []string{"record_name", "ttl", "content", "priority", "service", "weight", "target", "protocol", "port"}
	for _, attr := range optionalAttrs {
		a, ok := schema.Attributes[attr]
		if !ok {
			t.Errorf("Schema missing optional attribute: %s", attr)
			continue
		}
		if !a.IsOptional() && !a.IsComputed() {
			t.Errorf("Attribute %s should be optional or computed", attr)
		}
	}

	// Verify computed attributes
	computedAttrs := []string{"id"}
	for _, attr := range computedAttrs {
		a, ok := schema.Attributes[attr]
		if !ok {
			t.Errorf("Schema missing computed attribute: %s", attr)
			continue
		}
		if !a.IsComputed() {
			t.Errorf("Attribute %s should be computed", attr)
		}
	}
}

func TestNewDNSRecordDataSource(t *testing.T) {
	d := NewDNSRecordDataSource()
	if d == nil {
		t.Fatal("NewDNSRecordDataSource() returned nil")
	}
}

func TestDNSRecordDataSource_Metadata(t *testing.T) {
	d := &DNSRecordDataSource{}

	req := datasource.MetadataRequest{
		ProviderTypeName: "combell",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(context.Background(), req, resp)

	expected := "combell_dns_record"
	if resp.TypeName != expected {
		t.Errorf("Metadata() TypeName = %s, want %s", resp.TypeName, expected)
	}
}

func TestDNSRecordDataSource_Schema(t *testing.T) {
	d := &DNSRecordDataSource{}

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics)
	}

	schema := resp.Schema

	// Verify required attributes
	requiredAttrs := []string{"id", "domain_name"}
	for _, attr := range requiredAttrs {
		a, ok := schema.Attributes[attr]
		if !ok {
			t.Errorf("Schema missing required attribute: %s", attr)
			continue
		}
		if !a.IsRequired() {
			t.Errorf("Attribute %s should be required", attr)
		}
	}

	// Verify computed attributes
	computedAttrs := []string{"type", "record_name", "ttl", "content", "priority", "service", "weight", "target", "protocol", "port"}
	for _, attr := range computedAttrs {
		a, ok := schema.Attributes[attr]
		if !ok {
			t.Errorf("Schema missing computed attribute: %s", attr)
			continue
		}
		if !a.IsComputed() {
			t.Errorf("Attribute %s should be computed", attr)
		}
	}
}

func TestDNSRecordResource_BuildAPIRecord(t *testing.T) {
	r := &DNSRecordResource{}

	tests := []struct {
		name     string
		model    DNSRecordResourceModel
		wantType string
	}{
		{
			name: "A record",
			model: DNSRecordResourceModel{
				Type:       types.StringValue("A"),
				RecordName: types.StringValue("www"),
				TTL:        types.Int64Value(3600),
				Content:    types.StringValue("192.168.1.1"),
			},
			wantType: "A",
		},
		{
			name: "MX record",
			model: DNSRecordResourceModel{
				Type:       types.StringValue("MX"),
				RecordName: types.StringValue(""),
				TTL:        types.Int64Value(3600),
				Content:    types.StringValue("mail.example.com"),
				Priority:   types.Int64Value(10),
			},
			wantType: "MX",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.buildAPIRecord(tt.model)
			if result.Type != tt.wantType {
				t.Errorf("buildAPIRecord() Type = %s, want %s", result.Type, tt.wantType)
			}
		})
	}
}
