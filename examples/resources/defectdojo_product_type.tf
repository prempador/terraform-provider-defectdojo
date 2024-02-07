resource "defectdojo_product_type" "test_product_type" {
  name             = "Test Product Type"
  description      = "This is the description of the Test Product Type"
  critical_product = true
  key_product      = true
}