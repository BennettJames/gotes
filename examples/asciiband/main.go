package main

import (
	"log"

	"github.com/gdamore/tcell"
)

func main() {
	ctx, cancel := RootContext()
	defer cancel()

	screen, screenErr := tcell.NewScreen()
	if screenErr != nil {
		log.Fatal(screenErr)
	}
	if err := screen.Init(); err != nil {
		log.Fatal(err)
	}
	defer screen.Fini()

	err := runRenderer(ctx, screen)
	if err != nil {
		log.Print("Fatal error in render:", err)
	}
	return
}
