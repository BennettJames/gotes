package main

import (
	"log"

	"github.com/bennettjames/gotes/internal/iutil"
)

func main() {
	ctx, cancel := iutil.RootContext()
	defer cancel()

	// screen, screenErr := tcell.NewScreen()
	// if screenErr != nil {
	// 	log.Fatal(screenErr)
	// }
	// if err := screen.Init(); err != nil {
	// 	log.Fatal(err)
	// }
	// defer screen.Fini()

	// runErr := runBasicRenderer(ctx, screen)
	runErr := runComponentRenderer(ctx)
	if runErr != nil {
		log.Print("Fatal error in render:", runErr)
	}
	return
}
