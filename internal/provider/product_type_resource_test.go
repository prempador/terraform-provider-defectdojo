// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProductTypeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
			resource "defectdojo_product_type" "test" {
				name = "ProductType"
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify set fields
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "name", "ProductType"),
					// Verify default fields
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "description", ""),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "critical_product", "false"),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "key_product", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_product_type.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
			resource "defectdojo_product_type" "test" {
				name 				= "UpdatedProductType"
				description 		= "UpdatedDescription"
				critical_product 	= true
				key_product 		= true
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify updated fields
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "name", "UpdatedProductType"),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "description", "UpdatedDescription"),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "critical_product", "true"),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "key_product", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
