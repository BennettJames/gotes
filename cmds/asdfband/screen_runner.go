package main

import (
	"context"

	"github.com/gdamore/tcell"
)

type ScreenRunner struct {
	screen tcell.Screen
}

func NewScreenRunner() (*ScreenRunner, error) {
	screen, screenErr := tcell.NewScreen()
	if screenErr != nil {
		return nil, screenErr
	}
	return &ScreenRunner{
		screen: screen,
	}, nil
}

func (sr *ScreenRunner) Run(
	ctx context.Context,
	onEv func(ev tcell.Event),
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if initErr := sr.screen.Init(); initErr != nil {
		return initErr
	}
	go func() {
		<-ctx.Done()
		sr.screen.Fini()
	}()

	for {
		ev := sr.screen.PollEvent()
		if ev != nil {
			onEv(ev)
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
