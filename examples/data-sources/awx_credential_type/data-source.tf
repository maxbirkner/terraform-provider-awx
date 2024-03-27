data "awx_credential_type" "example" {
  credential_id = 10
}

output "credential_type_id" {
  value = data.awx_credential_type.example.id
}
