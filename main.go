package main

import (
	"flag"
	"gioui.org/app"
	"github.com/maxazimi/qr/ui"
	"log"
	"os"
)

func main() {
	flag.Parse()

	gui := ui.New()
	go func() {
		if err := gui.Run(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
