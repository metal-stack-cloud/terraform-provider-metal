package kubeconfig

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/terraform-provider-metal/internal/session"
	"google.golang.org/protobuf/types/known/durationpb"
	"gopkg.in/yaml.v3"
)

var (
	_ datasource.DataSource              = &KubeconfigDataSource{}
	_ datasource.DataSourceWithConfigure = &KubeconfigDataSource{}
)

func NewKubeconfigDataSource() datasource.DataSource {
	return &KubeconfigDataSource{}
}

type KubeconfigDataSource struct {
	session *session.Session
}

// Metadata implements datasource.DataSource.
func (d *KubeconfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubeconfig"
}

// Schema implements datasource.DataSource.
func (d *KubeconfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         `Allows generating a new kubeconfig to be able to access and operate in the given cluster within a given time frame. If you need non-expiring access, use a ServiceAccount instead.`,
		MarkdownDescription: `Allows generating a new kubeconfig to be able to access and operate in the given cluster within a given time frame. If you need non-expiring access, use a [ServiceAccount](https://kubernetes.io/docs/concepts/security/service-accounts/) instead.`,
		Attributes:          kubeconfigDataSourceAttributes(),
	}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *KubeconfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	session, ok := req.ProviderData.(*session.Session)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *session.Session, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.session = session
}

func (d *KubeconfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data kubeconfigDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	project := data.Project.ValueString()
	if data.Project.IsNull() {
		project = d.session.Project
	}

	expiration, err := time.ParseDuration(data.Expiration.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("expiration"), "invalid duration", "must be of form `1h30m`")
		return
	}

	kcResp, err := d.session.Client.Apiv1().Cluster().GetCredentials(ctx, connect.NewRequest(&apiv1.ClusterServiceGetCredentialsRequest{
		Project:    project,
		Uuid:       data.Uuid.ValueString(),
		Expiration: durationpb.New(expiration),
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to generate kubeconfig", err.Error())
		return
	}
	tflog.Trace(ctx, "generated kubeconfig")

	data.Raw = types.StringValue(kcResp.Msg.GetKubeconfig())

	var kubeconfig rawKubeconfig
	err = yaml.Unmarshal([]byte(kcResp.Msg.GetKubeconfig()), &kubeconfig)
	if err != nil {
		resp.Diagnostics.AddAttributeWarning(path.Root("external"), "parsing raw kubeconfig failed", err.Error())
	}
	data.External = &kubeconfigContent{}

	for _, c := range kubeconfig.Clusters {
		if !strings.HasSuffix("external", c.Name) {
			continue
		}
		data.External.Host = types.StringValue(c.Cluster.Server)

		ca, err := base64.StdEncoding.DecodeString(c.Cluster.CertificateAuthorityData)
		if err != nil {
			resp.Diagnostics.AddAttributeWarning(path.Root("external").AtName("cluster_ca_certificate"), "decoding failed", err.Error())
		}
		data.External.ClusterCaCertificate = types.StringValue(string(ca))
	}

	for _, u := range kubeconfig.Users {
		if !strings.HasSuffix("external", u.Name) {
			continue
		}
		clientKey, err := base64.StdEncoding.DecodeString(u.User.ClientKeyData)
		if err != nil {
			resp.Diagnostics.AddAttributeWarning(path.Root("external").AtName("client_key"), "decoding failed", err.Error())
		}
		data.External.ClientKey = types.StringValue(string(clientKey))

		clientCert, err := base64.StdEncoding.DecodeString(u.User.ClientCertificateData)
		if err != nil {
			resp.Diagnostics.AddAttributeWarning(path.Root("external").AtName("client_certificate"), "decoding failed", err.Error())
		}
		data.External.ClientCertificate = types.StringValue(string(clientCert))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type rawKubeconfig struct {
	Clusters []struct {
		Name    string `yaml:"name"`
		Cluster struct {
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
			Server                   string `yaml:"server"`
		} `yaml:"cluster"`
	} `yaml:"clusters"`
	Users []struct {
		Name string `yaml:"name"`
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
		} `yaml:"user"`
	} `yaml:"users"`
}
