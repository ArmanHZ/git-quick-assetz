package ui

import (
	"os"

	"github.com/ArmanHZ/git-release-downloader/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TODO: Implement and refactor using this.
func (a *App) resetFocus(toFocus int) {}

// XXX: Maybe some parts of this function need to go to the events. But, I don't
// want to globalize more variables. We'll see.
// We could call this once and then use the reference to eliminate the modal focusables reset.
func (a *App) buildDownloadPage() *tview.Frame {
	// Reset the focusables each time, because the modal can be created multiple times.
	a.modalFocusables = []tview.Primitive{}
	a.modalFocusIndex = 0
	a.isModalActive = true

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(0, 0).
		SetGap(2, 1)

	pwd, _ := os.Getwd()
	fileSaveInput := tview.NewInputField().
		SetLabel("Save location: ").
		SetText(pwd)
	a.modalFocusables = append(a.modalFocusables, fileSaveInput)

	downloadList := tview.NewList().
		SetSelectedBackgroundColor(tcell.ColorDefault).
		SetSelectedTextColor(tcell.ColorWhite)

	urlList := [][]string{}
	for _, k := range a.downloadList {
		downloadList.AddItem(k[0], k[1], 0, nil)
		urlList = append(urlList, k)
	}

	closeModalButton := tview.NewButton("Cancel").
		SetSelectedFunc(func() {
			a.focusIndex = int(ReleaseView) // If the user clicks Cancel, most likely they'll select more files.
			a.app.SetRoot(a.mainGrid, true).SetFocus(a.focusables[a.focusIndex])
			a.isModalActive = false
		})
	a.modalFocusables = append(a.modalFocusables, closeModalButton)

	downloadButton := tview.NewButton("Download").
		SetSelectedFunc(func() {
			modal := tview.NewModal().
				SetText("Download started. Pop-up will close after the files are downloaded.")

			a.app.SetRoot(modal, false)

			go func() {
				err := utils.DownloadAssets(urlList, fileSaveInput.GetText())

				if err != nil {
					panic(err)
				}

				a.app.SetRoot(a.mainGrid, true)
				a.focusIndex = int(UrlInput)
				a.app.SetFocus(a.focusables[a.focusIndex])
				a.isModalActive = false

				a.app.Draw()
			}()

		})
	a.modalFocusables = append(a.modalFocusables, downloadButton)

	grid.AddItem(fileSaveInput, 0, 0, 1, 2, 0, 0, true).
		AddItem(downloadList, 1, 0, 1, 2, 0, 0, false).
		AddItem(closeModalButton, 2, 0, 1, 1, 0, 0, true).
		AddItem(downloadButton, 2, 1, 1, 1, 0, 0, true)

	return tview.NewFrame(grid).SetBorders(1, 1, 0, 0, 2, 2)
}

// TODO: Better name. This part will have the user input interactibles.
func (a *App) buildHeader() *tview.Frame {
	grid := tview.NewGrid().
		SetRows(0, 1, 0).
		SetColumns(0)

	a.urlInput = tview.NewInputField().
		SetLabel("GitHub URL: ").
		SetPlaceholder("https://github.com/owner/repo")
	a.focusables = append(a.focusables, a.urlInput)

	a.downloadButton = tview.NewButton("Download Assets")
	a.focusables = append(a.focusables, a.downloadButton)

	grid.AddItem(a.urlInput, 0, 0, 1, 1, 0, 0, true).
		AddItem(a.downloadButton, 2, 0, 1, 1, 0, 0, true)

	// Frame for adding padding without having bg color difference.
	return tview.NewFrame(grid).SetBorders(1, 1, 0, 0, 2, 2)
}

// TODO: Better name. This part will have the tree view.
func (a *App) buildMainBody() *tview.Frame {
	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(0)

	rootNode := tview.NewTreeNode("[::b]Releases[-]").
		SetColor(tcell.ColorBlue)
	a.releaseView = tview.NewTreeView().
		SetRoot(rootNode)
	a.focusables = append(a.focusables, a.releaseView)

	grid.AddItem(a.releaseView, 0, 0, 1, 1, 0, 0, true)

	return tview.NewFrame(grid).SetBorders(1, 1, 0, 0, 2, 2)
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
