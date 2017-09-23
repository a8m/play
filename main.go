package main

import (
	"bytes"
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
		usage("Missing target program.")
		os.Exit(1)
	}

	b := new(bytes.Buffer)
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Env = os.Environ()
	cmd.Stdout = b
	cmd.Stderr = b

	rand.Seed(time.Now().UnixNano())
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}

	go tetris.NewGame().Start()
	cmd.Run()
	termbox.Close()
	io.Copy(os.Stdout, b)
}

const usageText = `Usage: play <program> [args...]`

func usage(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Fprintln(os.Stderr, usageText)
}
