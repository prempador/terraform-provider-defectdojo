package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDojoGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
			resource "defectdojo_dojo_group" "test" {
				name = "DojoGroup"
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify set fields
					resource.TestCheckResourceAttr("defectdojo_dojo_group.test", "name", "DojoGroup"),
					// Verify default fields
					resource.TestCheckResourceAttr("defectdojo_dojo_group.test", "description", ""),
					resource.TestCheckResourceAttr("defectdojo_dojo_group.test", "social_provider", ""),
					// Verify default configuration_permissions length
					resource.TestCheckResourceAttr("defectdojo_dojo_group.test", "configuration_permissions.#", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_dojo_group.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
			resource "defectdojo_dojo_group" "test" {
				name 			= "UpdatedDojoGroup"
				description 	= "UpdatedDescription"
				social_provider = "AzureAD" 			# currently only supports AzureAD
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify updated fields
					resource.TestCheckResourceAttr("defectdojo_dojo_group.test", "name", "UpdatedDojoGroup"),
					resource.TestCheckResourceAttr("defectdojo_dojo_group.test", "description", "UpdatedDescription"),
					resource.TestCheckResourceAttr("defectdojo_dojo_group.test", "social_provider", "AzureAD"),
					// can't update configuration_permissions for now
					resource.TestCheckResourceAttr("defectdojo_dojo_group.test", "configuration_permissions.#", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
