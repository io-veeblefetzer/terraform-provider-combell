# Terraform Provider Combell - Project Specification

## Overview

This document outlines the project plan for developing a Terraform provider for the Combell API (v2). The initial scope focuses on two primary resource groups: **DNS Records** and **Domains**.

## API Reference

- **Base URL**: `https://api.combell.com/v2`
- **Authentication**: HMAC-based authentication
- **OpenAPI Spec**: `resources/combell-openapi-spec.json`

---

## Phase 1: DNS Records

### Resource: `combell_dns_record`

Manages DNS records for a domain in Combell.

#### API Endpoints

| Operation | Method | Endpoint |
|-----------|--------|----------|
| List all records | GET | `/dns/{domainName}/records` |
| Create record | POST | `/dns/{domainName}/records` |
| Get record | GET | `/dns/{domainName}/records/{recordId}` |
| Update record | PUT | `/dns/{domainName}/records/{recordId}` |
| Delete record | DELETE | `/dns/{domainName}/records/{recordId}` |

#### Schema Definition

```hcl
resource "combell_dns_record" "example" {
  domain_name = "example.com"        # Required - The domain name
  type        = "A"                  # Required - Record type (A, AAAA, CAA, CNAME, MX, TXT, SRV, ALIAS, TLSA)
  record_name = "www"                # Optional - Host name/alias (empty or '@' equals domain name)
  ttl         = 3600                 # Optional - Time to live (60-86400), default 3600
  content     = "192.168.1.1"        # Required for most types - Variable data based on record type
  priority    = 10                   # Optional - For MX/SRV records (0-9999), default 10
  
  # SRV-specific fields
  service     = "_sip"               # Optional - Service name for SRV records
  protocol    = "TCP"                # Optional - Protocol for SRV records (TCP, UDP), default TCP
  weight      = 0                    # Optional - Weight for SRV records, default 0
  target      = "sip.example.com"    # Optional - Target hostname for SRV records
  port        = 5060                 # Optional - Port for SRV records
}
```

#### Attribute Reference

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `domain_name` | string | Yes | The domain name for the DNS record |
| `type` | string | Yes | Record type: A, AAAA, CAA, CNAME, MX, TXT, SRV, ALIAS, TLSA |
| `record_name` | string | No | Host name or alias. Empty or '@' equals the domain name |
| `ttl` | int | No | Time to live in seconds (60-86400). Default: 3600 |
| `content` | string | Conditional | Variable data depending on record type (see content rules below) |
| `priority` | int | No | Priority for MX/SRV records (0-9999). Default: 10 |
| `service` | string | No | Service name for SRV records (e.g., `_sip`) |
| `protocol` | string | No | Protocol for SRV records: TCP, UDP. Default: TCP |
| `weight` | int | No | Weight for SRV records with same priority. Default: 0 |
| `target` | string | No | Target hostname for SRV records |
| `port` | int | No | Port number for SRV records |

#### Content Field Rules

| Type | Content Description |
|------|---------------------|
| A | IPv4 address |
| AAAA | IPv6 address |
| CNAME | Canonical name of an alias |
| MX | Fully qualified domain name of a mail host |
| TXT | Free form text data |
| CAA | Format: `{flag} {tag} {ca}` (e.g., `0 issue "letsencrypt.org"`) |
| ALIAS | Canonical name of an alias |
| TLSA | Format: `{usage} {selector} {matching_type} {data}` |
| SRV | Not applicable (use service, target, port, weight fields) |

#### Computed Attributes (Read-Only)

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | string | The unique identifier of the DNS record |

#### Import

DNS records can be imported using the domain name and record ID:

```bash
terraform import combell_dns_record.example example.com/record-id-123
```

---

### Data Source: `combell_dns_record`

Retrieves a specific DNS record by ID.

#### Schema Definition

```hcl
data "combell_dns_record" "example" {
  domain_name = "example.com"
  record_id   = "record-id-123"
}
```

#### Attribute Reference

| Attribute | Type | Description |
|-----------|------|-------------|
| `domain_name` | string | The domain name |
| `record_id` | string | The record ID |
| `type` | string | Record type |
| `record_name` | string | Host name or alias |
| `ttl` | int | Time to live |
| `content` | string | Record content |
| `priority` | int | Priority (MX/SRV) |
| `service` | string | Service name (SRV) |
| `protocol` | string | Protocol (SRV) |
| `weight` | int | Weight (SRV) |
| `target` | string | Target (SRV) |
| `port` | int | Port (SRV) |

