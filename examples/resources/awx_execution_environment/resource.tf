resource "awx_execution_environment" "example" {
  name         = "Example"
  image        = "quay.io/ansible/awx-ee:1.0.0"
  credential   = awx_credential.example.id
  description  = "Example Execution Environment"
  organization = awx_organization.example.id
  pull         = "always"
}
