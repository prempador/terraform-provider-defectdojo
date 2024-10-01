// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDojoGroupMemberResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "defectdojo_dojo_group" "test_group" {
					name = "DojoGroupMemberTestGroup"
				}

				resource "defectdojo_user" "test_user" {
					username = "DojoGroupMemberTestUser"
					email	 = "email@email.com"
					password = "veryHardPassword1234!"
				}

				resource "defectdojo_dojo_group_member" "test" {
					group = defectdojo_dojo_group.test_group.id
					user  = defectdojo_user.test_user.id
					role  = "Writer"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify set fields
					resource.TestCheckResourceAttrPair("defectdojo_dojo_group_member.test", "group", "defectdojo_dojo_group.test_group", "id"),
					resource.TestCheckResourceAttrPair("defectdojo_dojo_group_member.test", "user", "defectdojo_user.test_user", "id"),
					resource.TestCheckResourceAttr("defectdojo_dojo_group_member.test", "role", "Writer"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_dojo_group_member.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Currently unable to test Update and Read as the Defectdojo API deletes the group member when the group or user is deleted

			// Delete testing automatically occurs in TestCase
		},
	})
}
