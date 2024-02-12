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

resource "defectdojo_engagement" "test_engagement" {
  name        = "Test Engagement"
  description = "This is the description of the Test Engagement"
  product     = defectdojo_product.test_product.id
  start_date  = "2024-01-01"
  end_date    = "2024-01-31"
}