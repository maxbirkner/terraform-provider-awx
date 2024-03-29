terraform {
  required_providers {
    awx = {
      source  = "josh-silvas/awx"
      version = "1.0.6"
    }
  }
}

// Example configuration for the AWX provider using a username and password
provider "awx_with_username_password" {
  hostname = "https://awx.example.com"
  username = "admin"
  password = "password"
}

// Example configuration for the AWX provider using a token
provider "awx_with_token" {
  hostname = "https://awx.example.com"
  token    = "token"
}
