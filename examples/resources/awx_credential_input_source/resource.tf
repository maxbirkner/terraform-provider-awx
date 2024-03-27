resource "awx_credential_input_source" "example" {
  description      = "example"
  input_field_name = "ssh_key_data"
  target           = 10
  source           = 10
  metadata = {
  }
}
