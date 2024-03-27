data "awx_organization" "default" {
  name = "Default"
}

resource "awx_inventory" "example" {
  name            = "Example Inventory"
  organization_id = data.awx_organization.default.id
  variables       = <<YAML

---
system_supporters:
  - pi

YAML
}
