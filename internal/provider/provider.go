// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	client "github.com/metal-stack-cloud/api/go/client"
	cluster "github.com/metal-stack-cloud/terraform-provider-metal/internal/cluster"
	"github.com/metal-stack-cloud/terraform-provider-metal/internal/kubeconfig"
	ipaddress "github.com/metal-stack-cloud/terraform-provider-metal/internal/public_ip"
	session "github.com/metal-stack-cloud/terraform-provider-metal/internal/session"
	"github.com/metal-stack-cloud/terraform-provider-metal/internal/snapshot"
	"github.com/metal-stack-cloud/terraform-provider-metal/internal/volume"
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
		MarkdownDescription: "Manage bare-metal Kubernetes clusters on [metalstack.cloud](https://metalstack.cloud).\n\n" +
			"> **Note:** As [metalstack.cloud](https://metalstack.cloud) does not yet provide API Tokens, you currently need to pick your JWT from your web session. This is obviously going to change. After you logged in, open the Developer Tools of your browser, head to the console, filter for `token` and copy the JWT starting with `eyJ`.\n" +
			"> To get the project id, with dev tools open, switch to your project and open the clusters view. In the dev tools' network tab, search for `api.v1.ClusterService/List`, select a request with status 200, head to the payload and copy the project id. This is obviously going to change.\n",
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

	err := readConfigFile()
	if err != nil {
		resp.Diagnostics.AddError("Unable to read metalstack.cloud config", err.Error())
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
	if data.Organization.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("organization"),
			"Unknown metalstack.cloud API organization",
			"The provider cannot create the metalstack.cloud API client as there is an unknown configuration value for the metalstack.cloud API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the METAL_STACK_CLOUD_ORGANIZATION environment variable.",
		)
	}
	if data.Project.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("project"),
			"Unknown metalstack.cloud API project",
			"The provider cannot create the metalstack.cloud API client as there is an unknown configuration value for the metalstack.cloud API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the METAL_STACK_CLOUD_PROJECT environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiUrl := viper.GetString("api-url")
	if !data.ApiUrl.IsNull() {
		apiUrl = data.ApiUrl.ValueString()
	}
	apiToken := viper.GetString("api-token")
	if !data.ApiToken.IsNull() {
		apiToken = data.ApiToken.ValueString()
	}
	project := viper.GetString("project")
	if !data.Project.IsNull() {
		project = data.Project.ValueString()
	}
	organization := viper.GetString("organization")
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
		UserAgent: "terraform-provider-metal/" + p.version,
		Debug:     viper.GetBool("debug"),
	}
	apiClient := client.New(dialConfig)
	session := &session.Session{
		Client:       apiClient,
		Organization: organization,
		Project:      project,
	}
	resp.DataSourceData = session
	resp.ResourceData = session
}

func (p *MetalstackCloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		cluster.NewClusterResource,
		ipaddress.NewPublicIpResource,
	}
}

func (p *MetalstackCloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		cluster.NewClusterDataSource,
		ipaddress.NewPublicIpDataSource,
		volume.NewVolumeDataSource,
		snapshot.NewSnapshotDataSource,
		kubeconfig.NewKubeconfigDataSource,
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
