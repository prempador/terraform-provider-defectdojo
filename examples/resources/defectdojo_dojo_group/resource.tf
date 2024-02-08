resource "defectdojo_dojo_group" "test_dojo_group" {
  name            = "Dojo Group"
  description     = "Dojo Group Description"
  social_provider = "AzureAD" # currently only supports AzureAD
}
