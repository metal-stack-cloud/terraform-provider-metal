package kubeconfig_test

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

func TestAccKubeconfigDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleClusterSeed,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("metal_cluster.panda", "name", "tf-kubeconf"),
				),
			},
			{
				Config: testAccExampleClusterSeed + testAccExampleDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.metal_kubeconfig.panda_kubeconfig", "raw"),
					resource.TestCheckResourceAttrSet("data.metal_kubeconfig.panda_kubeconfig", "external"),
					resource.TestCheckResourceAttrSet("data.metal_kubeconfig.panda_kubeconfig.external", "host"),
					resource.TestCheckResourceAttrSet("data.metal_kubeconfig.panda_kubeconfig.external", "client_certificate"),
					resource.TestCheckResourceAttrSet("data.metal_kubeconfig.panda_kubeconfig.external", "client_key"),
					resource.TestCheckResourceAttrSet("data.metal_kubeconfig.panda_kubeconfig.external", "cluster_ca_certificate"),
				),
			},
		},
	})
}

const testAccExampleClusterSeed = `
resource "metal_cluster" "panda" {
	name = "tf-kubeconf"
	kubernetes = "1.27.9"
	workers = [
		{
			name = "group-0"
			machine_type = "c1-medium-x86"
			max_size = 2
			min_size = 1
		}
	]
	// FIXME: https://github.com/metal-stack-cloud/terraform-provider-metal/issues/51
	// maintenance = {
	// 	time_window = {
	// 		begin = "05:00 AM"
	// 		duration = 2
	// 	}
	// }
}
`

const testAccExampleDataSource = `
data "metal_kubeconfig" "panda_kubeconfig" {
	id = resource.metal_cluster.panda.id
	expiration = "1h02m"
}
output "panda_kubeconfig" {
  value = data.metal_kubeconfig.panda_kubeconfig
}
`
