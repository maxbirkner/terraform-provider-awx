resource "awx_organization" "example" {
  name = "example"
}

resource "awx_credential_machine" "example" {
  name            = "Example Machine Credential"
  description     = "Example Machine Credential"
  organization_id = awx_organization.example.id
  username        = "admin"
  ssh_key_data    = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDZ6k1"
}


resource "awx_organization_galaxy_credential" "baseconfig" {
  organization_id = awx_organization.example.id
  credential_id   = awx_credential_machine.example.id
}