---

### Data Source: `combell_dns_records`

Retrieves all DNS records for a domain with optional filtering.

#### Schema Definition

```hcl
data "combell_dns_records" "all" {
  domain_name = "example.com"
  
  # Optional filters
  type        = "A"                  # Filter by record type
  record_name = "www"                # Filter by record name
  service     = "_sip"               # Filter by service (SRV only)
}
```

#### Attribute Reference

| Attribute | Type | Description |
|-----------|------|-------------|
| `domain_name` | string | The domain name to query |
| `type` | string | Optional filter by record type |
| `record_name` | string | Optional filter by record name |
| `service` | string | Optional filter by service (SRV records) |
| `records` | list | List of DNS records matching the criteria |

#### Output: `records` Block

Each record in the `records` list contains:

```hcl
records = [
  {
    id          = "record-id-123"
    type        = "A"
    record_name = "www"
    ttl         = 3600
    content     = "192.168.1.1"
    priority    = null
    service     = null
    protocol    = null
    weight      = null
    target      = null
    port        = null
  }
]
```

---

## Phase 2: Domains

### Data Source: `combell_domains`

Lists all domains in the Combell account.

#### API Endpoint

| Operation | Method | Endpoint |
|-----------|--------|----------|
| List domains | GET | `/domains` |

#### Schema Definition

```hcl
data "combell_domains" "all" {
  # Optional pagination
  skip = 0
  take = 100
}
```

#### Attribute Reference

| Attribute | Type | Description |
|-----------|------|-------------|
| `skip` | int | Number of items to skip (pagination) |
| `take` | int | Number of items to return (pagination) |
| `domains` | list | List of domains |
| `total_results` | int | Total number of domains available |

#### Output: `domains` Block

```hcl
domains = [
  {
    domain_name     = "example.com"
    expiration_date = "2025-12-31T23:59:59Z"
    will_renew      = true
  }
]
```

---

### Data Source: `combell_domain`

Retrieves detailed information about a specific domain.

#### API Endpoint

| Operation | Method | Endpoint |
|-----------|--------|----------|
| Get domain details | GET | `/domains/{domainName}` |

#### Schema Definition

```hcl
data "combell_domain" "example" {
  domain_name = "example.com"
}
```

#### Attribute Reference

| Attribute | Type | Description |
|-----------|------|-------------|
| `domain_name` | string | The domain name |
| `expiration_date` | string | Domain expiration date (ISO 8601) |
| `will_renew` | bool | Whether domain will auto-renew |
| `can_toggle_renew` | bool | Whether renew state can be changed |
| `name_servers` | list | List of nameserver objects |
| `registrant` | object | Registrant information |

#### Output: `name_servers` Block

```hcl
name_servers = [
  {
    name = "ns1.combell.net"
    ip   = "195.130.131.1"
  }
]
```

#### Output: `registrant` Block

```hcl
registrant = {
  first_name        = "John"
  last_name         = "Doe"
  address           = "123 Main St"
  postal_code       = "1000"
  city              = "Brussels"
  country_code      = "BE"
  email             = "john@example.com"
  phone             = "+32.123456789"
  fax               = null
  language_code     = "en"
  company_name      = "Example Corp"
  enterprise_number = "BE0123456789"
}
```

---

### Resource: `combell_domain_nameservers`

Manages the nameservers for a domain.

#### API Endpoint

| Operation | Method | Endpoint |
|-----------|--------|----------|
| Update nameservers | PUT | `/domains/{domainName}/nameservers` |

#### Schema Definition

```hcl
resource "combell_domain_nameservers" "example" {
  domain_name = "example.com"
  
  name_servers = [
    "ns1.example.com",
    "ns2.example.com"
  ]
}
```

#### Attribute Reference

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `domain_name` | string | Yes | The domain name |
| `name_servers` | list(string) | Yes | List of nameserver hostnames |

#### Import

Domain nameservers can be imported using the domain name:

```bash
terraform import combell_domain_nameservers.example example.com
```

#### Notes

- This resource only manages nameservers; domain registration/transfer is out of scope for this phase
- Reading the current state requires fetching from the domain details endpoint
- The resource implements Create, Read, Update operations (Delete reverts to defaults or is no-op)

---

### Resource: `combell_domain_renew`

Manages the auto-renewal state for a domain.

#### API Endpoint

