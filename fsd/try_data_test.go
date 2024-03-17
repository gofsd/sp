package fsd

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcctryDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "fsd_try" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of try returned
					resource.TestCheckResourceAttr("data.fsd_try.test", "try.#", "9"),
					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttr("data.fsd_try.test", "try.0.description", ""),
					resource.TestCheckResourceAttr("data.fsd_try.test", "try.0.id", "1"),
					resource.TestCheckResourceAttr("data.fsd_try.test", "try.0.image", "/hashicorp.png"),
					resource.TestCheckResourceAttr("data.fsd_try.test", "try.0.ingredients.#", "1"),
					resource.TestCheckResourceAttr("data.fsd_try.test", "try.0.ingredients.0.id", "6"),
					resource.TestCheckResourceAttr("data.fsd_try.test", "try.0.name", "HCP Aeropress"),
					resource.TestCheckResourceAttr("data.fsd_try.test", "try.0.price", "200"),
					resource.TestCheckResourceAttr("data.fsd_try.test", "try.0.teaser", "Automation in a cup"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("data.fsd_try.test", "id", "placeholder"),
				),
			},
		},
	})
}
