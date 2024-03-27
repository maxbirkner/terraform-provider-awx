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


resource "awx_inventory_source" "example" {
  depends_on = [
    awx_project.example,
    awx_inventory.example
  ]
  inventory_id      = awx_inventory.example.id
  name              = "Example Inventory Source"
  description       = "Inventory source. [Managed by Terraform]"
  source_project_id = awx_project.example.id
  overwrite         = true
  source            = "scm"
  verbosity         = 1
  credential_id     = data.awx_credential.example.id
  source_path       = "inventories/example.yml"
}
