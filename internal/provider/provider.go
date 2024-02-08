// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/prempador/go-defectdojo"
)

// Ensure DefectDojoProvider satisfies various provider interfaces.
var _ provider.Provider = &DefectdojoProvider{}

// DefectdojoProvider defines the provider implementation.
type DefectdojoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// DefectdojoProviderModel describes the provider data model.
type DefectdojoProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Token    types.String `tfsdk:"token"`
}

func (p *DefectdojoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "defectdojo"
	resp.Version = p.version
}

func (p *DefectdojoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "The host of the defectdojo instance",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username of the defectdojo user (required if token is not set)",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password of the defectdojo user (required if token is not set)",
				Optional:            true,
				Sensitive:           true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "The token of the defectdojo user (required if username and password are not set)",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *DefectdojoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data DefectdojoProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("DEFECTDOJO_HOST")
	username := os.Getenv("DEFECTDOJO_USERNAME")
	password := os.Getenv("DEFECTDOJO_PASSWORD")
	token := os.Getenv("DEFECTDOJO_TOKEN")

	if !data.Host.IsNull() {
		host = data.Host.ValueString()
	}

	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	}

	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}

	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Defectdojo API Host",
			"The provider cannot create the Defectdojo API client as there is a missing or empty value for the Defectdojo API host. "+
				"Set the host value in the configuration or use the DEFECTDOJO_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	// if token is empty we require username and password to be set to fetch a token
	if token == "" {
		if username == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("username"),
				"Missing Defectdojo API Username",
				"The provider cannot create the Defectdojo API client as there is a missing or empty value for the Defectdojo API username. "+
					"Set the username value in the configuration or use the DEFECTDOJO_USERNAME environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}

		if password == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("password"),
				"Missing Defectdojo API Password",
				"The provider cannot create the Defectdojo API client as there is a missing or empty value for the Defectdojo API password. "+
					"Set the password value in the configuration or use the DEFECTDOJO_PASSWORD environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	transport := cleanhttp.DefaultPooledTransport()
	httpClient := retryablehttp.NewClient()
	httpClient.HTTPClient.Transport = transport
	httpClient.Logger = nil // disable logging

	parsedHost, err := url.Parse(host)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse Defectdojo API Host", "Failed to parse Defectdojo API Host: "+err.Error())
		return
	}

	cfg := defectdojo.NewConfiguration()
	cfg.Host = parsedHost.Host
	cfg.Scheme = parsedHost.Scheme
	cfg.HTTPClient = httpClient.StandardClient()
	cfg.UserAgent = "terraform-provider-defectdojo/" + p.version

	// if token is empty we are fetching a new one with the username and password
	if token == "" {
		// we have to go oldschool here because the openapi definition for the api token endpoint is not working
		body := []byte(fmt.Sprintf(`{
			"username": "%s",
			"password": "%s"
		}`, username, password))

		// Create a HTTP post request
		r, err := http.NewRequest("POST", host+"/api/v2/api-token-auth/", bytes.NewBuffer(body))
		if err != nil {
			resp.Diagnostics.AddError("Failed to create POST Request to fetch API Token", "Failed to create POST Request to fetch API Token: "+err.Error())

			return
		}

		r.Header.Add("Content-Type", "application/json")
		r.Header.Add("user-agent", "terraform-provider-defectdojo/"+p.version)

		res, err := httpClient.StandardClient().Do(r)
		if err != nil {
			resp.Diagnostics.AddError("Failed to create POST Request to fetch API Token", "Failed to create POST Request to fetch API Token: "+err.Error())

			return
		}

		defer res.Body.Close()

		response := struct {
			Token string `json:"token"`
		}{}

		derr := json.NewDecoder(res.Body).Decode(&response)
		if derr != nil {
			resp.Diagnostics.AddError("Failed to decode Defectdojo API Response", "Failed to decode Defectdojo API Response: "+err.Error())

			return
		}

		// this would be how it would work if the definition was correct

		// c := defectdojo.NewAPIClient(cfg)

		// ddtoken, res, err := c.ApiTokenAuthAPI.ApiTokenAuthCreate(context.Background()).Username(username).Password(password).Execute()
		// if err != nil {
		// 	resp.Diagnostics.AddError("Failed to authenticate with Defectdojo API", "Failed to authenticate with Defectdojo API: "+err.Error()+fmt.Sprintf("%+v %s %s", res, username, password))
		// 	return
		// }

		token = response.Token
	}

	cfg.AddDefaultHeader("Authorization", "Token "+token)
	client := defectdojo.NewAPIClient(cfg)

	// Make the defectdojo client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *DefectdojoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDojoGroupResource,
		NewDojoGroupMemberResource,
		NewProductTypeResource,
		NewUserResource,
	}
}

func (p *DefectdojoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProductTypesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DefectdojoProvider{
			version: version,
		}
	}
}
