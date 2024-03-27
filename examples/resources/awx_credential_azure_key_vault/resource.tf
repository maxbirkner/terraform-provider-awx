resource "awx_organization" "example" {
  name = "example"
}

resource "awx_credential_azure_key_vault" "example" {
  name            = "example"
  organization_id = awx_organization.example.id
  secret          = "example"
  client          = "example"
  tenant          = "example"
  url             = "example"
}
