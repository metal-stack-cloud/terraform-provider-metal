// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"slices"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/golang-jwt/jwt/v4"
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
	ApiToken types.String `tfsdk:"api_token"`
	Project  types.String `tfsdk:"project"`
}

func (p *MetalstackCloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "metal"
	resp.Version = p.version
}

func (p *MetalstackCloudProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage bare-metal Kubernetes clusters on [metalstack.cloud](https://metalstack.cloud).\n\n" +
			"To obtain an `api token` for creating resources, visit [metalstack.cloud](https://metalstack.cloud). Head to the the `Access Tokens` section and create a new one with the desired permissions, name and validity. \n" +
			"**Note:** Watch out to first select the desired organization and project you want the token to be valid for. \n\n" +
			"All provider defaults can be derived from the environment variables `METAL_STACK_CLOUD_*` or `~/.metal-stack-cloud/config.yaml`.",
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The API token to use for authentication. Defaults to `METAL_STACK_CLOUD_API_TOKEN`.",
				Optional:            true,
				Sensitive:           true,
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The project to use for authentication. Defaults to `METAL_STACK_CLOUD_PROJECT` or derived from `api_token`.",
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
	if data.ApiToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Unknown metalstack.cloud api_token",
			"The provider cannot create the metalstack.cloud API client as there is an unknown configuration value for the metalstack.cloud API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the METAL_STACK_CLOUD_API_TOKEN environment variable.",
		)
	}
	if data.Project.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("project"),
			"Unknown metalstack.cloud project",
			"The provider cannot create the metalstack.cloud API client as there is an unknown configuration value for the metalstack.cloud API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the METAL_STACK_CLOUD_PROJECT environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiToken := viper.GetString("api-token")
	if !data.ApiToken.IsNull() {
		apiToken = data.ApiToken.ValueString()
	}
	err = assumeDefaultsFromApiToken(apiToken)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Invalid API Token",
			err.Error(),
		)
	}
	apiUrl := viper.GetString("api-url")
	project := viper.GetString("project")
	if !data.Project.IsNull() {
		project = data.Project.ValueString()
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing metalstack.cloud api_token",
			"The provider cannot create the metalstack.cloud API client as there is an unknown configuration value for the metalstack.cloud API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the METAL_STACK_CLOUD_API_TOKEN environment variable.",
		)
	}
	if project == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("project"),
			"Missing metalstack.cloud project",
			"The provider cannot create the metalstack.cloud API client as there is an unknown configuration value for the metalstack.cloud API project. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the METAL_STACK_CLOUD_PROJECT environment variable.",
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
		Client:  apiClient,
		Project: project,
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

func assumeDefaultsFromApiToken(apiToken string) error {
	parser := jwt.NewParser()

	var claims Claims
	_, _, err := parser.ParseUnverified(apiToken, &claims)
	if err != nil {
		return err
	}

	viper.Set("api-url", claims.Issuer)

	var projects []string

	subjects := make([]string, 0, len(claims.Roles)+len(claims.Permissions))
	for subject := range claims.Roles {
		subjects = append(subjects, subject)
	}
	for subject := range claims.Permissions {
		if slices.Contains(subjects, subject) {
			continue
		}
		subjects = append(subjects, subject)
	}
	for _, subject := range subjects {
		// All UUIDs are projects
		if _, err := uuid.ParseUUID(subject); err == nil {
			projects = append(projects, subject)
		}
	}
	if len(projects) == 1 {
		viper.SetDefault("project", projects[0])
	}
	return nil
}
