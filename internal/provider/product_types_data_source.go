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
	_ datasource.DataSource              = &ProductTypesDataSource{}
	_ datasource.DataSourceWithConfigure = &ProductTypesDataSource{}
)

func NewProductTypesDataSource() datasource.DataSource {
	return &ProductTypesDataSource{}
}

// ProductTypesDataSource defines the data source implementation.
type ProductTypesDataSource struct {
	client *defectdojo.APIClient
}

// ProductTypesDataSourceModel describes the data source data model.
type ProductTypesDataSourceModel struct {
	ProductTypes []productTypeModel `tfsdk:"product_types"`
}

// productTypeModel describes the data source data model.
type productTypeModel struct {
	ID              types.Int64   `tfsdk:"id"`
	Name            types.String  `tfsdk:"name"`
	Description     types.String  `tfsdk:"description"`
	CriticalProduct types.Bool    `tfsdk:"critical_product"`
	KeyProduct      types.Bool    `tfsdk:"key_product"`
	Members         []types.Int64 `tfsdk:"members"`
	Authorization   []types.Int64 `tfsdk:"authorization_groups"`
	Created         types.String  `tfsdk:"created"`
	Updated         types.String  `tfsdk:"updated"`
}

func (d *ProductTypesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_types"
}

func (d *ProductTypesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"product_types": schema.ListNestedAttribute{
				Description: "Product Types",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "The unique identifier for the product type",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the product type",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the product type",
							Computed:    true,
						},
						"critical_product": schema.BoolAttribute{
							Description: "Whether the product type is a critical product",
							Computed:    true,
						},
						"key_product": schema.BoolAttribute{
							Description: "Whether the product type is a key product",
							Computed:    true,
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
						"created": schema.StringAttribute{
							Description: "The date the product type was created",
							Computed:    true,
						},
						"updated": schema.StringAttribute{
							Description: "The date the product type was last updated",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *ProductTypesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProductTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ProductTypesDataSourceModel

	productTypes, res, err := d.client.ProductTypesAPI.ProductTypesList(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Product Types",
			"Could not red product types, unexpected error: "+err.Error()+"\nDefectdojo responded with status: "+fmt.Sprintf("%v", res.Body),
		)
		return
	}

	// Map response body to model
	for _, productType := range productTypes.Results {
		productTypeState := productTypeModel{
			ID:              types.Int64Value(int64(productType.GetId())),
			Name:            types.StringValue(productType.GetName()),
			Description:     types.StringValue(productType.GetDescription()),
			CriticalProduct: types.BoolValue(productType.GetCriticalProduct()),
			KeyProduct:      types.BoolValue(productType.GetKeyProduct()),
			Created:         types.StringValue(productType.GetCreated().String()),
			Updated:         types.StringValue(productType.GetUpdated().String()),
		}

		// for _, member := range productType.GetMembers() {
		// 	productTypeState.Members = append(productTypeState.Members, types.Int64Value(int64(member)))
		// }

		// for _, authorization := range productType.GetAuthorizationGroups() {
		// 	productTypeState.Authorization = append(productTypeState.Authorization, types.Int64Value(int64(authorization)))
		// }

		state.ProductTypes = append(state.ProductTypes, productTypeState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
