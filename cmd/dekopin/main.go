package main

import (
	"context"
	"log"
	"os"

	"github.com/iwashi623/dekopin"
)

func main() {
	ctx := context.Background()
	exitCode, err := dekopin.Run(ctx)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	os.Exit(exitCode)
}
