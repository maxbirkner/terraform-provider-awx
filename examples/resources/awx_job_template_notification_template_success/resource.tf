data "awx_organization" "example" {
  name = "Default"
}

resource "awx_project" "example" {
  name            = "example-ansible-main"
  organization_id = data.awx_organization.example.id
  scm_type        = "git"
  scm_url         = "git@github.com/josh-silvas/example-ansible.git"
  scm_branch      = "main"
}

data "awx_inventory" "example" {
  name            = "private_services"
  organization_id = data.awx_organization.example.id
}

resource "awx_job_template" "example" {
  name           = "baseconfig"
  job_type       = "run"
  inventory_id   = data.awx_inventory.example.id
  project_id     = awx_project.example.id
  playbook       = "master-configure-system.yml"
  become_enabled = true
}

resource "awx_notification_template" "example" {
  name = "notification_template-test"
}

resource "awx_job_template_notification_template_success" "example" {
  job_template_id          = awx_job_template.example.id
  notification_template_id = awx_notification_template.example.id
}
