data "awx_organization" "example" {
  name = "Default"
}


data "awx_inventory" "example" {
  name = "My Inventory"
}

data "awx_inventory_role" "example" {
  name         = "Admin"
  inventory_id = data.awx_inventory.example.id
}

data "awx_project" "example" {
  name = "My Project"
}

data "awx_project_role" "example" {
  name       = "Admin"
  project_id = data.awx_project.example.id
}

resource "awx_team" "example" {
  name            = "admins-team"
  organization_id = data.awx_organization.example.id

  role_entitlement {
    role_id = data.awx_inventory_role.example.id
  }
  role_entitlement {
    role_id = data.awx_project_role.example.id
  }
}
