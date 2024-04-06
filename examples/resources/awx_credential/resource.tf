resource "awx_organization" "example" {
  name = "example"
}

resource "awx_credential" "example" {
  name               = "example"
  description        = "Example of ansible vault credential"
  organization_id    = awx_organization.example.id
  credential_type_id = 3 # ansible vault
  inputs = {
    vault_password = "admin"
    vault_id       = "password"
  }
}
