package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TODO: Better name. This part will have the user input interactibles.
func (a *App) buildHeader() *tview.Frame {
	grid := tview.NewGrid().
		SetRows(1, 0, 1, 0, 1).
		SetColumns(2, 0, 2)

	a.urlInput = tview.NewInputField().
		SetLabel("GitHub URL: ").
		SetPlaceholder("https://github.com/owner/repo").
		// TODO: Remove the debugging link
		SetText("https://github.com/peass-ng/PEASS-ng")
	a.focusables = append(a.focusables, a.urlInput)

	a.downloadButton = tview.NewButton("Download Assets")
	a.focusables = append(a.focusables, a.downloadButton)

	grid.AddItem(a.urlInput, 1, 1, 1, 1, 0, 0, true).
		AddItem(a.downloadButton, 3, 1, 1, 1, 0, 0, true)

	// Frame for adding padding without having bg color difference.
	return tview.NewFrame(grid).SetBorders(0, 0, 1, 1, 2, 2)
}

// TODO: Better name. This part will have the tree view.
func (a *App) buildMainBody() *tview.Frame {
	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(2, 0, 2)

	rootNode := tview.NewTreeNode("[::b]Releases[-]").
		SetColor(tcell.ColorBlue)
	a.releaseView = tview.NewTreeView().
		SetRoot(rootNode)
	a.focusables = append(a.focusables, a.releaseView)

	grid.AddItem(a.releaseView, 1, 1, 1, 1, 0, 0, true)

	return tview.NewFrame(grid).SetBorders(0, 0, 1, 1, 2, 2)
}

func (a *App) buildUI() {
	a.mainGrid = tview.NewGrid().
		SetRows(5, 0).
		SetColumns(0).
		SetBorders(true)

	header := a.buildHeader()
	a.mainGrid.AddItem(header, 0, 0, 1, 1, 0, 0, true)

	mainBody := a.buildMainBody()
	a.mainGrid.AddItem(mainBody, 1, 0, 1, 1, 0, 0, true)

	a.app.SetRoot(a.mainGrid, true)
}
