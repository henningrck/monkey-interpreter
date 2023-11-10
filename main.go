package main

import (
	"os"

	"github.com/henningrck/monkey-interpreter/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
