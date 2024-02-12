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
	"github.com/prempador/go-defectdojo"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &productResource{}
	_ resource.ResourceWithConfigure   = &productResource{}
	_ resource.ResourceWithImportState = &productResource{}
)

// NewProductResource is a helper function to simplify the provider implementation.
func NewProductResource() resource.Resource {
	return &productResource{}
}

// productResource is the data source implementation.
type productResource struct {
	client *defectdojo.APIClient
}

type productResourceModel struct {
	ID                            types.Int64  `tfsdk:"id"`
	Name                          types.String `tfsdk:"name"`
	Description                   types.String `tfsdk:"description"`
	ProdNumericGrade              types.Int64  `tfsdk:"prod_numeric_grade"`
	BusinessCriticality           types.String `tfsdk:"business_criticality"`
	Platform                      types.String `tfsdk:"platform"`
	Lifecycle                     types.String `tfsdk:"product_lifecycle"`
	Origin                        types.String `tfsdk:"origin"`
	UserRecords                   types.Int64  `tfsdk:"user_records"`
	Revenue                       types.String `tfsdk:"revenue"`
	ExternalAudience              types.Bool   `tfsdk:"external_audience"`
	InternetAccessible            types.Bool   `tfsdk:"internet_accessible"`
	EnableProductTagInheritance   types.Bool   `tfsdk:"enable_product_tag_inheritance"`
	EnableSimpleRiskAcceptance    types.Bool   `tfsdk:"enable_simple_risk_acceptance"`
	EnableFullRiskAcceptance      types.Bool   `tfsdk:"enable_full_risk_acceptance"`
	DisableSlaBreachNotifications types.Bool   `tfsdk:"disable_sla_breach_notifications"`
	ProductManager                types.Int64  `tfsdk:"product_manager"`
	TechnicalContact              types.Int64  `tfsdk:"technical_contact"`
	TeamManager                   types.Int64  `tfsdk:"team_manager"`
	ProdType                      types.Int64  `tfsdk:"prod_type"`
	SlaConfiguration              types.Int64  `tfsdk:"sla_configuration"`
	Regulations                   types.List   `tfsdk:"regulations"`
	Tags                          types.List   `tfsdk:"tags"`
}

// Metadata returns the data source type name.
func (r *productResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product"
}

