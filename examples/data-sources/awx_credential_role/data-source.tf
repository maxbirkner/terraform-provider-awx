data "awx_credential_role" "example" {
  credential_id = 10
}

output "credential_role_id" {
  value = data.awx_credential_role.example.id
}
