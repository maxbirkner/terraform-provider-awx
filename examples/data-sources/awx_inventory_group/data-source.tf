data "awx_inventory_group" "example" {
  inventory_id = data.awx_inventory.example.id
  name         = "example"
}
