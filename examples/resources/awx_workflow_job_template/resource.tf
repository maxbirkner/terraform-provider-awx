data "awx_organization" "example" {
  name = "Default"
}

resource "awx_inventory" "example" {
  name            = "Example Inventory"
  organization_id = data.awx_organization.example.id
}

resource "awx_workflow_job_template" "example" {
  name            = "workflow-job"
  organization_id = data.awx_organization.example.id
  inventory_id    = awx_inventory.example.id
}
