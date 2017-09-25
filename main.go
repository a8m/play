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
	"github.com/nsf/termbox-go"
)

func main() {
	if len(os.Args) < 2 {
		usage("play: Missing target program.")
		os.Exit(1)
	}

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Env = os.Environ()

	r1, err := cmd.StdoutPipe()
	failOnErr(err)
	r2, err := cmd.StderrPipe()
	failOnErr(err)

	done := make(chan bool, 1)
	rand.Seed(time.Now().UnixNano())
	failOnErr(termbox.Init())

	go func() { failOnErr(cmd.Run()); done <- true }()
	go func() { tetris.NewGame().Start(); done <- true }()

	<-done
	termbox.Close()
	io.Copy(os.Stdout, io.MultiReader(r1, r2))
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
