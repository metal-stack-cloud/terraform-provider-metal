package volume_test

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
					resource.TestCheckResourceAttr("data.metal_volume.existing", "name", "pvc-9326d0bb-6d2a-4a1f-9498-58854ad038d7"),
					resource.TestCheckResourceAttr("data.metal_volume.existing", "labels.purpose", "terraform-tests"),
				),
			},
		},
	})
}

const testAccExampleClusterSeed = `
data "metal_volume" "existing" {
	name = "pvc-9326d0bb-6d2a-4a1f-9498-58854ad038d7"
}
`
