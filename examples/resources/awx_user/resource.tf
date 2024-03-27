data "awx_organization" "example" {
  name = "Default"
}

data "awx_organization_role" "example" {
  name            = "Read"
  organization_id = awx_organization.example.id
}

resource "awx_user" "example" {
  username = "my_user"
  password = "my_password"
  role_entitlement {
    role_id = data.awx_organization_role.example.id
  }
}
