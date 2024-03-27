resource "awx_organization" "example" {
  name = "example"
}

resource "awx_credential_machine" "example" {
  name            = "example"
  description     = "Example Machine Credential"
  organization_id = awx_organization.example.id
  credential_type = "Machine"
  inputs = {
    username = "admin"
    password = "password"
  }
}
