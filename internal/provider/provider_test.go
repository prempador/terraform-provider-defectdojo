// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Defectdojo client is properly configured.
	// It is also possible to use the DEFECTDOJO_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
	provider "defectdojo" {}
	`
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"defectdojo": providerserver.NewProtocol6WithError(New("test")()),
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func TestAccPreCheck(t *testing.T) {
	testDefectdojoHost(t)
	if v := os.Getenv("DEFECTDOJO_TOKEN"); v == "" {
		testDefectdojoUsername(t)
		testDefectdojoPassword(t)
	}
}

func testDefectdojoUsername(t *testing.T) {
	if v := os.Getenv("DEFECTDOJO_USERNAME"); v == "" {
		t.Fatal("DEFECTDOJO_USERNAME must be set for this acceptance test")
	}
}

func testDefectdojoPassword(t *testing.T) {
	if v := os.Getenv("DEFECTDOJO_PASSWORD"); v == "" {
		t.Fatal("DEFECTDOJO_PASSWORD must be set for this acceptance test")
	}
}

func testDefectdojoHost(t *testing.T) {
	if v := os.Getenv("DEFECTDOJO_HOST"); v == "" {
		t.Fatal("DEFECTDOJO_HOST must be set for this acceptance test")
	}
}
