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

resource "awx_credential_machine" "example" {
  name            = "Example Machine Credential"
  description     = "Example Machine Credential"
  organization_id = data.awx_organization.default.id
  username        = "admin"
  password        = "password"
  become_method   = "sudo"
  become_username = "root"
  become_password = "password"
}

resource "awx_job_template" "example" {
  name           = "baseconfig"
  job_type       = "run"
  inventory_id   = awx_inventory.example.id
  project_id     = awx_project.example.id
  playbook       = "master-configure-system.yml"
  become_enabled = true
}

resource "awx_job_template_credential" "baseconfig" {
  job_template_id = awx_job_template.example.id
  credential_id   = awx_credential_machine.example.id
}