// Schema defines the schema for the data source.
func (r *productResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The unique identifier of the product",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the product",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the product",
				Required:    true,
			},
			"prod_numeric_grade": schema.Int64Attribute{
				Description: "The numeric grade of the product",
				Computed:    true,
				Optional:    true,
			},
			"business_criticality": schema.StringAttribute{
				Description: "The business criticality of the product",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("very high", "high", "medium", "low", "very low", "none"),
				},
			},
			"platform": schema.StringAttribute{
				Description: "The platform of the product",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("web service", "desktop", "iot", "mobile", "web"),
				},
			},
			"product_lifecycle": schema.StringAttribute{
				Description: "The lifecycle of the product (renamed to product_lifecycle from the API to avoid conflict with lifecycle attribute)",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("construction", "production", "retirement"),
				},
			},
			"origin": schema.StringAttribute{
				Description: "The origin of the product",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("third party library", "purchased", "contractor", "internal", "open source", "outsourced"),
				},
			},
			"user_records": schema.Int64Attribute{
				Description: "Estimate the number of user records within the application",
				Computed:    true,
				Optional:    true,
			},
			"revenue": schema.StringAttribute{
				Description: "Estimate the application's revenue",
				Computed:    true,
				Optional:    true,
			},
			"external_audience": schema.BoolAttribute{
				Description: "Specify if the application is used by people outside the organization",
				Computed:    true,
				Optional:    true,
			},
			"internet_accessible": schema.BoolAttribute{
				Description: "Specify if the application is accessible from the public internet",
				Computed:    true,
				Optional:    true,
			},
			"enable_product_tag_inheritance": schema.BoolAttribute{
				Description: "Enables product tag inheritance. Any tags added on a product will automatically be added to all Engagements, Tests, and Findings",
				Computed:    true,
				Optional:    true,
			},
			"enable_simple_risk_acceptance": schema.BoolAttribute{
				Description: "Allows simple risk acceptance by checking/unchecking a checkbox",
				Computed:    true,
				Optional:    true,
			},
			"enable_full_risk_acceptance": schema.BoolAttribute{
				Description: "Allows full risk acceptance using a risk acceptance form, expiration date, uploaded proof, etc.",
				Computed:    true,
				Optional:    true,
			},
			"disable_sla_breach_notifications": schema.BoolAttribute{
				Description: "Disable SLA breach notifications if configured in the global settings",
				Computed:    true,
				Optional:    true,
			},
			"product_manager": schema.Int64Attribute{
				Description: "The product manager of the product",
				Computed:    true,
				Optional:    true,
			},
			"technical_contact": schema.Int64Attribute{
				Description: "The technical contact of the product",
				Computed:    true,
				Optional:    true,
			},
			"team_manager": schema.Int64Attribute{
				Description: "The team manager of the product",
				Computed:    true,
				Optional:    true,
			},
			"prod_type": schema.Int64Attribute{
				Description: "The product type of the product",
				Required:    true,
			},
			"sla_configuration": schema.Int64Attribute{
				Description: "The SLA configuration of the product",
				Computed:    true,
				Optional:    true,
			},
			"regulations": schema.ListAttribute{
				ElementType: types.Int64Type,
				Description: "List of regulations for the product",
				Computed:    true,
				Optional:    true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "List of tags for the product",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *productResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *productResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan productResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert string to float
	var revenue *float64

	if plan.Revenue.ValueStringPointer() != nil && plan.Revenue.ValueString() != "" {
		r, err := strconv.ParseFloat(plan.Revenue.ValueString(), 64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Revenue",
				"Revenue must be a valid float: "+err.Error(),
			)
			return
		}

		revenue = &r
	}

	tags := make([]string, 0)
	diags = plan.Tags.ElementsAs(ctx, &tags, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	regulations := make([]int32, 0)
	diags = plan.Regulations.ElementsAs(ctx, &regulations, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate request from plan
	productRequest := defectdojo.ProductRequest{
		Name:                          plan.Name.ValueString(),
		Description:                   plan.Description.ValueString(),
		ProdNumericGrade:              basetypesInt64ValueToDefectdojoNullableInt32(plan.ProdNumericGrade),
		BusinessCriticality:           *defectdojo.NewNullableString(plan.BusinessCriticality.ValueStringPointer()),
		Platform:                      *defectdojo.NewNullableString(plan.Platform.ValueStringPointer()),
		Lifecycle:                     *defectdojo.NewNullableString(plan.Lifecycle.ValueStringPointer()),
		Origin:                        *defectdojo.NewNullableString(plan.Origin.ValueStringPointer()),
		UserRecords:                   basetypesInt64ValueToDefectdojoNullableInt32(plan.UserRecords),
		Revenue:                       *defectdojo.NewNullableFloat64(revenue),
		ExternalAudience:              plan.ExternalAudience.ValueBoolPointer(),
		InternetAccessible:            plan.InternetAccessible.ValueBoolPointer(),
		EnableProductTagInheritance:   plan.EnableProductTagInheritance.ValueBoolPointer(),
		EnableSimpleRiskAcceptance:    plan.EnableSimpleRiskAcceptance.ValueBoolPointer(),
		EnableFullRiskAcceptance:      plan.EnableFullRiskAcceptance.ValueBoolPointer(),
		DisableSlaBreachNotifications: plan.DisableSlaBreachNotifications.ValueBoolPointer(),
		ProductManager:                basetypesInt64ValueToDefectdojoNullableInt32(plan.ProductManager),
		TechnicalContact:              basetypesInt64ValueToDefectdojoNullableInt32(plan.TechnicalContact),
		TeamManager:                   basetypesInt64ValueToDefectdojoNullableInt32(plan.TeamManager),
		ProdType:                      int32(plan.ProdType.ValueInt64()),
		SlaConfiguration:              basetypesInt64ValueToInt32Pointer(plan.SlaConfiguration),
	}

	// Create new product
	product, res, err := r.client.ProductsAPI.ProductsCreate(ctx).ProductRequest(productRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Defectdojo Product",
			"Could not create product, unexpected error: "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(int64(product.GetId()))
	plan.Name = types.StringValue(product.GetName())
	plan.Description = types.StringValue(product.GetDescription())
	plan.ProdNumericGrade = types.Int64Value(int64(product.GetProdNumericGrade()))
	plan.BusinessCriticality = types.StringValue(product.GetBusinessCriticality())
	plan.Platform = types.StringValue(product.GetPlatform())
	plan.Lifecycle = types.StringValue(product.GetLifecycle())
	plan.Origin = types.StringValue(product.GetOrigin())
	plan.UserRecords = types.Int64Value(int64(product.GetUserRecords()))
	plan.Revenue = types.StringValue(fmt.Sprintf("%f", product.GetRevenue()))
	plan.ExternalAudience = types.BoolValue(product.GetExternalAudience())
	plan.InternetAccessible = types.BoolValue(product.GetInternetAccessible())
	plan.EnableProductTagInheritance = types.BoolValue(product.GetEnableProductTagInheritance())
	plan.EnableSimpleRiskAcceptance = types.BoolValue(product.GetEnableSimpleRiskAcceptance())
	plan.EnableFullRiskAcceptance = types.BoolValue(product.GetEnableFullRiskAcceptance())
	plan.DisableSlaBreachNotifications = types.BoolValue(product.GetDisableSlaBreachNotifications())
	plan.ProductManager = types.Int64Value(int64(product.GetProductManager()))
	plan.TechnicalContact = types.Int64Value(int64(product.GetTechnicalContact()))
	plan.TeamManager = types.Int64Value(int64(product.GetTeamManager()))
	plan.ProdType = types.Int64Value(int64(product.GetProdType()))
	plan.SlaConfiguration = types.Int64Value(int64(product.GetSlaConfiguration()))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("regulations"), product.Regulations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("tags"), product.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *productResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state productResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed product value from Defectdojo
	product, res, err := r.client.ProductsAPI.ProductsRetrieve(ctx, int32(state.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Defectdojo Product",
			"Could not read product with ID "+state.ID.String()+": "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}

	// Overwrite state with refreshed state
	state.ID = types.Int64Value(int64(product.GetId()))
	state.Name = types.StringValue(product.GetName())
	state.Description = types.StringValue(product.GetDescription())
	state.ProdNumericGrade = types.Int64Value(int64(product.GetProdNumericGrade()))
	state.BusinessCriticality = types.StringValue(product.GetBusinessCriticality())
	state.Platform = types.StringValue(product.GetPlatform())
	state.Lifecycle = types.StringValue(product.GetLifecycle())
	state.Origin = types.StringValue(product.GetOrigin())
	state.UserRecords = types.Int64Value(int64(product.GetUserRecords()))
	state.Revenue = types.StringValue(fmt.Sprintf("%f", product.GetRevenue()))
	state.ExternalAudience = types.BoolValue(product.GetExternalAudience())
	state.InternetAccessible = types.BoolValue(product.GetInternetAccessible())
	state.EnableProductTagInheritance = types.BoolValue(product.GetEnableProductTagInheritance())
	state.EnableSimpleRiskAcceptance = types.BoolValue(product.GetEnableSimpleRiskAcceptance())
	state.EnableFullRiskAcceptance = types.BoolValue(product.GetEnableFullRiskAcceptance())
	state.DisableSlaBreachNotifications = types.BoolValue(product.GetDisableSlaBreachNotifications())
	state.ProductManager = types.Int64Value(int64(product.GetProductManager()))
	state.TechnicalContact = types.Int64Value(int64(product.GetTechnicalContact()))
	state.TeamManager = types.Int64Value(int64(product.GetTeamManager()))
	state.ProdType = types.Int64Value(int64(product.GetProdType()))
	state.SlaConfiguration = types.Int64Value(int64(product.GetSlaConfiguration()))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("regulations"), product.Regulations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("tags"), product.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *productResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan productResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert string to float
	var revenue *float64

	if plan.Revenue.ValueStringPointer() != nil && plan.Revenue.ValueString() != "" {
		r, err := strconv.ParseFloat(plan.Revenue.ValueString(), 64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Revenue",
				"Revenue must be a valid float: "+err.Error(),
			)
			return
		}

		revenue = &r
	}

	tags := make([]string, 0)
	diags = plan.Tags.ElementsAs(ctx, &tags, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	regulations := make([]int32, 0)
	diags = plan.Regulations.ElementsAs(ctx, &regulations, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate request from plan
	productRequest := defectdojo.ProductRequest{
		Name:                          plan.Name.ValueString(),
		Description:                   plan.Description.ValueString(),
		ProdNumericGrade:              basetypesInt64ValueToDefectdojoNullableInt32(plan.ProdNumericGrade),
		BusinessCriticality:           *defectdojo.NewNullableString(plan.BusinessCriticality.ValueStringPointer()),
		Platform:                      *defectdojo.NewNullableString(plan.Platform.ValueStringPointer()),
		Lifecycle:                     *defectdojo.NewNullableString(plan.Lifecycle.ValueStringPointer()),
		Origin:                        *defectdojo.NewNullableString(plan.Origin.ValueStringPointer()),
		UserRecords:                   basetypesInt64ValueToDefectdojoNullableInt32(plan.UserRecords),
		Revenue:                       *defectdojo.NewNullableFloat64(revenue),
		ExternalAudience:              plan.ExternalAudience.ValueBoolPointer(),
		InternetAccessible:            plan.InternetAccessible.ValueBoolPointer(),
		EnableProductTagInheritance:   plan.EnableProductTagInheritance.ValueBoolPointer(),
		EnableSimpleRiskAcceptance:    plan.EnableSimpleRiskAcceptance.ValueBoolPointer(),
		EnableFullRiskAcceptance:      plan.EnableFullRiskAcceptance.ValueBoolPointer(),
		DisableSlaBreachNotifications: plan.DisableSlaBreachNotifications.ValueBoolPointer(),
		ProductManager:                basetypesInt64ValueToDefectdojoNullableInt32(plan.ProductManager),
		TechnicalContact:              basetypesInt64ValueToDefectdojoNullableInt32(plan.TechnicalContact),
		TeamManager:                   basetypesInt64ValueToDefectdojoNullableInt32(plan.TeamManager),
		ProdType:                      int32(plan.ProdType.ValueInt64()),
		SlaConfiguration:              basetypesInt64ValueToInt32Pointer(plan.SlaConfiguration),
	}

	// Update existing product
	_, res, err := r.client.ProductsAPI.ProductsUpdate(ctx, int32(plan.ID.ValueInt64())).ProductRequest(productRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Defectdojo Product",
			"Could not update product with ID "+plan.ID.String()+": "+err.Error()+"\nDefectdojo responded with status:"+res.Status,
		)
		return
	}

	// Get refreshed product value from Defectdojo
	product, res, err := r.client.ProductsAPI.ProductsRetrieve(ctx, int32(plan.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Defectdojo Product",
			"Could not read product with ID "+plan.ID.String()+": "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(int64(product.GetId()))
	plan.Name = types.StringValue(product.GetName())
	plan.Description = types.StringValue(product.GetDescription())
	plan.ProdNumericGrade = types.Int64Value(int64(product.GetProdNumericGrade()))
	plan.BusinessCriticality = types.StringValue(product.GetBusinessCriticality())
	plan.Platform = types.StringValue(product.GetPlatform())
	plan.Lifecycle = types.StringValue(product.GetLifecycle())
	plan.Origin = types.StringValue(product.GetOrigin())
	plan.UserRecords = types.Int64Value(int64(product.GetUserRecords()))
	plan.Revenue = types.StringValue(fmt.Sprintf("%f", product.GetRevenue()))
	plan.ExternalAudience = types.BoolValue(product.GetExternalAudience())
	plan.InternetAccessible = types.BoolValue(product.GetInternetAccessible())
	plan.EnableProductTagInheritance = types.BoolValue(product.GetEnableProductTagInheritance())
	plan.EnableSimpleRiskAcceptance = types.BoolValue(product.GetEnableSimpleRiskAcceptance())
	plan.EnableFullRiskAcceptance = types.BoolValue(product.GetEnableFullRiskAcceptance())
	plan.DisableSlaBreachNotifications = types.BoolValue(product.GetDisableSlaBreachNotifications())
	plan.ProductManager = types.Int64Value(int64(product.GetProductManager()))
	plan.TechnicalContact = types.Int64Value(int64(product.GetTechnicalContact()))
	plan.TeamManager = types.Int64Value(int64(product.GetTeamManager()))
	plan.ProdType = types.Int64Value(int64(product.GetProdType()))
	plan.SlaConfiguration = types.Int64Value(int64(product.GetSlaConfiguration()))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("regulations"), product.Regulations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("tags"), product.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *productResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state productResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing product
	res, err := r.client.ProductsAPI.ProductsDestroy(ctx, int32(state.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Defectdojo Product",
			"Could not delete product, unexpected error: "+err.Error()+"\nDefectdojo responded with status: "+res.Status,
		)
		return
	}
}

func (r *productResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
