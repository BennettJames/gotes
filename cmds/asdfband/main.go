package main

import (
	"log"

	"github.com/bennettjames/gotes/internal/iutil"
)

func main() {
	ctx, cancel := iutil.RootContext()
	defer cancel()

	runErr := runComponentRenderer(ctx)
	if runErr != nil {
		log.Print("Fatal error in render:", runErr)
	}
	return
}
