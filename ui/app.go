package ui

import (
	"grd/utils"

	// "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FocusablePrimitives int

// TODO FIXME: The order of these are kept manually for far. Need to research a good
// way of autumating this part. Specially if we add more components.
const (
	UrlInput FocusablePrimitives = iota
	DownloadButton
	ReleaseView
)

type App struct {
	app *tview.Application

	focusables    []tview.Primitive
	focusIndex    int
	downloadList  map[utils.Asset][]string // FIXME: I mean, we're only using the key value.
	isModalActive bool

	mainGrid *tview.Grid

	urlInput       *tview.InputField
	downloadButton *tview.Button
	releaseView    *tview.TreeView
}

func New() *App {
	a := &App{
		app: tview.NewApplication(),
	}

	a.buildUI()
	a.bindEvents()
	a.initInputCapture()

	// a.app.SetFocus(a.focusables[0])
	a.app.SetFocus(a.focusables[UrlInput])
	a.isModalActive = false

	return a
}

func (a *App) Run() error {
	return a.app.Run()
}
