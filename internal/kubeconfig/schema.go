package kubeconfig

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dataschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func kubeconfigDataSourceAttributes() map[string]dataschema.Attribute {
	return map[string]dataschema.Attribute{
		"id": dataschema.StringAttribute{
			Required:    true,
			Description: "The ID of the cluster to connect to with the kubeconfig.",
		},
		"project": dataschema.StringAttribute{
			Computed:    true,
			Optional:    true,
			Description: "If the cluster is in a different project, than configured in the provider.",
		},
		"expiration": dataschema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(regexp.MustCompile(`(\d+h)?(\d+m)?`), "not a valid time duration"),
				stringvalidator.LengthAtLeast(2),
			},
			Description: "Indicates how long the kubeconfig is valid as duration string in the for of `1h02m`.",
		},
		"raw": dataschema.StringAttribute{
			Computed:    true,
			Description: "The actual kubeconfig that can be used to connect to the given cluster.",
		},
	}
}
