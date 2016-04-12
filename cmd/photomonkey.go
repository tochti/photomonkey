package main

import (
	"log"

	"github.com/tochti/photomonkey/libs"
)

func main() {

	specs, err := photomonkey.ReadSpecs()
	if err != nil {
		log.Fatal(err)
	}

	photomonkey.Start(specs.Token)
}
