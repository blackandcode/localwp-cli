package main

import (
	"os"

	"github.com/blackandcode/localwp-cli/internal/localwp"
)

var version = "1.0.0"

func main() {
	os.Exit(localwp.Run(os.Args[1:], version, os.Stdin, os.Stdout, os.Stderr))
}
