package main

import (
	"flag"
	"os"

	"github.com/geodro/lerd/internal/tray"
)

func main() {
	var mono bool
	flag.BoolVar(&mono, "mono", true, "Use monochrome (white) icon")
	flag.Parse()
	if err := tray.Run(mono); err != nil {
		os.Exit(1)
	}
}
