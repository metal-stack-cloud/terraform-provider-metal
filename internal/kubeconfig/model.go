package kubeconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type kubeconfigDataSourceModel struct {
	Uuid       types.String `tfsdk:"id"`
	Project    types.String `tfsdk:"project"`
	Expiration types.String `tfsdk:"expiration"`

	Raw types.String `tfsdk:"raw"`
}
