// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/prempador/go-defectdojo"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &UsersDataSource{}
	_ datasource.DataSourceWithConfigure = &UsersDataSource{}
)

func NewUsersDataSource() datasource.DataSource {
	return &UsersDataSource{}
}

// UsersDataSource defines the data source implementation.
type UsersDataSource struct {
	client *defectdojo.APIClient
}

// UsersDataSourceModel describes the data source data model.
type UsersDataSourceModel struct {
	Users []userModel `tfsdk:"users"`
}

// userModel describes the data source data model.
type userModel struct {
	ID                       types.Int64   `tfsdk:"id"`
	Username                 types.String  `tfsdk:"username"`
	FirstName                types.String  `tfsdk:"first_name"`
	LastName                 types.String  `tfsdk:"last_name"`
	Email                    types.String  `tfsdk:"email"`
	DateJoined               types.String  `tfsdk:"date_joined"`
	LastLogin                types.String  `tfsdk:"last_login"`
	IsActive                 types.Bool    `tfsdk:"is_active"`
	IsSuperuser              types.Bool    `tfsdk:"is_superuser"`
	ConfigurationPermissions []types.Int64 `tfsdk:"configuration_permissions"`
}

// Metadata returns the data source type name.
func (d *UsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the data source.
func (d *UsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"users": schema.ListNestedAttribute{
				Description: "List of users",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "The unique identifier for the user",
							Computed:    true,
						},
						"username": schema.StringAttribute{
							Description: "The username of the user",
							Computed:    true,
						},
						"first_name": schema.StringAttribute{
							Description: "The first name of the user",
							Computed:    true,
						},
						"last_name": schema.StringAttribute{
							Description: "The last name of the user",
							Computed:    true,
						},
						"email": schema.StringAttribute{
							Description: "The email of the user",
							Computed:    true,
						},
						"date_joined": schema.StringAttribute{
							Description: "The date the user joined",
							Computed:    true,
						},
						"last_login": schema.StringAttribute{
							Description: "The last login date of the user",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Whether the user is active",
							Computed:    true,
						},
						"is_superuser": schema.BoolAttribute{
							Description: "Whether the user is a superuser",
							Computed:    true,
						},
						"configuration_permissions": schema.ListAttribute{
							Description: "Configuration permissions of the user",
							Computed:    true,
							ElementType: types.Int64Type,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *UsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*defectdojo.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *defectdojo.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *UsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state UsersDataSourceModel

	// Fetch data from the API
	users, res, err := d.client.UsersAPI.UsersList(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Users",
			"Could not red users, unexpected error: "+err.Error()+"\nDefectdojo responded with status: "+fmt.Sprintf("%v", res.Body),
		)
		return
	}

	// Map response body to model
	for _, user := range users.Results {
		userModel := userModel{
			ID:          types.Int64Value(int64(user.GetId())),
			Username:    types.StringValue(user.GetUsername()),
			FirstName:   types.StringValue(user.GetFirstName()),
			LastName:    types.StringValue(user.GetLastName()),
			Email:       types.StringValue(user.GetEmail()),
			DateJoined:  types.StringValue(user.GetDateJoined().String()),
			LastLogin:   types.StringValue(user.GetLastLogin().String()),
			IsActive:    types.BoolValue(user.GetIsActive()),
			IsSuperuser: types.BoolValue(user.GetIsSuperuser()),
		}

		for _, configurationPermissions := range user.GetConfigurationPermissions() {
			userModel.ConfigurationPermissions = append(userModel.ConfigurationPermissions, int32PointerToBasetypesInt64Value(configurationPermissions))
		}

		state.Users = append(state.Users, userModel)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
