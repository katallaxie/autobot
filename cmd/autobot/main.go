package main

import (
	"log"

	"github.com/katallaxie/autobot/cmd/autobot/cmd"
)

func main() {
	if err := cmd.RootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
