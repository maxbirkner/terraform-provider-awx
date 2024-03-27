data "awx_organization" "example" {
  name = "My Organization"
}


data "awx_organization_role" "example" {
  name            = "Admin"
  organization_id = awx_organization.example.id
}
