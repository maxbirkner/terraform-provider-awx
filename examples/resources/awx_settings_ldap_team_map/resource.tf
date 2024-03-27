data "awx_organization" "example" {
  name = "Default"
}

resource "awx_team" "example" {
  name            = "Admins"
  organization_id = data.awx_organization.example.id
}

resource "awx_settings_ldap_team_map" "example" {
  name         = awx_team.example.name
  users        = ["CN=MyTeam,OU=Groups,DC=example,DC=com"]
  organization = data.awx_organization.example.name
  remove       = true
}
