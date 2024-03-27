data "awx_organization" "default" {
  name = "Default"
}

resource "awx_inventory" "example" {
  name            = "Example Inventory"
  organization_id = data.awx_organization.default.id
}


resource "awx_inventory_group" "example" {
  name         = "Example"
  description  = "Example Inventory Group"
  inventory_id = awx_inventory.example.id
  variables    = <<EOF
{
    "example": "example"
}
EOF
}
