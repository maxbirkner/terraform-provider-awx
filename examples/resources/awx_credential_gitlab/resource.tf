resource "awx_organization" "example" {
  name = "example"
}

resource "awx_credential_gitlab" "example" {
  organization_id = awx_organization.example.id
  name            = "awx-scm-credential"
  description     = "test"
  token           = "My_TOKEN"
}
