// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProductTypesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: providerConfig + `
				resource "defectdojo_product_type" "test1" {
					name = "ProductType1"
				}

				resource "defectdojo_product_type" "test2" {
					name = "ProductType2"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify product types were created
					resource.TestCheckResourceAttr("defectdojo_product_type.test1", "name", "ProductType1"),
					resource.TestCheckResourceAttr("defectdojo_product_type.test2", "name", "ProductType2"),
				),
			},
			{
				Config: providerConfig + `data "defectdojo_product_types" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of product types returned
					resource.TestCheckResourceAttr("data.defectdojo_product_types.test", "product_types.#", "3"), // need to check for 3 because the default product type is created
				),
			},
		},
	})
}
