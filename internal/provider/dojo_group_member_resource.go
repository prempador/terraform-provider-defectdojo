// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/prempador/go-defectdojo"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dojoGroupMemberResource{}
	_ resource.ResourceWithConfigure   = &dojoGroupMemberResource{}
	_ resource.ResourceWithImportState = &dojoGroupMemberResource{}
)

var (
	roleName = map[int32]basetypes.StringValue{
		1: types.StringValue("API_Importer"),
		2: types.StringValue("Writer"),
		3: types.StringValue("Maintainer"),
		4: types.StringValue("Owner"),
		5: types.StringValue("Reader"),
	}

	roleID = map[basetypes.StringValue]int32{
		types.StringValue("API_Importer"): 1,
		types.StringValue("Writer"):       2,
		types.StringValue("Maintainer"):   3,
		types.StringValue("Owner"):        4,
		types.StringValue("Reader"):       5,
	}
)

// NewDojoGroupMemberResource is a helper function to simplify the provider implementation.
func NewDojoGroupMemberResource() resource.Resource {
	return &dojoGroupMemberResource{}
}

// dojoGroupMemberResource is the data source implementation.
type dojoGroupMemberResource struct {
	client *defectdojo.APIClient
}

type dojoGroupMemberResourceModel struct {
	ID    types.Int64  `tfsdk:"id"`
	User  types.Int64  `tfsdk:"user"`
	Group types.Int64  `tfsdk:"group"`
	Role  types.String `tfsdk:"role"`
}

// Metadata returns the data source type name.
func (r *dojoGroupMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dojo_group_member"
}

// Schema defines the schema for the data source.
func (r *dojoGroupMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The unique identifier of the dojo group",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"user": schema.Int64Attribute{
				Description: "The unique identifier of the user",
				Required:    true,
			},
			"group": schema.Int64Attribute{
				Description: "The unique identifier of the group",
				Required:    true,
			},
			"role": schema.StringAttribute{
				Description: "This role determines the permissions of the user to manage the group. The available roles are: API_Importer, Writer, Maintainer, Owner, Reader",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("API_Importer", "Writer", "Maintainer", "Owner", "Reader"),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *dojoGroupMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *dojoGroupMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan dojoGroupMemberResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate request from plan
	dojoGroupRequest := defectdojo.DojoGroupMemberRequest{
		User:  int32(plan.User.ValueInt64()),
		Group: int32(plan.Group.ValueInt64()),
		Role:  roleID[plan.Role],
	}

	// Create new group member
	dojoGroupMember, res, err := r.client.DojoGroupMembersAPI.DojoGroupMembersCreate(ctx).DojoGroupMemberRequest(dojoGroupRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Defectdojo Dojo Group Member",
			"Could not create dojo group member unexpected error: "+err.Error()+"\nDefectdojo responded with status:"+res.Status,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(int64(dojoGroupMember.GetId()))
	plan.User = types.Int64Value(int64(dojoGroupMember.GetUser()))
	plan.Group = types.Int64Value(int64(dojoGroupMember.GetGroup()))
	plan.Role = roleName[dojoGroupMember.GetRole()]

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *dojoGroupMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state dojoGroupMemberResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed group member value from Defectdojo
	dojoGroupMember, res, err := r.client.DojoGroupMembersAPI.DojoGroupMembersRetrieve(ctx, int32(state.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Defectdojo Dojo Group Member",
			"Could not read dojo group member with ID "+state.ID.String()+": "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}

	// Overwrite state with refreshed state
	state.ID = types.Int64Value(int64(dojoGroupMember.GetId()))
	state.User = types.Int64Value(int64(dojoGroupMember.GetUser()))
	state.Group = types.Int64Value(int64(dojoGroupMember.GetGroup()))
	state.Role = roleName[dojoGroupMember.GetRole()]

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dojoGroupMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan dojoGroupMemberResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate request from plan
	dojoGroupRequest := defectdojo.DojoGroupMemberRequest{
		User:  int32(plan.User.ValueInt64()),
		Group: int32(plan.Group.ValueInt64()),
		Role:  roleID[plan.Role],
	}

	// Update existing group member
	_, res, err := r.client.DojoGroupMembersAPI.DojoGroupMembersUpdate(ctx, int32(plan.ID.ValueInt64())).DojoGroupMemberRequest(dojoGroupRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Defectdojo Dojo Group Member",
			"Could not update dojo group member with ID "+plan.ID.String()+": "+err.Error()+"\nDefectdojo responded with status:"+res.Status,
		)
		return
	}

	// Get refreshed group member value from Defectdojo
	dojoGroupMember, res, err := r.client.DojoGroupMembersAPI.DojoGroupMembersRetrieve(ctx, int32(plan.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Defectdojo Dojo Group Member",
			"Could not read dojo group member with ID "+plan.ID.String()+": "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(int64(dojoGroupMember.GetId()))
	plan.User = types.Int64Value(int64(dojoGroupMember.GetUser()))
	plan.Group = types.Int64Value(int64(dojoGroupMember.GetGroup()))
	plan.Role = roleName[dojoGroupMember.GetRole()]

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dojoGroupMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state dojoGroupMemberResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing group member
	res, err := r.client.DojoGroupMembersAPI.DojoGroupMembersDestroy(ctx, int32(state.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Defectdojo Dojo Group Member",
			"Could not delete dojo group member, unexpected error: "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}
}

func (r *dojoGroupMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
