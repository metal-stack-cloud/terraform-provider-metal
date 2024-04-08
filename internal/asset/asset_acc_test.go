package asset_test

import (
	"errors"
	"strconv"
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
						"data.metal_assets.assets",
						tfjsonpath.New("items"),
						knownvalue.NotNull(),
					),
				},
			},
			{
				Config: providerConfig + testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith("data.metal_assets.assets", "items.0.kubernetes.#", func(value string) error {
						count, err := strconv.Atoi(value)
						if err != nil {
							return err
						}
						if count < 3 {
							return errors.New("Retrieved too less supported Kubernetes versions")
						}
						return nil
					}),
				),
			},
			{
				Config: providerConfig + testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.region.id",
						"muc",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.region.name",
						"Munich",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.region.partitions.0.id",
						"muc-1",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.region.partitions.0.name",
						"Munic 1",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.region.defaults.machine_type",
						"n1-medium-x86",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.region.defaults.worker_max",
						"3",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.region.defaults.worker_min",
						"1",
					),
				),
			},
			{
				Config: providerConfig + testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.machine_types.0.cpus",
						"8",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.machine_types.0.id",
						"n1-medium-x86",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.machine_types.0.memory",
						"34359738368",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.machine_types.0.storage",
						"960000000000",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.machine_types.1.cpus",
						"8",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.machine_types.1.id",
						"c1-medium-x86",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.machine_types.1.memory",
						"137438953472",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.machine_types.1.storage",
						"960000000000",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.machine_types.2.cpus",
						"24",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.machine_types.2.id",
						"c1-large-x86",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.machine_types.2.memory",
						"206158430208",
					),
					resource.TestCheckResourceAttr(
						"data.metal_assets.assets",
						"items.0.machine_types.2.storage",
						"960000000000",
					),
				),
			},
		},
	})

}

const testAccExampleDataSourceConfig = `
data "metal_assets" "assets" {
}
`
