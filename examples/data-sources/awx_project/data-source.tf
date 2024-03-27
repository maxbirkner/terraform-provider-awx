data "awx_project" "my_project" {
  name = "my_project"
}

data "awx_project" "my_other_project" {
  id = 10
}

