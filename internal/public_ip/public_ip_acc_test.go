package ipaddress_test

import (
	"fmt"
	"strconv"
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

func TestAccPublicIPResourceAndDataSource(t *testing.T) {
	var initialIPCount int

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith("data.metal_public_ips.addrs", "items.#", func(value string) error {
						count, err := strconv.Atoi(value)
						if err != nil {
							return err
						}
						initialIPCount = count
						return nil
					}),
				),
			},
			{Config: testAccPublicIpSeedFirst},
			{
				Config: testAccPublicIpSeedFirst + testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckResourceAttrResolve("data.metal_public_ips.addrs", "items.#", func() string {
						return strconv.Itoa(initialIPCount + 1)
					}),
				),
			},
			{Config: testAccPublicIpSeedFirst + testAccPublicIpSeedSecond},
			{
				Config: testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckResourceAttrResolve("data.metal_public_ips.addrs", "items.#", func() string {
						return strconv.Itoa(initialIPCount + 2)
					}),
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

func testCheckResourceAttrResolve(name, key string, derefValue func() string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(name, key, func(value string) error {
		v := derefValue()
		if v != value {
			return fmt.Errorf(
				"%s: Attribute '%s' expected %#v, got %#v",
				name,
				key,
				value,
				v,
			)
		}
		return nil
	})
}
