// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEngagementResource(t *testing.T) {
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
				  
				resource "defectdojo_product" "test_product" {
					name        = "Test Product"
					description = "This is the description of the Test Product"
					prod_type   = defectdojo_product_type.test_product_type.id
				}
				  
				resource "defectdojo_engagement" "test" {
					name          = "Test Engagement"
					description   = "This is the description of the Test Engagement"
					product       = defectdojo_product.test_product.id 
					target_start  = "2024-01-01"
					target_end    = "2024-01-31"
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify set fields
					resource.TestCheckResourceAttr("defectdojo_engagement.test", "name", "Test Engagement"),
					resource.TestCheckResourceAttr("defectdojo_engagement.test", "description", "This is the description of the Test Engagement"),
					resource.TestCheckResourceAttr("defectdojo_engagement.test", "target_start", "2024-01-01"),
					resource.TestCheckResourceAttr("defectdojo_engagement.test", "target_end", "2024-01-31"),
					resource.TestCheckResourceAttrPair("defectdojo_engagement.test", "product", "defectdojo_product.test_product", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_engagement.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
