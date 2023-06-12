package main

import (
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/seal-io/terraform-provider-byteset/byteset"
	"github.com/seal-io/terraform-provider-byteset/utils/signalx"
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
		signalx.Context(),
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
