package main

import (
	"github.com/ArmanHZ/git-quick-assetz/ui"
)

func main() {
	app := ui.New()
	if err := app.Run(); err != nil {
		panic(err)
	}

}
