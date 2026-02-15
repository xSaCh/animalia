package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/xSaCh/animalia/internal/game"
)

const (
	TICKS_PER_SECOND = 5
)

func main() {
	milliseconds := 1000 / TICKS_PER_SECOND

	ticker := time.NewTicker(time.Millisecond * time.Duration(milliseconds))
	renderTicker := time.NewTicker(time.Millisecond * time.Duration(500))
	defer ticker.Stop()
	defer renderTicker.Stop()

	ctx, cancel := newCtrlCContext()
	defer cancel()

	keyChan := make(chan string)
	// Keyboard reader goroutine
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			keyChan <- scanner.Text()
		}
	}()

	world := game.NewWorld(30)
	for i := range 1 {
		pos := world.GetRandomWalkablePosition()
		goat := game.NewGoat(i+1, pos)
		world.Entities = append(world.Entities, goat)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			world.Tick()
		case <-renderTicker.C:
			clearConsole()
			world.DrawAsciiWorld()
			world.PrintEntities()

		case key := <-keyChan:
			switch key {
			case "c":
				// Manual state changes are now handled by behavior tree
				// world.Goats[0].State = ... (states are determined by BT)
				_ = key
			}
		}
	}

}

func clearConsole() {
	fmt.Print("\033[H\033[2J")
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
