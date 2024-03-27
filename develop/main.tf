terraform {
  required_providers {
    awx = {
      source = "registry.terraform.io/josh-silvas/awx"
    }
  }
}

locals {
  url   = "awx.example.com"
  token = "admin"
}

// Configure the AWX provider. This example relies on AWX_HOSTNAME and AWX_TOKEN to exist
// in the environment. If not, you will need to specify them here.
provider "awx" {
#    hostname = local.url
#    token    = local.token
#    insecure = true
}

data "awx_organization" "my_org" {
    name = "my_org"
}

resource "awx_credential_scm" "my_credential" {
    name            = "SCM Service Account"
    username        = "my_username"
    description     = "Service Account for git access to the repo. [Managed by Terraform]"
    organization_id = data.awx_organization.my_org.id
    ssh_key_data    = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDQq3"
}

resource "awx_project" "ansible_main" {
    name                 = "Ansible Repository"
    description          = "This project is associated to the project main branch. [Managed by Terraform]"
    scm_type             = "git"
    scm_url              = "git@github.com:org/repo.git"
    scm_branch           = "main"
    scm_update_on_launch = true
    organization_id      = data.awx_organization.my_org.id
    scm_credential_id    = awx_credential_scm.my_credential.id
}

data "awx_project_role" "cns_ansible_main_admins" {
    name       = "Admin"
    project_id = awx_project.ansible_main.id
}
