package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/cespare/go-tetris/tetris"
	"github.com/ensky/snake-game/snake"
	"github.com/nsf/termbox-go"
)

// games holds all available games in play.
var games = []interface {
	Start()
}{
	snake.NewGame(),
	tetris.NewGame(),
}

func main() {
	if len(os.Args) < 2 {
		usage("play: Missing target program.")
		os.Exit(1)
	}

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Env = os.Environ()

	r, w, err := os.Pipe()
	failOnErr(err)
	cmd.Stdout = w
	cmd.Stderr = w

	done := make(chan bool, 1)
	rand.Seed(time.Now().UnixNano())
	failOnErr(termbox.Init())

	go func() {
		games[rand.Intn(len(games))].Start()
		done <- true
	}()

	go func() {
		failOnErr(cmd.Run())
		w.Close()
		done <- true
	}()

	<-done
	defer termbox.Close()
	io.Copy(os.Stdout, r)
}

const usageText = `Usage: play <program> [args...]`

func usage(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Fprintln(os.Stderr, usageText)
}

func failOnErr(err error) {
	if err != nil {
		log.Fatalf("play: %s", err)
	}
}
