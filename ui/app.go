package ui

import (
	"os"
	"path/filepath"

	"github.com/ArmanHZ/git-release-downloader/utils"

	"github.com/rivo/tview"
)

type App struct {
	app          *tview.Application
	histFilePath string

	mainFocus          *FocusManager
	repoHistoryFocus   *FocusManager
	downloadModalFocus *FocusManager

	activeFocus *FocusManager

	downloadList map[utils.Asset][]string // FIXME: I mean, we're only using the key value.

	mainGrid *tview.Grid

	urlInput       *tview.InputField
	downloadButton *tview.Button
	releaseView    *tview.TreeView
}

func initHistFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(home, ".config", "grd", "hist.json")
}

func New() *App {
	a := &App{
		app:          tview.NewApplication(),
		histFilePath: initHistFilePath(),
	}

	a.buildUI()
	a.bindEvents()
	a.initInputCapture()

	a.app.SetFocus(a.urlInput)

	return a
}

func (a *App) Run() error {
	return a.app.Run()
}
