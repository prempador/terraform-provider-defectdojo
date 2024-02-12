// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProductResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "defectdojo_product_type" "test_product_type" {
					name             = "Test Product Type"
					description      = "This is the description of the Test Product Type"
					critical_product = true
					key_product      = true
				}
				  
				resource "defectdojo_product" "test" {
					name             = "Test Product"
					description      = "This is the description of the Test Product"
					prod_type        = defectdojo_product_type.test_product_type.id
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify set fields
					resource.TestCheckResourceAttr("defectdojo_product.test", "name", "Test Product"),
					resource.TestCheckResourceAttr("defectdojo_product.test", "description", "This is the description of the Test Product"),
					resource.TestCheckResourceAttrPair("defectdojo_product.test", "prod_type", "defectdojo_product_type.test_product_type", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_product.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
