resource "awx_organization" "example" {
  name        = "Example Organization"
  description = "Example Organization"
}

resource "awx_credential_machine" "example_1" {
  name            = "Example Machine Credential"
  description     = "Example Machine Credential"
  organization_id = awx_organization.example.id
  username        = "admin"
  password        = "password"
  become_method   = "sudo"
  become_username = "root"
  become_password = "password"
}

resource "awx_credential_machine" "example_2" {
  name            = "Example 2 Machine Credential"
  description     = "Example 2 Machine Credential"
  organization_id = awx_organization.example.id
  username        = "admin"
  ssh_key_data    = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDZ6k1"
}
