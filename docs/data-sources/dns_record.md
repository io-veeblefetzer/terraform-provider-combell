---
page_title: "combell_dns_record Data Source - terraform-provider-combell"
subcategory: "DNS"
description: |-
  Retrieves a DNS record from a Combell domain zone.
---

# combell_dns_record (Data Source)

Retrieves a DNS record from a Combell domain zone by its ID.

## Example Usage

```terraform
data "combell_dns_record" "www" {
  domain_name = "example.com"
  id          = "12345"
}

output "record_content" {
  value = data.combell_dns_record.www.content
}
```

## Schema

### Required

- `domain_name` (String) The domain name the record belongs to.
- `id` (String) The unique identifier of the DNS record.

### Read-Only

- `content` (String) The content/value of the record.
- `port` (Number) The port number for SRV records.
- `priority` (Number) Priority for MX and SRV records.
- `protocol` (String) The protocol for SRV records.
- `record_name` (String) The name of the record (subdomain).
- `service` (String) The service name for SRV records.
- `target` (String) The target hostname for SRV records.
- `ttl` (Number) Time to live in seconds.
- `type` (String) The type of DNS record (`A`, `AAAA`, `CAA`, `CNAME`, `MX`, `TXT`, `SRV`, `ALIAS`, `TLSA`).
- `weight` (Number) Weight for SRV records.
