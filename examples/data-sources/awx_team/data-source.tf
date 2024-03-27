data "awx_team" "example" {
  name = "My Team"
}

data "awx_team" "default" {
  id = 5
}
