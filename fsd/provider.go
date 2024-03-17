package fsd

import (
	"context"
	"os"

	typs "github.com/gofsd/fsd-types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &fsdProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &fsdProvider{}
}

// fsdProvider is the provider implementation.
type fsdProvider struct{}

// fsdProviderModel maps provider schema data to a Go type.
type fsdProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *fsdProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "fsd"
}

// Schema defines the provider-level schema for configuration data.
func (p *fsdProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with fsd.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for fsd API. May also be provided via fsd_HOST environment variable.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for fsd API. May also be provided via fsd_USERNAME environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for fsd API. May also be provided via fsd_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a fsd API client for data sources and resources.
func (p *fsdProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring fsd client")

	// Retrieve provider data from configuration
	var config fsdProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown fsd API Host",
			"The provider cannot create the fsd API client as there is an unknown configuration value for the fsd API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the fsd_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown fsd API Username",
			"The provider cannot create the fsd API client as there is an unknown configuration value for the fsd API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the fsd_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown fsd API Password",
			"The provider cannot create the fsd API client as there is an unknown configuration value for the fsd API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the fsd_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("fsd_HOST")
	username := os.Getenv("fsd_USERNAME")
	password := os.Getenv("fsd_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing fsd API Host",
			"The provider cannot create the fsd API client as there is a missing or empty value for the fsd API host. "+
				"Set the host value in the configuration or use the fsd_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing fsd API Username",
			"The provider cannot create the fsd API client as there is a missing or empty value for the fsd API username. "+
				"Set the username value in the configuration or use the fsd_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing fsd API Password",
			"The provider cannot create the fsd API client as there is a missing or empty value for the fsd API password. "+
				"Set the password value in the configuration or use the fsd_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "fsd_host", host)
	ctx = tflog.SetField(ctx, "fsd_username", username)
	ctx = tflog.SetField(ctx, "fsd_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "fsd_password")

	tflog.Debug(ctx, "Creating fsd client")

	// Create a new fsd client using the configuration values
	client, err := typs.NewClient(&host, &username, &password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create fsd API Client",
			"An unexpected error occurred when creating the fsd API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"fsd Client Error: "+err.Error(),
		)
		return
	}

	// Make the fsd client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
	tflog.Info(ctx, "Configured fsd client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *fsdProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCoffeesDataSource,
		NewTryDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *fsdProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOrderResource,
		NewTryResource,
	}
}
