resource "awx_credential_type" "example" {
  name        = "Example Credential Type"
  description = "Example Credential Type"
  kind        = "cloud"
  inputs      = "{\"fields\": [{\"id\": \"password\", \"label\": \"Password\", \"type\": \"string\", \"secret\": true}]}"
  injectors   = "{\"extra_vars\": {\"password\": \"password\"}}"
}
