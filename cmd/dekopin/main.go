package main

import (
	"context"
	"os"

	"github.com/iwashi623/dekopin"
)

func main() {
	ctx := context.Background()
	exitCode := dekopin.Run(ctx)

	os.Exit(exitCode)
}
