package fsd

import (
	"context"

	typs "github.com/gofsd/fsd-types"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &tryResource{}
	_ resource.ResourceWithConfigure   = &tryResource{}
	_ resource.ResourceWithImportState = &tryResource{}
)

// tryResourceModel maps the resource schema data.
type tryResourceModel struct {
	ID          types.String   `tfsdk:"id"`
	Items       []tryItemModel `tfsdk:"items"`
	LastUpdated types.String   `tfsdk:"last_updated"`
}

// tryItemModel maps try item data.
type tryItemModel struct {
	Coffee   tryItemCoffeeModel `tfsdk:"coffee"`
	Quantity types.Int64        `tfsdk:"quantity"`
}

// tryItemCoffeeModel maps coffee try item data.
type tryItemCoffeeModel struct {
	ID          types.Int64   `tfsdk:"id"`
	Name        types.String  `tfsdk:"name"`
	Teaser      types.String  `tfsdk:"teaser"`
	Description types.String  `tfsdk:"description"`
	Price       types.Float64 `tfsdk:"price"`
	Image       types.String  `tfsdk:"image"`
}

// NewTryResource is a helper function to simplify the provider implementation.
func NewTryResource() resource.Resource {
	return &tryResource{}
}

// tryResource is the resource implementation.
type tryResource struct {
	client *typs.Client
}

// Configure adds the provider configured client to the resource.
func (r *tryResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*typs.Client)
}

// Metadata returns the resource type name.
func (r *tryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_try"
}

// Schema defines the schema for the resource.
func (r *tryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an try.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the try.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the try.",
				Computed:    true,
			},
			"items": schema.ListNestedAttribute{
				Description: "List of items in the try.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"quantity": schema.Int64Attribute{
							Description: "Count of this item in the try.",
							Required:    true,
						},
						"coffee": schema.SingleNestedAttribute{
							Description: "Coffee item in the try.",
							Required:    true,
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Description: "Numeric identifier of the coffee.",
									Required:    true,
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
							},
						},
					},
				},
			},
		},
	}
}

// Create a new resource
func (r *tryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan tryResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	// var items []typs.tryItem
	// for _, item := range plan.Items {
	// 	items = append(items, typs.tryItem{
	// 		Coffee: typs.Coffee{
	// 			ID: int(item.Coffee.ID.ValueInt64()),
	// 		},
	// 		Quantity: int(item.Quantity.ValueInt64()),
	// 	})
	// }

	// // Create new try
	// try, err := r.client.CreateTry(items)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error creating try",
	// 		"Could not create try, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }

	// Map response body to schema and populate Computed attribute values
	// plan.ID = types.StringValue(strconv.Itoa(try.ID))
	// for tryItemIndex, tryItem := range try.Items {
	// 	plan.Items[tryItemIndex] = tryItemModel{
	// 		Coffee: tryItemCoffeeModel{
	// 			ID:          types.Int64Value(int64(tryItem.Coffee.ID)),
	// 			Name:        types.StringValue(tryItem.Coffee.Name),
	// 			Teaser:      types.StringValue(tryItem.Coffee.Teaser),
	// 			Description: types.StringValue(tryItem.Coffee.Description),
	// 			Price:       types.Float64Value(tryItem.Coffee.Price),
	// 			Image:       types.StringValue(tryItem.Coffee.Image),
	// 		},
	// 		Quantity: types.Int64Value(int64(tryItem.Quantity)),
	// 	}
	// }
	// plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *tryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state tryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed try value from fsd
	// try, err := r.client.GetTry(state.ID.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Reading fsd try",
	// 		"Could not read fsd try ID "+state.ID.ValueString()+": "+err.Error(),
	// 	)
	// 	return
	// }

	// Overwrite items with refreshed state
	// state.Items = []tryItemModel{}
	// for _, item := range try.Items {
	// 	state.Items = append(state.Items, tryItemModel{
	// 		Coffee: tryItemCoffeeModel{
	// 			ID:          types.Int64Value(int64(item.Coffee.ID)),
	// 			Name:        types.StringValue(item.Coffee.Name),
	// 			Teaser:      types.StringValue(item.Coffee.Teaser),
	// 			Description: types.StringValue(item.Coffee.Description),
	// 			Price:       types.Float64Value(item.Coffee.Price),
	// 			Image:       types.StringValue(item.Coffee.Image),
	// 		},
	// 		Quantity: types.Int64Value(int64(item.Quantity)),
	// 	})
	// }

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *tryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan tryResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	// var fsdItems []typs.tryItem
	// for _, item := range plan.Items {
	// 	fsdItems = append(fsdItems, typs.tryItem{
	// 		Coffee: typs.Coffee{
	// 			ID: int(item.Coffee.ID.ValueInt64()),
	// 		},
	// 		Quantity: int(item.Quantity.ValueInt64()),
	// 	})
	// }

	// // Update existing try
	// _, err := r.client.UpdateTry(plan.ID.ValueString(), fsdItems)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Updating fsd try",
	// 		"Could not update try, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }

	// // Fetch updated items from Gettry as Updatetry items are not
	// // populated.
	// try, err := r.client.GetTry(plan.ID.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Reading fsd try",
	// 		"Could not read fsd try ID "+plan.ID.ValueString()+": "+err.Error(),
	// 	)
	// 	return
	// }

	// // Update resource state with updated items and timestamp
	// plan.Items = []tryItemModel{}
	// for _, item := range try.Items {
	// 	plan.Items = append(plan.Items, tryItemModel{
	// 		Coffee: tryItemCoffeeModel{
	// 			ID:          types.Int64Value(int64(item.Coffee.ID)),
	// 			Name:        types.StringValue(item.Coffee.Name),
	// 			Teaser:      types.StringValue(item.Coffee.Teaser),
	// 			Description: types.StringValue(item.Coffee.Description),
	// 			Price:       types.Float64Value(item.Coffee.Price),
	// 			Image:       types.StringValue(item.Coffee.Image),
	// 		},
	// 		Quantity: types.Int64Value(int64(item.Quantity)),
	// 	})
	// }
	// plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *tryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state tryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing try
	// err := r.client.DeleteTry(state.ID.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Deleting fsd try",
	// 		"Could not delete try, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }
}

func (r *tryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
