/a# Terraform Provider for Combell

[![Go](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Terraform](https://img.shields.io/badge/Terraform-1.0+-7B42BC?style=flat&logo=terraform)](https://www.terraform.io/)
[![License: MPL-2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

A Terraform provider for managing resources in [Combell](https://www.combell.com/) hosting platform via the [Combell API v2](https://api.combell.com/v2/documentation).

## Features

- 🔐 **HMAC Authentication** - Secure API access using HMAC-SHA256 signatures
- 🌐 **DNS Record Management** - Full CRUD operations for DNS records
- 📥 **Import Support** - Import existing resources into Terraform state
- 📚 **Documentation** - Comprehensive docs with examples

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.18 (for development)
- Combell API credentials ([obtain from control panel](https://my.combell.com))

## Installation

### From Terraform Registry (Coming Soon)

```hcl
terraform {
  required_providers {
    combell = {
      source  = "io-veeblefetzer/combell"
      version = "~> 0.1"
    }
  }
}
```

### Manual Installation

1. Build the provider:
   ```bash
   go build -o terraform-provider-combell
   ```

2. Install to local plugins directory:
   ```bash
   mkdir -p ~/.terraform.d/plugins/io-veeblefetzer/combell/0.1.0/darwin_amd64
   mv terraform-provider-combell ~/.terraform.d/plugins/io-veeblefetzer/combell/0.1.0/darwin_amd64/
   ```

## Authentication

The provider uses HMAC authentication. Obtain your API key and secret from the [Combell control panel](https://my.combell.com).

### Environment Variables (Recommended)

```bash
export COMBELL_API_KEY="your-api-key"
export COMBELL_API_SECRET="your-api-secret"
```

### Provider Configuration

```hcl
provider "combell" {
  api_key    = "your-api-key"    # Or use COMBELL_API_KEY
  api_secret = "your-api-secret" # Or use COMBELL_API_SECRET
}
```

> ⚠️ **Warning**: Avoid hardcoding credentials in configuration files. Use environment variables or a secrets management system.

### IP Whitelisting

The Combell API restricts access by IP address. Ensure your IP is whitelisted in the Combell control panel.

## Quick Start

### Managing DNS Records

```hcl
# Configure the provider
provider "combell" {}

# Create an A record
resource "combell_dns_record" "www" {
  domain_name = "example.com"
  type        = "A"
  record_name = "www"
  content     = "192.168.1.1"
  ttl         = 3600
}

# Create an MX record
resource "combell_dns_record" "mail" {
  domain_name = "example.com"
  type        = "MX"
  record_name = ""
  content     = "mail.example.com"
  priority    = 10
  ttl         = 3600
}

# Create a TXT record for SPF
resource "combell_dns_record" "spf" {
  domain_name = "example.com"
  type        = "TXT"
  record_name = ""
  content     = "v=spf1 include:_spf.google.com ~all"
  ttl         = 3600
}
```

### Reading Existing DNS Records

```hcl
data "combell_dns_record" "existing" {
  domain_name = "example.com"
  id          = "12345"
}

output "record_type" {
  value = data.combell_dns_record.existing.type
}
```

### Importing Existing Resources

```bash
terraform import combell_dns_record.www "example.com/12345"
```

## Resources

| Resource | Description |
|----------|-------------|
| `combell_dns_record` | Manages DNS records (A, AAAA, CAA, CNAME, MX, TXT, SRV, ALIAS, TLSA) |

## Data Sources

| Data Source | Description |
|-------------|-------------|
| `combell_dns_record` | Retrieves an existing DNS record by ID |

## Documentation

- [Provider Documentation](docs/index.md)
- [DNS Record Resource](docs/resources/dns_record.md)
- [DNS Record Data Source](docs/data-sources/dns_record.md)

## Development

### Building

```bash
go build -v ./...
```

### Testing

```bash
# Unit tests
go test ./...

# Acceptance tests (requires API credentials)
export COMBELL_API_KEY="your-key"
export COMBELL_API_SECRET="your-secret"
export TF_ACC=1
go test ./... -v
```

### Project Structure

```
├── internal/
│   ├── client/          # HTTP client with HMAC authentication
│   ├── provider/        # Terraform provider configuration
│   └── resources/       # Resource and data source implementations
│       ├── dns_record/  # DNS record resource/data source
│       └── domain/      # Domain resource/data source (planned)
├── docs/                # Provider documentation
├── examples/            # Example configurations
└── resources/           # API specifications
```

## Contributing

Contributions are welcome! Please follow the guidelines in [.clinerules](.clinerules).

### Commit Messages

All commits must follow [Conventional Commits](https://www.conventionalcommits.org/) and reference a GitHub issue:

```
feat(dns): implement DNS record resource (#12)
fix(client): handle rate limiting (#45)
docs(readme): add installation instructions (#23)
```

## Roadmap

- [x] Phase 1: Project setup and provider skeleton
- [x] Phase 2: HMAC authentication client
- [x] Phase 3: DNS record resource and data source
- [ ] Phase 4: Domain data source
- [ ] Phase 5: CI/CD pipeline
- [ ] Phase 6: Terraform Registry publishing

## License

[Mozilla Public License 2.0](LICENCE)

## Support

- [Report Issues](https://github.com/io-veeblefetzer/terraform-provider-combell/issues)
- [Combell API Documentation](https://api.combell.com/v2/documentation)
