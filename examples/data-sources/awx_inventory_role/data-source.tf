data "awx_inventory_role" "example" {
  inventory_id = data.awx_inventory.example.id
}
