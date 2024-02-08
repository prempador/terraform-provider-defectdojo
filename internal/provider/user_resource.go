// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/prempador/go-defectdojo"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

// NewUserResource is a helper function to simplify the provider implementation.
func NewUserResource() resource.Resource {
	return &userResource{}
}

// userResource is the data source implementation.
type userResource struct {
	client *defectdojo.APIClient
}

type userResourceModel struct {
	ID                       types.Int64  `tfsdk:"id"`
	Username                 types.String `tfsdk:"username"`
	FirstName                types.String `tfsdk:"first_name"`
	LastName                 types.String `tfsdk:"last_name"`
	Email                    types.String `tfsdk:"email"`
	IsActive                 types.Bool   `tfsdk:"is_active"`
	IsSuperUser              types.Bool   `tfsdk:"is_superuser"`
	Password                 types.String `tfsdk:"password"`
	ConfigurationPermissions types.List   `tfsdk:"configuration_permissions"`
}

// Metadata returns the data source type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the data source.
func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The unique identifier of the user",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Description: "The username of the user",
				Required:    true,
			},
			"first_name": schema.StringAttribute{
				Description: "The first name of the user",
				Computed:    true,
				Optional:    true,
			},
			"last_name": schema.StringAttribute{
				Description: "The last name of the user",
				Computed:    true,
				Optional:    true,
			},
			"email": schema.StringAttribute{
				Description: "The email of the user",
				Computed:    true,
				Optional:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "The active status of the user",
				Computed:    true,
				Optional:    true,
			},
			"is_superuser": schema.BoolAttribute{
				Description: "The superuser status of the user",
				Computed:    true,
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password of the user",
				Optional:    true,
				Sensitive:   true,
			},
			"configuration_permissions": schema.ListAttribute{
				ElementType: types.Int64Type,
				Description: "The configuration permissions of the user",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*defectdojo.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *defectdojo.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	configurations := make([]*int32, 0)
	diags = plan.ConfigurationPermissions.ElementsAs(ctx, &configurations, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate request from plan
	userRequest := defectdojo.UserRequest{
		Username:                 plan.Username.ValueString(),
		FirstName:                plan.FirstName.ValueStringPointer(),
		LastName:                 plan.LastName.ValueStringPointer(),
		Email:                    plan.Email.ValueStringPointer(),
		IsActive:                 plan.IsActive.ValueBoolPointer(),
		IsSuperuser:              plan.IsSuperUser.ValueBoolPointer(),
		Password:                 plan.Password.ValueStringPointer(),
		ConfigurationPermissions: configurations,
	}

	// Create new product type
	user, res, err := r.client.UsersAPI.UsersCreate(ctx).UserRequest(userRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Defectdojo User",
			"Could not create user, unexpected error: "+err.Error()+"\nDefectdojo responded with status:"+res.Status,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(int64(user.GetId()))
	plan.Username = types.StringValue(user.GetUsername())
	plan.FirstName = types.StringValue(user.GetFirstName())
	plan.LastName = types.StringValue(user.GetLastName())
	plan.Email = types.StringValue(user.GetEmail())
	plan.IsActive = types.BoolValue(user.GetIsActive())
	plan.IsSuperUser = types.BoolValue(user.GetIsSuperuser())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("configuration_permissions"), user.ConfigurationPermissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed product type value from Defectdojo
	user, res, err := r.client.UsersAPI.UsersRetrieve(ctx, int32(state.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Defectdojo User",
			"Could not read user with ID "+state.ID.String()+": "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}

	// Overwrite state with refreshed state
	state.ID = types.Int64Value(int64(user.GetId()))
	state.Username = types.StringValue(user.GetUsername())
	state.FirstName = types.StringValue(user.GetFirstName())
	state.LastName = types.StringValue(user.GetLastName())
	state.Email = types.StringValue(user.GetEmail())
	state.IsActive = types.BoolValue(user.GetIsActive())
	state.IsSuperUser = types.BoolValue(user.GetIsSuperuser())

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("configuration_permissions"), user.ConfigurationPermissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Password.IsNull() {
		resp.Diagnostics.AddError(
			"Password Update Not Supported",
			"Password updates are not supported by the Defectdojo API. Please remove the password attribute from the update plan.",
		)
		return
	}

	configurations := make([]*int32, 0)
	diags = plan.ConfigurationPermissions.ElementsAs(ctx, &configurations, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate request from plan
	userRequest := defectdojo.UserRequest{
		Username:    plan.Username.ValueString(),
		FirstName:   plan.FirstName.ValueStringPointer(),
		LastName:    plan.LastName.ValueStringPointer(),
		Email:       plan.Email.ValueStringPointer(),
		IsActive:    plan.IsActive.ValueBoolPointer(),
		IsSuperuser: plan.IsSuperUser.ValueBoolPointer(),
	}

	// Update existing product type
	_, res, err := r.client.UsersAPI.UsersUpdate(ctx, int32(plan.ID.ValueInt64())).UserRequest(userRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Defectdojo User",
			"Could not update user with ID "+plan.ID.String()+": "+err.Error()+"\nDefectdojo responded with status:"+res.Status,
		)
		return
	}

	// Get refreshed product type value from Defectdojo
	user, res, err := r.client.UsersAPI.UsersRetrieve(ctx, int32(plan.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Defectdojo User",
			"Could not read user with ID "+plan.ID.String()+": "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(int64(user.GetId()))
	plan.Username = types.StringValue(user.GetUsername())
	plan.FirstName = types.StringValue(user.GetFirstName())
	plan.LastName = types.StringValue(user.GetLastName())
	plan.Email = types.StringValue(user.GetEmail())
	plan.IsActive = types.BoolValue(user.GetIsActive())
	plan.IsSuperUser = types.BoolValue(user.GetIsSuperuser())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("configuration_permissions"), user.ConfigurationPermissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing product type
	res, err := r.client.UsersAPI.UsersDestroy(ctx, int32(state.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Defectdojo User",
			"Could not delete user, unexpected error: "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			"Could not convert ID to integer: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(int64(id)))...)
}
