---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "awx_job_template_role Data Source - terraform-provider-awx"
subcategory: ""
description: |-
  Data source for AWX Job Template Role
---

# awx_job_template_role (Data Source)

Data source for AWX Job Template Role

## Example Usage

```terraform
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
  job_template_id = awx_job_template.example.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `job_template_id` (Number) The ID of the job template

### Optional

- `id` (Number) The ID of the role
- `name` (String) The name of the role
