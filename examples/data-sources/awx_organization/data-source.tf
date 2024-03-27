data "awx_organization" "my_organization" {
  name = "My Organization"
}

data "awx_organization" "my_other_organization" {
  id = 10
}
