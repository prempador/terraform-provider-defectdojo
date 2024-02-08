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
				resource "defectdojo_product_type" "test" {
					name = "ProductType"
				}

				resource "defectdojo_product_type" "test2" {
					name = "ProductType2"
				}

				data "defectdojo_product_types" "test" {}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of product types returned
					resource.TestCheckResourceAttr("data.defectdojo_product_types.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("data.defectdojo_product_types.test", "product_types.0.id", "1"),
					resource.TestCheckResourceAttr("data.defectdojo_product_types.test", "product_types.0.name", "Research and Development"),
				),
			},
		},
	})
}
