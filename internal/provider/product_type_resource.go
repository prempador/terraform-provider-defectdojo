package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/prempador/go-defectdojo"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &productTypeResource{}
	_ resource.ResourceWithConfigure   = &productTypeResource{}
	_ resource.ResourceWithImportState = &productTypeResource{}
)

// NewProductTypeResource is a helper function to simplify the provider implementation.
func NewProductTypeResource() resource.Resource {
	return &productTypeResource{}
}

// productTypeResource is the data source implementation.
type productTypeResource struct {
	client *defectdojo.APIClient
}

type productTypeResourceModel struct {
	ID              types.Int64  `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	CriticalProduct types.Bool   `tfsdk:"critical_product"`
	KeyProduct      types.Bool   `tfsdk:"key_product"`
	Members         types.List   `tfsdk:"members"`
	Authorization   types.List   `tfsdk:"authorization_groups"`
}

// Metadata returns the data source type name.
func (d *productTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_type"
}

// Schema defines the schema for the data source.
func (d *productTypeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The unique identifier for the product type",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the product type",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the product type",
				Computed:    true,
				Optional:    true,
			},
			"critical_product": schema.BoolAttribute{
				Description: "Whether the product type is a critical product",
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Optional:    true,
			},
			"key_product": schema.BoolAttribute{
				Description: "Whether the product type is a key product",
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Optional:    true,
			},
			"members": schema.ListAttribute{
				ElementType: types.Int64Type,
				Description: "The members of the product type",
				Computed:    true,
			},
			"authorization_groups": schema.ListAttribute{
				ElementType: types.Int64Type,
				Description: "The authorization groups of the product type",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *productTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *productTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan productTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate request from plan
	productTypeRequest := defectdojo.ProductTypeRequest{
		Name:            plan.Name.ValueString(),
		Description:     *defectdojo.NewNullableString(plan.Description.ValueStringPointer()),
		CriticalProduct: plan.CriticalProduct.ValueBoolPointer(),
		KeyProduct:      plan.KeyProduct.ValueBoolPointer(),
	}

	// Create new product type
	productType, res, err := r.client.ProductTypesAPI.ProductTypesCreate(ctx).ProductTypeRequest(productTypeRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Defectdojo Product Type",
			"Could not create product type, unexpected error: "+err.Error()+"\nDefectdojo responded with status:"+res.Status,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(int64(productType.GetId()))
	plan.Name = types.StringValue(productType.GetName())
	plan.Description = types.StringValue(productType.GetDescription())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("members"), productType.Members)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("authorization_groups"), productType.AuthorizationGroups)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *productTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state productTypeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed product type value from Defectdojo
	productType, res, err := r.client.ProductTypesAPI.ProductTypesRetrieve(ctx, int32(state.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Defectdojo Product Type",
			"Could not read product type with ID "+state.ID.String()+": "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}

	// Overwrite state with refreshed state
	state.ID = types.Int64Value(int64(productType.GetId()))
	state.Name = types.StringValue(productType.GetName())
	state.Description = types.StringValue(productType.GetDescription())
	state.CriticalProduct = types.BoolValue(productType.GetCriticalProduct())
	state.KeyProduct = types.BoolValue(productType.GetKeyProduct())

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("members"), productType.Members)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("authorization_groups"), productType.AuthorizationGroups)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *productTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan productTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate request from plan
	productTypeRequest := defectdojo.ProductTypeRequest{
		Name:            plan.Name.ValueString(),
		Description:     *defectdojo.NewNullableString(plan.Description.ValueStringPointer()),
		CriticalProduct: plan.CriticalProduct.ValueBoolPointer(),
		KeyProduct:      plan.KeyProduct.ValueBoolPointer(),
	}

	// Update existing product type
	_, res, err := r.client.ProductTypesAPI.ProductTypesUpdate(ctx, int32(plan.ID.ValueInt64())).ProductTypeRequest(productTypeRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Defectdojo Product Type",
			"Could not update product type with ID "+plan.ID.String()+": "+err.Error()+"\nDefectdojo responded with status:"+res.Status,
		)
		return
	}

	// Get refreshed product type value from Defectdojo
	productType, res, err := r.client.ProductTypesAPI.ProductTypesRetrieve(ctx, int32(plan.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Defectdojo Product Type",
			"Could not read product type with ID "+plan.ID.String()+": "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(int64(productType.GetId()))
	plan.Name = types.StringValue(productType.GetName())
	plan.Description = types.StringValue(productType.GetDescription())
	plan.CriticalProduct = types.BoolValue(productType.GetCriticalProduct())
	plan.KeyProduct = types.BoolValue(productType.GetKeyProduct())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("members"), productType.Members)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("authorization_groups"), productType.AuthorizationGroups)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *productTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state productTypeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing product type
	res, err := r.client.ProductTypesAPI.ProductTypesDestroy(ctx, int32(state.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Defectdojo Product Type",
			"Could not delete product type, unexpected error: "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}
}

func (r *productTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
