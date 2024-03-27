data "awx_project" "example" {
  name = "my_project"
}

data "awx_project_role" "example" {
  name       = "Admin"
  project_id = awx_project.example.id
}

