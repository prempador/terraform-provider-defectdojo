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
	_ resource.Resource                = &dojoGroupResource{}
	_ resource.ResourceWithConfigure   = &dojoGroupResource{}
	_ resource.ResourceWithImportState = &dojoGroupResource{}
)

// NewdojoGroupResource is a helper function to simplify the provider implementation.
func NewDojoGroupResource() resource.Resource {
	return &dojoGroupResource{}
}

// dojoGroupResource is the data source implementation.
type dojoGroupResource struct {
	client *defectdojo.APIClient
}

type dojoGroupResourceModel struct {
	ID                       types.Int64  `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Description              types.String `tfsdk:"description"`
	ConfigurationPermissions types.List   `tfsdk:"configuration_permissions"`
	SocialProvider           types.String `tfsdk:"social_provider"`
}

// Metadata returns the data source type name.
func (r *dojoGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dojo_group"
}

// Schema defines the schema for the data source.
func (r *dojoGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The unique identifier of the dojo group",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the dojo group",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the dojo group",
				Computed:    true,
				Optional:    true,
			},
			"configuration_permissions": schema.ListAttribute{
				ElementType: types.Int64Type,
				Description: "The configuration permissions of the dojo group",
				Computed:    true,
				Optional:    true,
			},
			"social_provider": schema.StringAttribute{
				Description: "The social provider the dojo group was imported from",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *dojoGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *dojoGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan dojoGroupResourceModel
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
	dojoGroupRequest := defectdojo.DojoGroupRequest{
		Name:                     plan.Name.ValueString(),
		Description:              *defectdojo.NewNullableString(plan.Description.ValueStringPointer()),
		ConfigurationPermissions: configurations,
		SocialProvider:           *defectdojo.NewNullableString(plan.SocialProvider.ValueStringPointer()),
	}

	// Create new dojo group
	dojoGroup, res, err := r.client.DojoGroupsAPI.DojoGroupsCreate(ctx).DojoGroupRequest(dojoGroupRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Defectdojo Dojo Group",
			"Could not create dojo group unexpected error: "+err.Error()+"\nDefectdojo responded with status:"+res.Status,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(int64(dojoGroup.GetId()))
	plan.Name = types.StringValue(dojoGroup.GetName())
	plan.Description = types.StringValue(dojoGroup.GetDescription())
	plan.SocialProvider = types.StringValue(dojoGroup.GetSocialProvider())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("configuration_permissions"), dojoGroup.ConfigurationPermissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *dojoGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state dojoGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed dojo group value from Defectdojo
	dojoGroup, res, err := r.client.DojoGroupsAPI.DojoGroupsRetrieve(ctx, int32(state.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Defectdojo Dojo Group",
			"Could not read dojo group with ID "+state.ID.String()+": "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}

	// Overwrite state with refreshed state
	state.ID = types.Int64Value(int64(dojoGroup.GetId()))
	state.Name = types.StringValue(dojoGroup.GetName())
	state.Description = types.StringValue(dojoGroup.GetDescription())
	state.SocialProvider = types.StringValue(dojoGroup.GetSocialProvider())

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("configuration_permissions"), dojoGroup.ConfigurationPermissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dojoGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan dojoGroupResourceModel
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
	dojoGroupRequest := defectdojo.DojoGroupRequest{
		Name:                     plan.Name.ValueString(),
		Description:              *defectdojo.NewNullableString(plan.Description.ValueStringPointer()),
		ConfigurationPermissions: configurations,
		SocialProvider:           *defectdojo.NewNullableString(plan.SocialProvider.ValueStringPointer()),
	}

	// Update existing dojo group
	_, res, err := r.client.DojoGroupsAPI.DojoGroupsUpdate(ctx, int32(plan.ID.ValueInt64())).DojoGroupRequest(dojoGroupRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Defectdojo Dojo Group",
			"Could not update dojo group with ID "+plan.ID.String()+": "+err.Error()+"\nDefectdojo responded with status:"+res.Status,
		)
		return
	}

	// Get refreshed dojo group value from Defectdojo
	dojoGroup, res, err := r.client.DojoGroupsAPI.DojoGroupsRetrieve(ctx, int32(plan.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Defectdojo Dojo Group",
			"Could not read dojo group with ID "+plan.ID.String()+": "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(int64(dojoGroup.GetId()))
	plan.Name = types.StringValue(dojoGroup.GetName())
	plan.Description = types.StringValue(dojoGroup.GetDescription())
	plan.SocialProvider = types.StringValue(dojoGroup.GetSocialProvider())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("configuration_permissions"), dojoGroup.ConfigurationPermissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dojoGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state dojoGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing dojo group
	res, err := r.client.DojoGroupsAPI.DojoGroupsDestroy(ctx, int32(state.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Defectdojo Dojo Group",
			"Could not delete dojo group, unexpected error: "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}
}

func (r *dojoGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
