terraform {
  required_providers {
    awx = {
      source  = "github.com/josh-silvas/awx"
      version = "0.1"
    }
  }
}

provider "awx" {
  hostname = "https://awx.example.com"
  username = "admin"
  password = "password"
  token    = "token"
}
