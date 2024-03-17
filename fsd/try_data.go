package fsd

import (
	"context"

	typs "github.com/gofsd/fsd-types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &tryDataSource{}
	_ datasource.DataSourceWithConfigure = &tryDataSource{}
)

// NewTryDataSource is a helper function to simplify the provider implementation.
func NewTryDataSource() datasource.DataSource {
	return &tryDataSource{}
}

// tryDataSource is the data source implementation.
type tryDataSource struct {
	client *typs.Client
}

// tryDataSourceModel maps the data source schema data.
type tryDataSourceModel struct {
	try []tryModel   `tfsdk:"try"`
	ID  types.String `tfsdk:"id"`
}

// tryModel maps try schema data.
type tryModel struct {
	ID          types.Int64           `tfsdk:"id"`
	Name        types.String          `tfsdk:"name"`
	Teaser      types.String          `tfsdk:"teaser"`
	Description types.String          `tfsdk:"description"`
	Price       types.Float64         `tfsdk:"price"`
	Image       types.String          `tfsdk:"image"`
	Ingredients []tryIngredientsModel `tfsdk:"ingredients"`
}

// tryIngredientsModel maps coffee ingredients data
type tryIngredientsModel struct {
	ID types.Int64 `tfsdk:"id"`
}

// Metadata returns the data source type name.
func (d *tryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_try"
}

// Schema defines the schema for the data source.
func (d *tryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of try.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute.",
				Computed:    true,
			},
			"try": schema.ListNestedAttribute{
				Description: "List of try.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Numeric identifier of the coffee.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Product name of the coffee.",
							Computed:    true,
						},
						"teaser": schema.StringAttribute{
							Description: "Fun tagline for the coffee.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Product description of the coffee.",
							Computed:    true,
						},
						"price": schema.Float64Attribute{
							Description: "Suggested cost of the coffee.",
							Computed:    true,
						},
						"image": schema.StringAttribute{
							Description: "URI for an image of the coffee.",
							Computed:    true,
						},
						"ingredients": schema.ListNestedAttribute{
							Description: "List of ingredients in the coffee.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.Int64Attribute{
										Description: "Numeric identifier of the coffee ingredient.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *tryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tryDataSourceModel

	// try, err := d.client.GetTry()
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable to Read fsd try",
	// 		err.Error(),
	// 	)
	// 	return
	// }

	// // Map response body to model
	// for _, coffee := range try {
	// 	trytate := tryModel{
	// 		ID:          types.Int64Value(int64(coffee.ID)),
	// 		Name:        types.StringValue(coffee.Name),
	// 		Teaser:      types.StringValue(coffee.Teaser),
	// 		Description: types.StringValue(coffee.Description),
	// 		Price:       types.Float64Value(coffee.Price),
	// 		Image:       types.StringValue(coffee.Image),
	// 	}

	// 	for _, ingredient := range coffee.Ingredient {
	// 		trytate.Ingredients = append(trytate.Ingredients, tryIngredientsModel{
	// 			ID: types.Int64Value(int64(ingredient.ID)),
	// 		})
	// 	}

	// 	state.try = append(state.try, trytate)
	// }

	state.ID = types.StringValue("placeholder")

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *tryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*typs.Client)
}
