// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
			resource "defectdojo_user" "test" {
				username = "User"
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify set fields
					resource.TestCheckResourceAttr("defectdojo_user.test", "username", "User"),
					// Verify default fields
					resource.TestCheckResourceAttr("defectdojo_user.test", "first_name", ""),
					resource.TestCheckResourceAttr("defectdojo_user.test", "last_name", ""),
					resource.TestCheckResourceAttr("defectdojo_user.test", "email", ""),
					resource.TestCheckResourceAttr("defectdojo_user.test", "is_active", "false"),
					resource.TestCheckResourceAttr("defectdojo_user.test", "is_superuser", "false"),
					// Verify default configuration_permissions length
					resource.TestCheckResourceAttr("defectdojo_user.test", "configuration_permissions.#", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_user.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
			resource "defectdojo_user" "test" {
				username 			= "UpdatedUser"
				first_name 			= "UpdatedFirstName"
				last_name 			= "UpdatedLastName"
				email 				= "email@email.com"
				is_active 			= true
				is_superuser 		= true
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify updated fields
					resource.TestCheckResourceAttr("defectdojo_user.test", "username", "UpdatedUser"),
					resource.TestCheckResourceAttr("defectdojo_user.test", "first_name", "UpdatedFirstName"),
					resource.TestCheckResourceAttr("defectdojo_user.test", "last_name", "UpdatedLastName"),
					resource.TestCheckResourceAttr("defectdojo_user.test", "email", "email@email.com"),
					resource.TestCheckResourceAttr("defectdojo_user.test", "is_active", "true"),
					resource.TestCheckResourceAttr("defectdojo_user.test", "is_superuser", "true"),
					// can't update configuration_permissions for now
					resource.TestCheckResourceAttr("defectdojo_user.test", "configuration_permissions.#", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
