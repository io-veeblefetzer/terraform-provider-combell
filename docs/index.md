---
page_title: "Combell Provider"
subcategory: ""
description: |-
  Terraform provider for managing Combell resources via the Combell API v2.
---

# Combell Provider

The Combell provider enables Terraform to manage resources in your Combell hosting account via the [Combell API v2](https://api.combell.com/v2/documentation).

## Features

- **DNS Records**: Create, read, update, and delete DNS records for your domains
- **Domains**: Query domain information and manage domain settings

## Example Usage

```hcl
terraform {
  required_providers {
    combell = {
      source  = "veeblefetzer/combell"
      version = "~> 0.1"
    }
  }
}

# Configure the Combell provider
provider "combell" {
  api_key    = var.combell_api_key
  api_secret = var.combell_api_secret
}
```

## Authentication

The Combell provider uses HMAC authentication. You need to provide your API key and secret, which can be obtained from the Combell control panel.

### Configuration

Credentials can be provided in the provider configuration or via environment variables:

```hcl
provider "combell" {
  api_key    = "your-api-key"
  api_secret = "your-api-secret"
}
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `COMBELL_API_KEY` | API key for authentication |
| `COMBELL_API_SECRET` | API secret for HMAC signature |
| `COMBELL_BASE_URL` | Optional API base URL override |

## Schema

### Optional

- `api_key` (String, Sensitive) - API key from Combell control panel. Can also be set via `COMBELL_API_KEY` environment variable.
- `api_secret` (String, Sensitive) - API secret from Combell control panel. Can also be set via `COMBELL_API_SECRET` environment variable.
- `base_url` (String) - API base URL. Defaults to `https://api.combell.com/v2`.
