// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.uber.org/zap"

	"github.com/metal-stack-cloud/api/go/client"
	ipaddress "github.com/metal-stack-cloud/terraform-provider-metal/internal/ip_address"
	"github.com/metal-stack-cloud/terraform-provider-metal/internal/session"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &MetalstackCloudProvider{}

// MetalstackCloudProvider defines the provider implementation.
type MetalstackCloudProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string

	// TODO: Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log

	log *zap.SugaredLogger
}

// MetalstackCloudProviderModel describes the provider data model.
type MetalstackCloudProviderModel struct {
	ApiUrl       types.String `tfsdk:"api_url"`
	ApiToken     types.String `tfsdk:"api_token"`
	Organization types.String `tfsdk:"organization"`
	Project      types.String `tfsdk:"project"`
}

func (p *MetalstackCloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "metal"
	resp.Version = p.version
}

func (p *MetalstackCloudProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				MarkdownDescription: "The api_url of the metalstack.cloud API.",
				Optional:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The API token to use for authentication.",
				Optional:            true,
				Sensitive:           true,
			},
			"organization": schema.StringAttribute{
				MarkdownDescription: "The organization to use for authentication.",
				Optional:            true,
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The project to use for authentication.",
				Optional:            true,
			},
		},
	}
}

func (p *MetalstackCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MetalstackCloudProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if data.ApiUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Unknown metalstack.cloud API api_url",
			"The provider cannot create the metalstack.cloud API client as there is an unknown configuration value for the metalstack.cloud API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the METAL_STACK_CLOUD_API_URL environment variable.",
		)
	}
	if data.ApiToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Unknown metalstack.cloud API api_token",
			"The provider cannot create the metalstack.cloud API client as there is an unknown configuration value for the metalstack.cloud API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the METAL_STACK_CLOUD_API_TOKEN environment variable.",
		)
	}

	// TODO: more unknown validation

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: use viper with prefix to read from env vars and config file
	apiUrl := os.Getenv("METAL_STACK_CLOUD_API_URL")
	if !data.ApiUrl.IsNull() {
		apiUrl = data.ApiUrl.ValueString()
	}
	apiToken := os.Getenv("METAL_STACK_CLOUD_API_TOKEN")
	if !data.ApiToken.IsNull() {
		apiToken = data.ApiToken.ValueString()
	}
	project := os.Getenv("METAL_STACK_CLOUD_PROJECT")
	if !data.Project.IsNull() {
		project = data.Project.ValueString()
	}
	organization := os.Getenv("METAL_STACK_CLOUD_ORGANIZATION")
	if !data.Organization.IsNull() {
		organization = data.Organization.ValueString()
	}

	if apiUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Missing metalstack.cloud API api_url",
			"The provider cannot create the metalstack.cloud API client as there is an unknown configuration value for the metalstack.cloud API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the METAL_STACK_CLOUD_API_URL environment variable.",
		)
	}
	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing metalstack.cloud API api_token",
			"The provider cannot create the metalstack.cloud API client as there is an unknown configuration value for the metalstack.cloud API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the METAL_STACK_CLOUD_API_TOKEN environment variable.",
		)
	}
	if project == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("project"),
			"Missing metalstack.cloud API project",
			"The provider cannot create the metalstack.cloud API client as there is an unknown configuration value for the metalstack.cloud API project. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the METAL_STACK_CLOUD_PROJECT environment variable.",
		)
	}
	if organization == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("organization"),
			"Missing metalstack.cloud API organization",
			"The provider cannot create the metalstack.cloud API client as there is an unknown configuration value for the metalstack.cloud API organization. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the METAL_STACK_CLOUD_ORGANIZATION environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	dialConfig := client.DialConfig{
		BaseURL:   apiUrl,
		Token:     apiToken,
		UserAgent: "terraform-provider-metalstackcloud/" + p.version,
		Log:       p.log.Named("metalstackcloud-api"),
		Debug:     true, // TODO
	}
	apiClient := client.New(dialConfig)
	session := &session.Session{
		Client:       apiClient,
		Organization: data.Organization.ValueString(),
		Project:      data.Project.ValueString(),
	}
	resp.DataSourceData = session
	resp.ResourceData = session
}

func (p *MetalstackCloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		ipaddress.NewIpResource,
	}
}

func (p *MetalstackCloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		ipaddress.NewIpDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MetalstackCloudProvider{
			version: version,
			log:     zap.S(),
		}
	}
}