| Operation | Method | Endpoint |
|-----------|--------|----------|
| Update renew state | PUT | `/domains/{domainName}/renew` |

#### Schema Definition

```hcl
resource "combell_domain_renew" "example" {
  domain_name = "example.com"
  will_renew  = true
}
```

#### Attribute Reference

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `domain_name` | string | Yes | The domain name |
| `will_renew` | bool | Yes | Whether the domain should auto-renew |

#### Computed Attributes (Read-Only)

| Attribute | Type | Description |
|-----------|------|-------------|
| `can_toggle_renew` | bool | Whether the renew state can currently be changed |

#### Import

Domain renewal configuration can be imported using the domain name:

```bash
terraform import combell_domain_renew.example example.com
```

#### Notes

- Changing `will_renew` is only allowed when `can_toggle_renew` is `true`
- Conditions for `can_toggle_renew`:
  - No unpaid invoices for the domain
  - Renewal won't start within 1 month
- Requires finance role for the requesting user

---

## Provider Configuration

### Schema Definition

```hcl
provider "combell" {
  api_key    = "your-api-key"        # Required - API key from control panel
  api_secret = "your-api-secret"     # Required - API secret from control panel
  base_url   = "https://api.combell.com/v2"  # Optional - API base URL
}
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `COMBELL_API_KEY` | API key for authentication |
| `COMBELL_API_SECRET` | API secret for HMAC signature |
| `COMBELL_BASE_URL` | Optional API base URL override |

### Authentication

The provider implements HMAC authentication as specified in the Combell API documentation:

1. Concatenate: `apikey + method + path + timestamp + nonce + content_hash`
2. Hash with SHA-256 using the API secret
3. Base64 encode the result
4. Send in Authorization header: `hmac apikey:signature:nonce:timestamp`

---

## Project Structure

```
terraform-provider-combell/
├── go.mod
├── go.sum
├── main.go
├── SPEC-TERRAFORM-PROVIDER-COMBELL.md
├── resources/
│   └── combell-openapi-spec.json
├── internal/
│   ├── provider/
│   │   ├── provider.go              # Provider configuration
│   │   └── provider_test.go
│   ├── client/
│   │   ├── client.go                # HTTP client with HMAC auth
│   │   ├── dns.go                   # DNS API methods
│   │   ├── domains.go               # Domains API methods
│   │   └── models.go                # API request/response models
│   └── resources/
│       ├── dns_record/
│       │   ├── resource.go          # combell_dns_record resource
│       │   ├── data_source.go       # combell_dns_record data source
│       │   ├── data_source_list.go  # combell_dns_records data source
│       │   └── schema.go            # Shared schema definitions
│       └── domain/
│           ├── data_source.go       # combell_domain data source
│           ├── data_source_list.go  # combell_domains data source
│           ├── nameservers.go       # combell_domain_nameservers resource
│           └── renew.go             # combell_domain_renew resource
├── examples/
│   ├── provider/
│   │   └── provider.tf
│   ├── resources/
│   │   ├── dns_record/
│   │   │   └── main.tf
│   │   ├── domain_nameservers/
│   │   │   └── main.tf
│   │   └── domain_renew/
│   │       └── main.tf
│   └── data-sources/
│       ├── dns_record/
│       │   └── main.tf
│       ├── dns_records/
│       │   └── main.tf
│       ├── domain/
│       │   └── main.tf
│       └── domains/
│           └── main.tf
└── docs/
    ├── index.md
    ├── resources/
    │   ├── dns_record.md
    │   ├── domain_nameservers.md
    │   └── domain_renew.md
    └── data-sources/
        ├── dns_record.md
        ├── dns_records.md
        ├── domain.md
        └── domains.md
```

---

## Implementation Phases

### Phase 1: Foundation
- [ ] Set up Go project structure
- [ ] Implement HMAC authentication client
- [ ] Create provider configuration with validation
- [ ] Add unit tests for authentication

### Phase 2: DNS Records
- [ ] Implement `combell_dns_record` resource (CRUD)
- [ ] Implement `combell_dns_record` data source
- [ ] Implement `combell_dns_records` data source (list with filters)
- [ ] Add acceptance tests for DNS resources
- [ ] Create documentation and examples

### Phase 3: Domains
- [ ] Implement `combell_domains` data source (list)
- [ ] Implement `combell_domain` data source (details)
- [ ] Implement `combell_domain_nameservers` resource
- [ ] Implement `combell_domain_renew` resource
- [ ] Add acceptance tests for domain resources
- [ ] Create documentation and examples

### Phase 4: Release
- [ ] Complete all documentation
- [ ] Set up CI/CD pipeline
- [ ] Create GitHub releases
- [ ] Publish to Terraform Registry

---

## Example Usage

### Complete DNS Management Example

```hcl
terraform {
  required_providers {
    combell = {
      source  = "veeblefetzer/combell"
      version = "~> 1.0"
    }
  }
}

