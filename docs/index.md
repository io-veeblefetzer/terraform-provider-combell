---
page_title: "Combell Provider"
subcategory: ""
description: |-
  Terraform provider for managing Combell resources via the Combell API v2.
---

# Combell Provider

The Combell provider enables Terraform to manage resources in your Combell hosting account via the [Combell API v2](https://api.combell.com/v2/documentation).

## Example Usage

```terraform
terraform {
  required_providers {
    combell = {
      source  = "io-veeblefetzer/combell"
      version = "~> 0.1"
    }
  }
}

# Configure the Combell provider
provider "combell" {
  # Credentials can be provided here or via environment variables
  # api_key    = "your-api-key"      # Or set COMBELL_API_KEY
  # api_secret = "your-api-secret"   # Or set COMBELL_API_SECRET
}
```

## Authentication

The provider uses HMAC authentication. You need to obtain your API key and secret from the [Combell control panel](https://my.combell.com).

### Credentials via Environment Variables

The recommended approach is to set credentials via environment variables:

```bash
export COMBELL_API_KEY="your-api-key"
export COMBELL_API_SECRET="your-api-secret"
```

### Credentials via Provider Configuration

Alternatively, you can set credentials directly in the provider configuration:

```terraform
provider "combell" {
  api_key    = "your-api-key"
  api_secret = "your-api-secret"
}
```

~> **Warning:** Hardcoding credentials in Terraform configuration is not recommended. Use environment variables or a secrets management system instead.

### IP Whitelisting

Access to the Combell API is restricted by IP address by default. You need to whitelist your IP address in the Combell control panel before you can use the API.

## Schema

### Optional

- `api_key` (String, Sensitive) - API key from Combell control panel. Can also be set via `COMBELL_API_KEY` environment variable.
- `api_secret` (String, Sensitive) - API secret from Combell control panel. Can also be set via `COMBELL_API_SECRET` environment variable.
- `base_url` (String) - API base URL. Defaults to `https://api.combell.com/v2`. Can also be set via `COMBELL_BASE_URL` environment variable.
