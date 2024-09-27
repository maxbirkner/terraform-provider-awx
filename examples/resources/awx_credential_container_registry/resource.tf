resource "awx_organization" "example" {
  name = "example"
}

resource "awx_credential_container_registry" "example" {
  name            = "container-registry-credential"
  organization_id = awx_organization.example.id
  description     = "This is a container registry credential"
  host            = "container-registry.example.com"
  username        = "username"
  password        = "password"
  verify_ssl      = true
}
