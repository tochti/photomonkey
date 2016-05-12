package main

import (
	"log"

	"github.com/tochti/photomonkey/bot"
)

func main() {

	specs, err := bot.ReadSpecs()
	if err != nil {
		log.Fatal(err)
	}

	bot.Start(specs.Token, specs.ImageDir)
}
