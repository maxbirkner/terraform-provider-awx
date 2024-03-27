resource "awx_host" "example" {
  name         = "some-host-node-1"
  description  = "Some host node 1"
  inventory_id = data.awx_inventory.default.id
  group_ids = [
    data.awx_inventory_group.default.id,
    data.awx_inventory_group.example.id,
  ]
  enabled   = true
  variables = <<YAML

---
ansible_host: 192.168.178.29
YAML
}
