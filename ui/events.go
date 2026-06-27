package ui

import (
	"fmt"
	"grd/utils"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) focusNext() {
	if len(a.focusables) == 0 {
		return
	}

	a.focusIndex = (a.focusIndex + 1) % len(a.focusables)
	a.app.SetFocus(a.focusables[a.focusIndex])
}

func (a *App) focusPrev() {
	if len(a.focusables) == 0 {
		return
	}

	a.focusIndex--
	if a.focusIndex < 0 {
		a.focusIndex = len(a.focusables) - 1
	}

	a.app.SetFocus(a.focusables[a.focusIndex])
}

func (a *App) initInputCapture() {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if !a.isModalActive {
			switch event.Key() {
			case tcell.KeyTab:
				a.focusNext()
				return nil

			case tcell.KeyBacktab:
				a.focusPrev()
				return nil
			}
		}

		return event
	})
}

func (a *App) resetReleaseTree() {
	a.releaseView.GetRoot().ClearChildren()
	a.releaseView.SetCurrentNode(a.releaseView.GetRoot())
	a.downloadList = map[utils.Asset][]string{}
}

func (a *App) populateReleaseTree(releases []utils.Release) {
	a.resetReleaseTree()

	root := a.releaseView.GetRoot()

	for _, release := range releases {
		rel := release

		// XXX: Published date or update date? idk
		t, err := time.Parse(time.RFC3339, rel.PublishedAt)
		if err != nil {
			log.Fatal(err)
		}

		formatted := t.Format("2006-01-02")

		releaseNodeText := fmt.Sprintf("%s [gray](%s)[-]", rel.TagName, formatted)
		releaseNode := tview.NewTreeNode(releaseNodeText).
			SetReference(rel)

		root.AddChild(releaseNode)
	}

	a.releaseView.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()

		if ref == nil {
			return
		}

		// TODO: More readable name
		switch v := ref.(type) {
		case utils.Release:
			// Only populate once
			if len(node.GetChildren()) == 0 {
				space := utils.AssetDigestSpaceCalc(v.Assets)
				for _, asset := range v.Assets {
					asset := asset
					assetNodeText := fmt.Sprintf("%-*s [gray]%s[-]", space, asset.Name, asset.Digest)
					assetNode := tview.NewTreeNode(assetNodeText).
						SetReference(asset)
					node.AddChild(assetNode)
				}
				node.SetExpanded(true)
			} else {
				node.SetExpanded(!node.IsExpanded())
			}

		// Asset selected
		case utils.Asset:
			if _, ok := a.downloadList[v]; ok {
				delete(a.downloadList, v)
				node.SetColor(tcell.ColorWhite)
			} else {
				a.downloadList[v] = []string{v.Name, v.BrowserDownloadURL}
				node.SetColor(tcell.ColorYellow)
			}
		default:
			node.SetExpanded(!node.IsExpanded())
		}
	})

}

func (a *App) urlAction() {
	a.urlInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			releases, err := utils.GetReleases(a.urlInput.GetText())
			if err != nil {
				// TODO: Pop-up saying error or something.
			}

			a.populateReleaseTree(releases)
			// XXX
			a.focusIndex = int(ReleaseView)
			a.app.SetFocus(a.focusables[a.focusIndex])
		}
	})

}

func (a *App) downloadButtonAction() {
	a.downloadButton.SetSelectedFunc(func() {
		a.isModalActive = true

		var text string
		if len(a.downloadList) == 0 {
			text = "No assets selected."
		} else {
			text = "Selected assets:\n\n"

			for _, asset := range a.downloadList {
				text += fmt.Sprintf(
					"%s\n%s\n\n",
					asset[0],
					asset[1],
				)
			}
		}

		modal := tview.NewModal().
			SetText(text).
			AddButtons([]string{"Download", "Cancel"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Cancel" {
					a.app.SetRoot(a.mainGrid, true).SetFocus(a.focusables[a.focusIndex])
					a.isModalActive = false
				}
			})
		a.app.SetRoot(modal, false).SetFocus(modal)
	})
}

func (a *App) bindEvents() {
	a.urlAction()
	a.downloadButtonAction()
}
