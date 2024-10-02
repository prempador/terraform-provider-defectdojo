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
	_ resource.Resource                = &engagementResource{}
	_ resource.ResourceWithConfigure   = &engagementResource{}
	_ resource.ResourceWithImportState = &engagementResource{}
)

// NewEngagementResource is a helper function to simplify the provider implementation.
func NewEngagementResource() resource.Resource {
	return &engagementResource{}
}

// engagementResource is the data source implementation.
type engagementResource struct {
	client *defectdojo.APIClient
}

type engagementResourceModel struct {
	ID                         types.Int64  `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Description                types.String `tfsdk:"description"`
	Version                    types.String `tfsdk:"version"`
	FirstContacted             types.String `tfsdk:"first_contacted"`
	TargetStart                types.String `tfsdk:"target_start"`
	TargetEnd                  types.String `tfsdk:"target_end"`
	Reason                     types.String `tfsdk:"reason"`
	Tracker                    types.String `tfsdk:"tracker"`
	TestStrategy               types.String `tfsdk:"test_strategy"`
	ThreatModel                types.Bool   `tfsdk:"threat_model"`
	APITest                    types.Bool   `tfsdk:"api_test"`
	PenTest                    types.Bool   `tfsdk:"pen_test"`
	CheckList                  types.Bool   `tfsdk:"check_list"`
	Status                     types.String `tfsdk:"status"`
	EngagementType             types.String `tfsdk:"engagement_type"`
	BuildID                    types.String `tfsdk:"build_id"`
	CommitHash                 types.String `tfsdk:"commit_hash"`
	BranchTag                  types.String `tfsdk:"branch_tag"`
	SourceCodeManagementURI    types.String `tfsdk:"source_code_management_uri"`
	DeduplicationOnEngagement  types.Bool   `tfsdk:"deduplication_on_engagement"`
	Lead                       types.Int64  `tfsdk:"lead"`
	Requester                  types.Int64  `tfsdk:"requester"`
	Preset                     types.Int64  `tfsdk:"preset"`
	ReportType                 types.Int64  `tfsdk:"report_type"`
	Product                    types.Int64  `tfsdk:"product"`
	BuildServer                types.Int64  `tfsdk:"build_server"`
	SourceCodeManagementServer types.Int64  `tfsdk:"source_code_management_server"`
	OrchestrationEngine        types.Int64  `tfsdk:"orchestration_engine"`
	Tags                       types.List   `tfsdk:"tags"`
}

// Metadata returns the data source type name.
func (r *engagementResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engagement"
}

// Schema defines the schema for the data source.
func (r *engagementResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The unique identifier of the engagement",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the engagement",
				Computed:    true,
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the engagement",
				Computed:    true,
				Optional:    true,
			},
			"version": schema.StringAttribute{
				Description: "Version of the product the engagement tested",
				Computed:    true,
				Optional:    true,
			},
			"first_contacted": schema.StringAttribute{
				Description: "The date the engagement was first contacted",
				Computed:    true,
				Optional:    true,
			},
			"target_start": schema.StringAttribute{
				Description: "The date the engagement is targeted to start",
				Required:    true,
			},
			"target_end": schema.StringAttribute{
				Description: "The date the engagement is targeted to end",
				Required:    true,
			},
			"reason": schema.StringAttribute{
				Description: "The reason for the engagement",
				Computed:    true,
				Optional:    true,
			},
			"tracker": schema.StringAttribute{
				Description: "Link to epic or ticket system with changes to version",
				Computed:    true,
				Optional:    true,
			},
			"test_strategy": schema.StringAttribute{
				Description: "The test strategy for the engagement",
				Computed:    true,
				Optional:    true,
			},
			"threat_model": schema.BoolAttribute{
				Description: "Whether the engagement includes a threat model",
				Computed:    true,
				Optional:    true,
			},
			"api_test": schema.BoolAttribute{
				Description: "Whether the engagement includes an API test",
				Computed:    true,
				Optional:    true,
			},
			"pen_test": schema.BoolAttribute{
				Description: "Whether the engagement includes a pen test",
				Computed:    true,
				Optional:    true,
			},
			"check_list": schema.BoolAttribute{
				Description: "Whether the engagement includes a check list",
				Computed:    true,
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the engagement",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Not Started", "Blocked", "Cancelled", "Completed", "In Progress", "On Hold", "Waiting for Resource"),
				},
			},
			"engagement_type": schema.StringAttribute{
				Description: "The type of engagement",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Interactive", "CI/CD"),
				},
			},
			"build_id": schema.StringAttribute{
				Description: "Build ID of the product the engagement tested",
				Computed:    true,
				Optional:    true,
			},
			"commit_hash": schema.StringAttribute{
				Description: "Commit hash from repo",
				Computed:    true,
				Optional:    true,
			},
			"branch_tag": schema.StringAttribute{
				Description: "Tag or branch of the product the engagement tested",
				Computed:    true,
				Optional:    true,
			},
			"source_code_management_uri": schema.StringAttribute{
				Description: "Resource link to source code",
				Computed:    true,
				Optional:    true,
			},
			"deduplication_on_engagement": schema.BoolAttribute{
				Description: "If enabled deduplication will only mark a finding in this engagement as duplicate of another finding if both findings are in this engagement. If disabled, deduplication is on the product level",
				Computed:    true,
				Optional:    true,
			},
			"lead": schema.Int64Attribute{
				Description: "The user ID of the engagement lead",
				Computed:    true,
				Optional:    true,
			},
			"requester": schema.Int64Attribute{
				Description: "The user ID of the engagement requester",
				Computed:    true,
				Optional:    true,
			},
			"preset": schema.Int64Attribute{
				Description: "Settings and notes for performing this engagement",
				Computed:    true,
				Optional:    true,
			},
			"report_type": schema.Int64Attribute{
				Description: "The report type for the engagement",
				Computed:    true,
				Optional:    true,
			},
			"product": schema.Int64Attribute{
				Description: "The product ID of the engagement",
				Required:    true,
			},
			"build_server": schema.Int64Attribute{
				Description: "Build server responsible for CI/CD test",
				Computed:    true,
				Optional:    true,
			},
			"source_code_management_server": schema.Int64Attribute{
				Description: "Source code server for CI/CD test",
				Computed:    true,
				Optional:    true,
			},
			"orchestration_engine": schema.Int64Attribute{
				Description: "Orchestration engine responsible for CI/CD test",
				Computed:    true,
				Optional:    true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "List of tags for the engagement",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *engagementResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *engagementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan engagementResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags := make([]string, 0)
	diags = plan.Tags.ElementsAs(ctx, &tags, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate request from plan
	engagementRequest := defectdojo.EngagementRequest{
		Name:                       *defectdojo.NewNullableString(plan.Name.ValueStringPointer()),
		Description:                *defectdojo.NewNullableString(plan.Description.ValueStringPointer()),
		Version:                    *defectdojo.NewNullableString(plan.Version.ValueStringPointer()),
		FirstContacted:             basetypesStringValueToDefectdojoNullableString(plan.FirstContacted),
		TargetStart:                plan.TargetStart.ValueString(),
		TargetEnd:                  plan.TargetEnd.ValueString(),
		Reason:                     *defectdojo.NewNullableString(plan.Reason.ValueStringPointer()),
		Tracker:                    *defectdojo.NewNullableString(plan.Tracker.ValueStringPointer()),
		TestStrategy:               *defectdojo.NewNullableString(plan.TestStrategy.ValueStringPointer()),
		ThreatModel:                plan.ThreatModel.ValueBoolPointer(),
		ApiTest:                    plan.APITest.ValueBoolPointer(),
		PenTest:                    plan.PenTest.ValueBoolPointer(),
		CheckList:                  plan.CheckList.ValueBoolPointer(),
		Status:                     basetypesStringValueToDefectdojoNullableString(plan.Status),
		EngagementType:             basetypesStringValueToDefectdojoNullableString(plan.EngagementType),
		BuildId:                    *defectdojo.NewNullableString(plan.BuildID.ValueStringPointer()),
		CommitHash:                 *defectdojo.NewNullableString(plan.CommitHash.ValueStringPointer()),
		BranchTag:                  *defectdojo.NewNullableString(plan.BranchTag.ValueStringPointer()),
		SourceCodeManagementUri:    *defectdojo.NewNullableString(plan.SourceCodeManagementURI.ValueStringPointer()),
		DeduplicationOnEngagement:  plan.DeduplicationOnEngagement.ValueBoolPointer(),
		Lead:                       basetypesInt64ValueToDefectdojoNullableInt32(plan.Lead),
		Requester:                  basetypesInt64ValueToDefectdojoNullableInt32(plan.Requester),
		Preset:                     basetypesInt64ValueToDefectdojoNullableInt32(plan.Preset),
		ReportType:                 basetypesInt64ValueToDefectdojoNullableInt32(plan.ReportType),
		Product:                    int32(plan.Product.ValueInt64()),
		BuildServer:                basetypesInt64ValueToDefectdojoNullableInt32(plan.BuildServer),
		SourceCodeManagementServer: basetypesInt64ValueToDefectdojoNullableInt32(plan.SourceCodeManagementServer),
		OrchestrationEngine:        basetypesInt64ValueToDefectdojoNullableInt32(plan.OrchestrationEngine),
		Tags:                       tags,
	}

	// Create new engagement
	engagement, res, err := r.client.EngagementsAPI.EngagementsCreate(ctx).EngagementRequest(engagementRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Defectdojo Engagement",
			"Could not create engagement, unexpected error: "+err.Error()+"\nDefectdojo responded with status: "+fmt.Sprintf("%+v", res),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(int64(engagement.GetId()))
	plan.Name = types.StringValue(engagement.GetName())
	plan.Description = types.StringValue(engagement.GetDescription())
	plan.Version = types.StringValue(engagement.GetVersion())
	plan.FirstContacted = types.StringValue(engagement.GetFirstContacted())
	plan.TargetStart = types.StringValue(engagement.GetTargetStart())
	plan.TargetEnd = types.StringValue(engagement.GetTargetEnd())
	plan.Reason = types.StringValue(engagement.GetReason())
	plan.Tracker = types.StringValue(engagement.GetTracker())
	plan.TestStrategy = types.StringValue(engagement.GetTestStrategy())
	plan.ThreatModel = types.BoolValue(engagement.GetThreatModel())
	plan.APITest = types.BoolValue(engagement.GetApiTest())
	plan.PenTest = types.BoolValue(engagement.GetPenTest())
	plan.CheckList = types.BoolValue(engagement.GetCheckList())
	plan.Status = types.StringValue(engagement.GetStatus())
	plan.EngagementType = types.StringValue(engagement.GetEngagementType())
	plan.BuildID = types.StringValue(engagement.GetBuildId())
	plan.CommitHash = types.StringValue(engagement.GetCommitHash())
	plan.BranchTag = types.StringValue(engagement.GetBranchTag())
	plan.SourceCodeManagementURI = types.StringValue(engagement.GetSourceCodeManagementUri())
	plan.DeduplicationOnEngagement = types.BoolValue(engagement.GetDeduplicationOnEngagement())
	plan.Lead = types.Int64Value(int64(engagement.GetLead()))
	plan.Requester = types.Int64Value(int64(engagement.GetRequester()))
	plan.Preset = types.Int64Value(int64(engagement.GetPreset()))
	plan.ReportType = types.Int64Value(int64(engagement.GetReportType()))
	plan.Product = types.Int64Value(int64(engagement.GetProduct()))
	plan.BuildServer = types.Int64Value(int64(engagement.GetBuildServer()))
	plan.SourceCodeManagementServer = types.Int64Value(int64(engagement.GetSourceCodeManagementServer()))
	plan.OrchestrationEngine = types.Int64Value(int64(engagement.GetOrchestrationEngine()))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("tags"), engagement.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *engagementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state engagementResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed engagement value from Defectdojo
	engagement, res, err := r.client.EngagementsAPI.EngagementsRetrieve(ctx, int32(state.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Defectdojo Engagement",
			"Could not read engagement with ID "+state.ID.String()+": "+err.Error()+"\nDefectdojo responded with status:  "+fmt.Sprintf("%v", res.Body),
		)
		return
	}

	// Overwrite state with refreshed state
	state.ID = types.Int64Value(int64(engagement.GetId()))
	state.Name = types.StringValue(engagement.GetName())
	state.Description = types.StringValue(engagement.GetDescription())
	state.Version = types.StringValue(engagement.GetVersion())
	state.FirstContacted = types.StringValue(engagement.GetFirstContacted())
	state.TargetStart = types.StringValue(engagement.GetTargetStart())
	state.TargetEnd = types.StringValue(engagement.GetTargetEnd())
	state.Reason = types.StringValue(engagement.GetReason())
	state.Tracker = types.StringValue(engagement.GetTracker())
	state.TestStrategy = types.StringValue(engagement.GetTestStrategy())
	state.ThreatModel = types.BoolValue(engagement.GetThreatModel())
	state.APITest = types.BoolValue(engagement.GetApiTest())
	state.PenTest = types.BoolValue(engagement.GetPenTest())
	state.CheckList = types.BoolValue(engagement.GetCheckList())
	state.Status = types.StringValue(engagement.GetStatus())
	state.EngagementType = types.StringValue(engagement.GetEngagementType())
	state.BuildID = types.StringValue(engagement.GetBuildId())
	state.CommitHash = types.StringValue(engagement.GetCommitHash())
	state.BranchTag = types.StringValue(engagement.GetBranchTag())
	state.SourceCodeManagementURI = types.StringValue(engagement.GetSourceCodeManagementUri())
	state.DeduplicationOnEngagement = types.BoolValue(engagement.GetDeduplicationOnEngagement())
	state.Lead = types.Int64Value(int64(engagement.GetLead()))
	state.Requester = types.Int64Value(int64(engagement.GetRequester()))
	state.Preset = types.Int64Value(int64(engagement.GetPreset()))
	state.ReportType = types.Int64Value(int64(engagement.GetReportType()))
	state.Product = types.Int64Value(int64(engagement.GetProduct()))
	state.BuildServer = types.Int64Value(int64(engagement.GetBuildServer()))
	state.SourceCodeManagementServer = types.Int64Value(int64(engagement.GetSourceCodeManagementServer()))
	state.OrchestrationEngine = types.Int64Value(int64(engagement.GetOrchestrationEngine()))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("tags"), engagement.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *engagementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan engagementResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags := make([]string, 0)
	diags = plan.Tags.ElementsAs(ctx, &tags, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate request from plan
	engagementRequest := defectdojo.EngagementRequest{
		Name:                       *defectdojo.NewNullableString(plan.Name.ValueStringPointer()),
		Description:                *defectdojo.NewNullableString(plan.Description.ValueStringPointer()),
		Version:                    *defectdojo.NewNullableString(plan.Version.ValueStringPointer()),
		FirstContacted:             *defectdojo.NewNullableString(plan.FirstContacted.ValueStringPointer()),
		TargetStart:                plan.TargetStart.ValueString(),
		TargetEnd:                  plan.TargetEnd.ValueString(),
		Reason:                     *defectdojo.NewNullableString(plan.Reason.ValueStringPointer()),
		Tracker:                    *defectdojo.NewNullableString(plan.Tracker.ValueStringPointer()),
		TestStrategy:               *defectdojo.NewNullableString(plan.TestStrategy.ValueStringPointer()),
		ThreatModel:                plan.ThreatModel.ValueBoolPointer(),
		ApiTest:                    plan.APITest.ValueBoolPointer(),
		PenTest:                    plan.PenTest.ValueBoolPointer(),
		CheckList:                  plan.CheckList.ValueBoolPointer(),
		Status:                     basetypesStringValueToDefectdojoNullableString(plan.Status),
		EngagementType:             basetypesStringValueToDefectdojoNullableString(plan.EngagementType),
		BuildId:                    *defectdojo.NewNullableString(plan.BuildID.ValueStringPointer()),
		CommitHash:                 *defectdojo.NewNullableString(plan.CommitHash.ValueStringPointer()),
		BranchTag:                  *defectdojo.NewNullableString(plan.BranchTag.ValueStringPointer()),
		SourceCodeManagementUri:    *defectdojo.NewNullableString(plan.SourceCodeManagementURI.ValueStringPointer()),
		DeduplicationOnEngagement:  plan.DeduplicationOnEngagement.ValueBoolPointer(),
		Lead:                       basetypesInt64ValueToDefectdojoNullableInt32(plan.Lead),
		Requester:                  basetypesInt64ValueToDefectdojoNullableInt32(plan.Requester),
		Preset:                     basetypesInt64ValueToDefectdojoNullableInt32(plan.Preset),
		ReportType:                 basetypesInt64ValueToDefectdojoNullableInt32(plan.ReportType),
		Product:                    int32(plan.Product.ValueInt64()),
		BuildServer:                basetypesInt64ValueToDefectdojoNullableInt32(plan.BuildServer),
		SourceCodeManagementServer: basetypesInt64ValueToDefectdojoNullableInt32(plan.SourceCodeManagementServer),
		OrchestrationEngine:        basetypesInt64ValueToDefectdojoNullableInt32(plan.OrchestrationEngine),
		Tags:                       tags,
	}

	// Update existing engagement
	_, res, err := r.client.EngagementsAPI.EngagementsUpdate(ctx, int32(plan.ID.ValueInt64())).EngagementRequest(engagementRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Defectdojo Engagement",
			"Could not update engagement with ID "+plan.ID.String()+": "+err.Error()+"\nDefectdojo responded with status: "+fmt.Sprintf("%v", res.Body),
		)
		return
	}

	// Get refreshed engagement value from Defectdojo
	engagement, res, err := r.client.EngagementsAPI.EngagementsRetrieve(ctx, int32(plan.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Defectdojo engagement",
			"Could not read engagement with ID "+plan.ID.String()+": "+err.Error()+"\nDefectdojo responded with status:  "+fmt.Sprintf("%v", res.Body),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(int64(engagement.GetId()))
	plan.Name = types.StringValue(engagement.GetName())
	plan.Description = types.StringValue(engagement.GetDescription())
	plan.Version = types.StringValue(engagement.GetVersion())
	plan.FirstContacted = types.StringValue(engagement.GetFirstContacted())
	plan.TargetStart = types.StringValue(engagement.GetTargetStart())
	plan.TargetEnd = types.StringValue(engagement.GetTargetEnd())
	plan.Reason = types.StringValue(engagement.GetReason())
	plan.Tracker = types.StringValue(engagement.GetTracker())
	plan.TestStrategy = types.StringValue(engagement.GetTestStrategy())
	plan.ThreatModel = types.BoolValue(engagement.GetThreatModel())
	plan.APITest = types.BoolValue(engagement.GetApiTest())
	plan.PenTest = types.BoolValue(engagement.GetPenTest())
	plan.CheckList = types.BoolValue(engagement.GetCheckList())
	plan.Status = types.StringValue(engagement.GetStatus())
	plan.EngagementType = types.StringValue(engagement.GetEngagementType())
	plan.BuildID = types.StringValue(engagement.GetBuildId())
	plan.CommitHash = types.StringValue(engagement.GetCommitHash())
	plan.BranchTag = types.StringValue(engagement.GetBranchTag())
	plan.SourceCodeManagementURI = types.StringValue(engagement.GetSourceCodeManagementUri())
	plan.DeduplicationOnEngagement = types.BoolValue(engagement.GetDeduplicationOnEngagement())
	plan.Lead = types.Int64Value(int64(engagement.GetLead()))
	plan.Requester = types.Int64Value(int64(engagement.GetRequester()))
	plan.Preset = types.Int64Value(int64(engagement.GetPreset()))
	plan.ReportType = types.Int64Value(int64(engagement.GetReportType()))
	plan.Product = types.Int64Value(int64(engagement.GetProduct()))
	plan.BuildServer = types.Int64Value(int64(engagement.GetBuildServer()))
	plan.SourceCodeManagementServer = types.Int64Value(int64(engagement.GetSourceCodeManagementServer()))
	plan.OrchestrationEngine = types.Int64Value(int64(engagement.GetOrchestrationEngine()))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.SetAttribute(ctx, path.Root("tags"), engagement.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *engagementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state engagementResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing engagement
	res, err := r.client.EngagementsAPI.EngagementsDestroy(ctx, int32(state.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Defectdojo Engagement",
			"Could not delete engagement, unexpected error: "+err.Error()+"\nDefectdojo responded with status:  "+fmt.Sprintf("%v", res.Body),
		)
		return
	}
}

func (r *engagementResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
