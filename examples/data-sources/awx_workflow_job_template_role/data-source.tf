resource "awx_workflow_job_template" "example" {
  name            = "workflow-job"
  organization_id = data.awx_organization.example.id
  inventory_id    = data.awx_inventory.default.id
}

data "awx_workflow_job_template_role" "example" {
  name                     = "Admin"
  workflow_job_template_id = awx_workflow_job_template.example.id
}
