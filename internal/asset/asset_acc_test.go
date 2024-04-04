package asset_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/metal-stack-cloud/terraform-provider-metal/internal/provider"
)

const (
	providerConfig = `
provider "metal" {
}
	`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"metal": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
)

func TestAccAssetDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccExampleDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.metal_asset.asset",
						tfjsonpath.New("items"),
						knownvalue.NotNull(),
					),
				},
			},
			// {
			// 	Config: providerConfig + testAccExampleDataSourceConfig,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr(
			// 			"data.metal_asset.asset",
			// 			"items.#",
			// 			"2",
			// 		),
			// 	),
			// },
			// {
			// 	Config: providerConfig + testAccExampleDataSourceConfig,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr(
			// 			"data.metal_asset.asset",
			// 			"items.0.machine_types.#",
			// 			"3",
			// 		),
			// 	),
			// },
			// {
			// 	Config: providerConfig + testAccExampleDataSourceConfig,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr(
			// 			"data.metal_asset.asset",
			// 			"items.0.machine_types.0.cpu_description",
			// 			"",
			// 		),
			// 		resource.TestCheckResourceAttr(
			// 			"data.metal_asset.asset",
			// 			"items.0.machine_types.0.cpus",
			// 			"8",
			// 		),
			// 		resource.TestCheckResourceAttr(
			// 			"data.metal_asset.asset",
			// 			"items.0.machine_types.0.id",
			// 			"n1-medium-x86",
			// 		),
			// 	),
			// },
		},
	})

}

const testAccExampleDataSourceConfig = `
data "metal_asset" "asset" {
}
`
