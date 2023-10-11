package cluster_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/metal-stack-cloud/terraform-provider-metal/internal/provider"
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
			{
				Config: testAccExampleClusterSeed,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("metal_cluster.acctest", "name", "tf-acctest"),
				),
			},
			{
				Config: testAccExampleClusterSeed + testAccExampleDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.metal_cluster.acctest-data", "kubernetes", "1.23.4"),
				),
			},
		},
	})
}

const testAccExampleClusterSeed = `
resource "metal_cluster" "acctest" {
	name = "tf-acctest"
	kubernetes = "1.23.4"
	workers = [
		{
			name = "default"
			machine_type = "c1-medium-x86"
			max_size = 2
			min_size = 1
		}
	]
}
`

const testAccExampleDataSource = `
data "metal_cluster" "acctest-data" {
	name = "tf-acctest"
}
`
