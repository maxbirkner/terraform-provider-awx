resource "awx_organization" "example" {
  name = "example"
}

resource "awx_credential_vault" "example" {
  name            = "vault-credential"
  organization_id = awx_organization.example.id
  description     = "This is a vault credential"
  vault_password  = "password"
}
