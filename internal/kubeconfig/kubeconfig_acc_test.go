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
					resource.TestCheckResourceAttr("data.metal_cluster.panda", "name", "tfix-panda"),
				),
			},
			{
				Config: testAccExampleClusterSeed + testAccExampleDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.metal_kubeconfig.panda_kubeconfig", "raw"),
					resource.TestCheckResourceAttrSet("data.metal_kubeconfig.panda_kubeconfig", "external.host"),
					resource.TestCheckResourceAttrSet("data.metal_kubeconfig.panda_kubeconfig", "external.client_certificate"),
					resource.TestCheckResourceAttrSet("data.metal_kubeconfig.panda_kubeconfig", "external.client_key"),
					resource.TestCheckResourceAttrSet("data.metal_kubeconfig.panda_kubeconfig", "external.cluster_ca_certificate"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config:  testAccExampleClusterSeed,
				Destroy: true,
			},
		},
	})
}

const testAccExampleClusterSeed = `
data "metal_cluster" "panda" {
	name = "tfix-panda"
}
`

const testAccExampleDataSource = `
data "metal_kubeconfig" "panda_kubeconfig" {
	id = data.metal_cluster.panda.id
	expiration = "1h02m"
}
output "panda_kubeconfig_external" {
  value = data.metal_kubeconfig.panda_kubeconfig.external
}
`
