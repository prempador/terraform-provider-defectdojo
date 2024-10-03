// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUsersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
				resource "defectdojo_user" "test1" {
					username = "User1"
					email	 = "email1@email.com"
					password = "veryHardPassword1234!"
				}

				resource "defectdojo_user" "test2" {
					username = "User2"
					email	 = "email2@email.com"
					password = "veryHardPassword1234!"
				}

                data "defectdojo_users" "test" {}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of users returned
					resource.TestCheckResourceAttr("data.defectdojo_users.test", "users.#", "3"), // need to check for 3 because the default user is created
					resource.TestCheckResourceAttr("data.defectdojo_users.test", "users.1.username", "User1"),
					resource.TestCheckResourceAttr("data.defectdojo_users.test", "users.1.email", "email1@email.com"),
					resource.TestCheckResourceAttr("data.defectdojo_users.test", "users.2.username", "User2"),
					resource.TestCheckResourceAttr("data.defectdojo_users.test", "users.2.email", "email2@email.com"),
				),
			},
		},
	})
}
