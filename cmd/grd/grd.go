package main

import (
	"github.com/ArmanHZ/git-release-downloader/ui"
)

func main() {
	app := ui.New()
	if err := app.Run(); err != nil {
		panic(err)
	}

}
