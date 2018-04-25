package invoke

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},

		DataSourcesMap: map[string]*schema.Resource{
			"lambda_function_invoke": dataSource(),
		},

		ResourcesMap: map[string]*schema.Resource{},
	}
}
