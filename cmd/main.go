package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/xSaCh/animalia/internal/game"
)

const (
	TICKS_PER_SECOND = 0.5
)

func main() {
	milliseconds := 1000 / TICKS_PER_SECOND

	ticker := time.NewTicker(time.Millisecond * time.Duration(milliseconds))
	defer ticker.Stop()

	ctx, cancel := newCtrlCContext()
	defer cancel()
	
	world := game.NewWorld()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Printf("Tick at %v\n", time.Now())
			world.Tick()
		}
	}

}

func newCtrlCContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		fmt.Println("Ctrl+C pressed, shutting down...")
		cancel()
	}()
	return ctx, cancel
}
