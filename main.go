package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/agy/terraform-provider-lambda/invoke"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: invoke.Provider,
	})
}
