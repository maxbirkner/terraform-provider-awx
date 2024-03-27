data "awx_organization" "default" {
  name = "Default"
}

resource "awx_project" "example" {
  name            = "example-ansible-main"
  organization_id = data.awx_organization.default.id
  scm_type        = "git"
  scm_url         = "git@github.com/josh-silvas/example-ansible.git"
  scm_branch      = "main"
}

resource "awx_inventory" "example" {
  name            = "Example Inventory"
  organization_id = data.awx_organization.default.id
}

resource "awx_job_template" "example" {
  name           = "baseconfig"
  job_type       = "run"
  inventory_id   = awx_inventory.example.id
  project_id     = awx_project.example.id
  playbook       = "master-configure-system.yml"
  become_enabled = true
}
