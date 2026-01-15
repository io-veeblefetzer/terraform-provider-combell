---
page_title: "combell_dns_record Resource - terraform-provider-combell"
subcategory: "DNS"
description: |-
  Manages a DNS record in a Combell domain zone.
---

# combell_dns_record (Resource)

Manages a DNS record in a Combell domain zone.

Supported record types: A, AAAA, CAA, CNAME, MX, TXT, SRV, ALIAS, TLSA.

## Example Usage

### A Record

```terraform
resource "combell_dns_record" "www" {
  domain_name = "example.com"
  type        = "A"
  record_name = "www"
  content     = "192.168.1.1"
  ttl         = 3600
}
```

### AAAA Record

```terraform
resource "combell_dns_record" "www_ipv6" {
  domain_name = "example.com"
  type        = "AAAA"
  record_name = "www"
  content     = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
  ttl         = 3600
}
```

### CNAME Record

```terraform
resource "combell_dns_record" "blog" {
  domain_name = "example.com"
  type        = "CNAME"
  record_name = "blog"
  content     = "www.example.com"
  ttl         = 3600
}
```

### MX Record

```terraform
resource "combell_dns_record" "mail" {
  domain_name = "example.com"
  type        = "MX"
  record_name = ""
  content     = "mail.example.com"
  priority    = 10
  ttl         = 3600
}
```

### TXT Record (SPF)

```terraform
resource "combell_dns_record" "spf" {
  domain_name = "example.com"
  type        = "TXT"
  record_name = ""
  content     = "v=spf1 include:_spf.google.com ~all"
  ttl         = 3600
}
```

### SRV Record

```terraform
resource "combell_dns_record" "sip" {
  domain_name = "example.com"
  type        = "SRV"
  service     = "_sip"
  protocol    = "TCP"
  target      = "sipserver.example.com"
  port        = 5060
  priority    = 10
  weight      = 5
  ttl         = 3600
}
```

## Schema

### Required

- `domain_name` (String) The domain name the record belongs to.
- `type` (String) The type of DNS record (`A`, `AAAA`, `CAA`, `CNAME`, `MX`, `TXT`, `SRV`, `ALIAS`, `TLSA`).

### Optional

- `content` (String) The content/value of the record.
  - **A**: IPv4 address
  - **AAAA**: IPv6 address
  - **CNAME**: target hostname
  - **MX**: mail server hostname
  - **TXT**: text content
  - **CAA**: `{flag} {tag} {ca}` format
  - **ALIAS**: target hostname
  - **TLSA**: `{usage} {selector} {matching_type} {data}` format
- `port` (Number) The port number for SRV records.
- `priority` (Number) Priority for MX and SRV records. Lower values have higher priority. Default: `10`.
- `protocol` (String) The protocol for SRV records (`TCP`, `UDP`).
- `record_name` (String) The name of the record (subdomain). Use `@` or empty string for the root domain.
- `service` (String) The service name for SRV records (e.g., `_sip`, `_http`).
- `target` (String) The target hostname for SRV records.
- `ttl` (Number) Time to live in seconds (60-86400). Default: `3600`.
- `weight` (Number) Weight for SRV records with the same priority. Higher values are preferred. Default: `0`.

### Read-Only

- `id` (String) The unique identifier of the DNS record.

## Import

DNS records can be imported using the format `domain_name/record_id`:

```shell
terraform import combell_dns_record.www "example.com/12345"
```
