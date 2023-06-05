package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/seal-io/terraform-provider-byteset/byteset"
)

func main() {
	var debug bool

	flag.BoolVar(
		&debug,
		"debug",
		false,
		"Start provider in stand-alone debug mode.",
	)
	flag.Parse()

	err := providerserver.Serve(
		context.Background(),
		byteset.NewProvider,
		providerserver.ServeOpts{
			Address: byteset.ProviderAddress,
			Debug:   debug,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
