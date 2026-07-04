package ui

import (
	"os"

	"github.com/ArmanHZ/git-quick-assetz/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) displayErrorModal(err error) {
	modal := tview.NewModal().
		SetText(err.Error()).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Ok" {
				a.app.SetRoot(a.mainGrid, true)
			}
		})

	a.app.SetRoot(modal, false)
}

// XXX: Maybe some parts of this function need to go to the events. But, I don't
// want to globalize more variables. We'll see.
// We could call this once and then use the reference to eliminate the modal focusables reset.
func (a *App) buildDownloadPage() *tview.Frame {
	// Reset the focusables each time, because the modal can be created multiple times.
	// ^ now that I think about it, only the download list changes, so this could be another struct or get its components defined in the `app.go` file.
	a.downloadModalFocus = NewFocusManager()
	a.activeFocus = a.downloadModalFocus

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(0, 0).
		SetGap(2, 1)

	pwd, _ := os.Getwd()
	fileSaveInput := tview.NewInputField().
		SetLabel("Save location: ").
		SetText(pwd)
	a.activeFocus.AddFocusable(fileSaveInput)

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
			a.activeFocus = a.mainFocus
			// If the user presses cancel, they most likely want to select more releases.
			a.activeFocus.index = int(ReleaseView)
			a.app.SetRoot(a.mainGrid, true).SetFocus(a.releaseView)
		})
	a.activeFocus.AddFocusable(closeModalButton)

	downloadButton := tview.NewButton("Download").
		SetSelectedFunc(func() {
			if len(urlList) == 0 {
				return
			}

			modal := tview.NewModal().
				SetText("Download started. Pop-up will close after the files are downloaded.")

			a.app.SetRoot(modal, false)

			go func() {
				err := utils.DownloadAssets(urlList, fileSaveInput.GetText())

				if err != nil {
					a.displayErrorModal(err)
				}

				// TODO: Reset the download list and the ReleaseView colors.
				a.activeFocus = a.mainFocus
				a.activeFocus.index = int(UrlInput)
				a.app.SetRoot(a.mainGrid, true)
				a.app.SetFocus(a.urlInput)

				a.app.Draw()
			}()

		})
	a.activeFocus.AddFocusable(downloadButton)

	grid.AddItem(fileSaveInput, 0, 0, 1, 2, 0, 0, true).
		AddItem(downloadList, 1, 0, 1, 2, 0, 0, false).
		AddItem(closeModalButton, 2, 0, 1, 1, 0, 0, true).
		AddItem(downloadButton, 2, 1, 1, 1, 0, 0, true)

	return tview.NewFrame(grid).SetBorders(1, 1, 0, 0, 2, 2)
}

// TODO: This function is long and ugly.
func (a *App) buildRepoHistoryPage() *tview.Frame {
	a.repoHistoryFocus = NewFocusManager()
	a.activeFocus = a.repoHistoryFocus

	grid := tview.NewGrid().
		SetRows(1, 0).
		SetColumns(0).
		SetGap(2, 0)

	urlHistory, err := utils.LoadURLHistory(a.histFilePath)
	if err != nil {
		panic(err)
	}

	urlList := tview.NewList()
	urlList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyRune && event.Rune() == 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case event.Key() == tcell.KeyRune && event.Rune() == 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		}

		return event
	})
	urlList.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		if r != 'q' {
			a.urlInput.SetText(s1)
			a.activeFocus = a.mainFocus
			a.app.SetRoot(a.mainGrid, true)
		}
	})

	for _, v := range urlHistory {
		urlList.AddItem(v, "", 0, nil)
	}

	urlList.AddItem("Back", "Press to return to the main page", 'q', func() {
		a.activeFocus = a.mainFocus
		a.app.SetRoot(a.mainGrid, true)
	})
	a.activeFocus.AddFocusable(urlList)

	grid.AddItem(tview.NewTextView().SetDynamicColors(true).SetText("URL History [gray]Press q to return[white]"), 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(urlList, 1, 0, 1, 1, 0, 0, true)

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
	a.activeFocus.AddFocusable(a.urlInput)

	a.downloadButton = tview.NewButton("Download Assets")
	a.activeFocus.AddFocusable(a.downloadButton)

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
	a.activeFocus.AddFocusable(a.releaseView)

	grid.AddItem(a.releaseView, 0, 0, 1, 1, 0, 0, true)

	return tview.NewFrame(grid).SetBorders(1, 1, 0, 0, 2, 2)
}

func (a *App) buildUI() {
	a.mainGrid = tview.NewGrid().
		SetRows(5, 0).
		SetColumns(0).
		SetBorders(true)

	a.mainFocus = NewFocusManager()
	a.activeFocus = a.mainFocus

	header := a.buildHeader()
	a.mainGrid.AddItem(header, 0, 0, 1, 1, 0, 0, true)

	mainBody := a.buildMainBody()
	a.mainGrid.AddItem(mainBody, 1, 0, 1, 1, 0, 0, true)

	a.app.SetRoot(a.mainGrid, true)
}
