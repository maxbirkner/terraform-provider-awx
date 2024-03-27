resource "awx_organization" "example" {
  name = "example"
}

resource "awx_credential_google_compute_engine" "example" {
  name            = "awx-gce-credential"
  organization_id = awx_organization.example.id
  description     = "This is a GCE credential"
  ssh_key_data    = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDQ6"
  username        = "admin"
  project         = "my-project"
}
