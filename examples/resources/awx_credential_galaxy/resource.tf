resource "awx_organization" "example" {
  name = "example"
}

resource "awx_credential_galaxy" "example" {
  name            = "example"
  organization_id = awx_organization.example.id
  url             = "example"
  token           = "example"
  auth_url        = "example"
}
