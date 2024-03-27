resource "awx_job_template" "example" {
  name           = "baseconfig"
  job_type       = "run"
  inventory_id   = data.awx_inventory.default.id
  project_id     = awx_project.base_service_config.id
  playbook       = "master-configure-system.yml"
  become_enabled = true
}

data "awx_job_template_role" "example" {
  name            = "Admin"
  job_template_id = data.awx_job_template.example.id
}
