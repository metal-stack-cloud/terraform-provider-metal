package kubeconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type kubeconfigDataSourceModel struct {
	Uuid       types.String `tfsdk:"id"`
	Project    types.String `tfsdk:"project"`
	Expiration types.String `tfsdk:"expiration"`

	Raw      types.String       `tfsdk:"raw"`
	External *kubeconfigContent `tfsdk:"external"`
}

type kubeconfigContent struct {
	Host types.String `tfsdk:"host"`

	ClientCertificate    types.String `tfsdk:"client_certificate"`
	ClientKey            types.String `tfsdk:"client_key"`
	ClusterCaCertificate types.String `tfsdk:"cluster_ca_certificate"`
}
