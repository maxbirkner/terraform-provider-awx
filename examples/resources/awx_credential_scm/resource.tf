resource "awx_organization" "example" {
  name        = "Example Organization"
  description = "Example Organization"
}

resource "awx_credential_scm" "example" {
  name            = "Example Machine Credential"
  description     = "Example Machine Credential"
  organization_id = awx_organization.example.id
  username        = "admin"
  password        = "password"
}
