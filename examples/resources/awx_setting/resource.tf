resource "awx_setting" "social_auth_saml_technical_contact" {
  name  = "SOCIAL_AUTH_SAML_TECHNICAL_CONTACT"
  value = <<EOF
	  {
	    "givenName": "Myorg",
	    "emailAddress": "test@foo.com"
	  }
	  EOF
}

resource "awx_setting" "social_auth_saml_sp_entity_id" {
  name  = "SOCIAL_AUTH_SAML_SP_ENTITY_ID"
  value = "test"
}

resource "awx_setting" "schedule_max_jobs" {
  name  = "SCHEDULE_MAX_JOBS"
  value = 15
}

resource "awx_setting" "remote_host_headers" {
  name  = "REMOTE_HOST_HEADERS"
  value = <<EOF
	  [
	    "HTTP_X_FORWARDED_FOR",
	    "REMOTE_ADDR",
	    "REMOTE_HOST"
	  ]
	  EOF
}
