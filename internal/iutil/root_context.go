package iutil

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func RootContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	// todo (bs): consider adding a flag here for running a trace.

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-ctx.Done():
		case sig := <-c:
			log.Printf("interrupt received; shutting down (signal: %s)", sig)
		}
		cancel()

		// will force-exit and dump goroutine trace if it doesn't quickly shut down.
		// Note that this is in a goroutine - if the main routine halts while this
		// is sleeping, the program will exit before the sleep completes.
		time.Sleep(1 * time.Second)
		buf := make([]byte, 1<<20)
		stacklen := runtime.Stack(buf, true)
		log.Printf("=== received SIGQUIT ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stacklen])
		os.Exit(1)
	}()

	return ctx, cancel
}
