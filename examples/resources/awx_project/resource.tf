resource "awx_organization" "example" {
  name = "example"
}

resource "awx_project" "example" {
  name                 = "base-service-configuration"
  scm_type             = "git"
  scm_url              = "https://github.com/nolte/ansible_playbook-baseline-online-server"
  scm_branch           = "feature/centos8-v2"
  scm_update_on_launch = true
  organization_id      = awx_organization.example.id
}
