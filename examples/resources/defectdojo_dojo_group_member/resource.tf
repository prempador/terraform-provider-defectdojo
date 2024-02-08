resource "defectdojo_dojo_group" "test_group" {
  name = "DojoGroup"
}

resource "defectdojo_user" "test_user" {
  username = "TestUser"
}

resource "defectdojo_dojo_group_member" "test_group_member" {
  group = defectdojo_dojo_group.test_group.id
  user  = defectdojo_user.test_user.id
  role  = "Writer"
}