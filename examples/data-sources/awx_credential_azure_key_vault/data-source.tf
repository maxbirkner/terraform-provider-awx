data "awx_credential_azure_key_vault" "example" {
  credential_id = 10
}

output "credential_name" {
  value = data.awx_credential_azure_key_vault.example.name
}
