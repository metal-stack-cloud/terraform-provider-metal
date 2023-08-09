package ipaddress_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

func TestAccExampleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.metal_public_ips.addrs", "items.#", "0"),
				),
			},
			{Config: testAccPublicIpSeedFirst},
			{
				Config: testAccPublicIpSeedFirst + testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.metal_public_ips.addrs", "items.#", "1"),
				),
			},
			{Config: testAccPublicIpSeedFirst + testAccPublicIpSeedSecond},
			{
				Config: testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.metal_public_ips.addrs", "items.#", "2"),
				),
			},
		},
	})
}

const testAccExampleDataSourceConfig = `
data "metal_public_ips" "addrs" {
}
`

const testAccPublicIpSeedFirst = `
resource "metal_public_ip" "first_ip" {
	name = "first"
}
	`
const testAccPublicIpSeedSecond = `
resource "metal_public_ip" "second_ip" {
	name = "second"
	description = "My description"
}
	`