provider "combell" {
  api_key    = var.combell_api_key
  api_secret = var.combell_api_secret
}

# Fetch domain information
data "combell_domain" "main" {
  domain_name = "example.com"
}

# Create A record for www
resource "combell_dns_record" "www" {
  domain_name = data.combell_domain.main.domain_name
  type        = "A"
  record_name = "www"
  content     = "192.168.1.100"
  ttl         = 3600
}

# Create MX records for email
resource "combell_dns_record" "mx_primary" {
  domain_name = data.combell_domain.main.domain_name
  type        = "MX"
  record_name = "@"
  content     = "mail.example.com"
  priority    = 10
  ttl         = 3600
}

resource "combell_dns_record" "mx_secondary" {
  domain_name = data.combell_domain.main.domain_name
  type        = "MX"
  record_name = "@"
  content     = "mail2.example.com"
  priority    = 20
  ttl         = 3600
}

# Create TXT record for SPF
resource "combell_dns_record" "spf" {
  domain_name = data.combell_domain.main.domain_name
  type        = "TXT"
  record_name = "@"
  content     = "v=spf1 include:_spf.google.com ~all"
  ttl         = 3600
}

# Configure custom nameservers
resource "combell_domain_nameservers" "custom" {
  domain_name = data.combell_domain.main.domain_name
  name_servers = [
    "ns1.cloudflare.com",
    "ns2.cloudflare.com"
  ]
}

# Enable auto-renewal
resource "combell_domain_renew" "auto" {
  domain_name = data.combell_domain.main.domain_name
  will_renew  = true
}

# Output domain details
output "domain_info" {
  value = {
    domain          = data.combell_domain.main.domain_name
    expires         = data.combell_domain.main.expiration_date
    will_renew      = data.combell_domain.main.will_renew
    can_toggle      = data.combell_domain.main.can_toggle_renew
  }
}

output "www_record_id" {
  value = combell_dns_record.www.id
}
```

---

## Error Handling

The provider should handle the following HTTP status codes appropriately:

| Status Code | Handling |
|-------------|----------|
| 200 OK | Success, return data |
| 201 Created | Resource created, read Location header |
| 204 No Content | Success, no data to return |
| 400 Bad Request | Return validation errors to user |
| 401 Unauthorized | Authentication failure, check credentials |
| 403 Forbidden | Permission denied |
| 404 Not Found | Resource not found (may indicate drift) |
| 429 Too Many Requests | Implement retry with backoff |
| 500 Internal Server Error | Retry with exponential backoff |

---

## Rate Limiting

The Combell API implements rate limiting. The provider should:

1. Read rate limit headers on each response:
   - `X-RateLimit-Limit`
   - `X-RateLimit-Usage`
   - `X-RateLimit-Remaining`
   - `X-RateLimit-Reset`

2. Implement automatic retry when receiving 429 status:
   - Use `Retry-After` header value
   - Implement exponential backoff

3. Consider implementing request throttling to prevent hitting limits

---

## Testing Strategy

### Unit Tests
- HMAC signature generation
- Request/response serialization
- Schema validation

### Acceptance Tests
- End-to-end CRUD operations
- Import functionality
- Error handling scenarios

### Test Environment Variables
```bash
export COMBELL_API_KEY="test-api-key"
export COMBELL_API_SECRET="test-api-secret"
export TF_ACC=1  # Enable acceptance tests
```

---

## Dependencies

```go
// go.mod
module github.com/veeblefetzer/terraform-provider-combell

go 1.21

require (
    github.com/hashicorp/terraform-plugin-framework v1.4.0
    github.com/hashicorp/terraform-plugin-go v0.19.0
    github.com/hashicorp/terraform-plugin-testing v1.5.1
)
```

---

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| 0.1.0 | TBD | Initial release with DNS and Domain support |

---

## References

- [Combell API Documentation](https://api.combell.com/v2/documentation)
- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- [Terraform Provider Development](https://developer.hashicorp.com/terraform/plugin)
