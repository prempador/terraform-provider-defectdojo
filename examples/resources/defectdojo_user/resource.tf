resource "defectdojo_user" "test_user" {
  username     = "Username"
  first_name   = "First"
  last_name    = "Last"
  email        = "terraform@provider.com"
  is_active    = true
  is_superuser = false
}
