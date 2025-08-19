//go:build ignore
// +build ignore

package main

import (
	"log"

	"github.com/fredrikmwold/rebasei-tui/internal/ui"
)

func main() {
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
