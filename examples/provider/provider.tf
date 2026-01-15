# Configure the Combell provider
provider "combell" {
  # API key from Combell control panel
  # Can also be set via COMBELL_API_KEY environment variable
  api_key = var.combell_api_key

  # API secret from Combell control panel
  # Can also be set via COMBELL_API_SECRET environment variable
  api_secret = var.combell_api_secret

  # Optional: Override the API base URL (defaults to https://api.combell.com/v2)
  # base_url = "https://api.combell.com/v2"
}

variable "combell_api_key" {
  type        = string
  description = "Combell API key"
  sensitive   = true
}

variable "combell_api_secret" {
  type        = string
  description = "Combell API secret"
  sensitive   = true
}
