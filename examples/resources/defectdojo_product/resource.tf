resource "defectdojo_product_type" "test_product_type" {
  name             = "Test Product Type"
  description      = "This is the description of the Test Product Type"
  critical_product = true
  key_product      = true
}

resource "defectdojo_product" "test_product" {
  name        = "Test Product"
  description = "This is the description of the Test Product"
  prod_type   = defectdojo_product_type.test_product_type.id
}
