package ui

import (
	"fmt"
	"grd/utils"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) focusNext(focusables []tview.Primitive, focusIndex *int) {
	if len(focusables) == 0 {
		return
	}

	*focusIndex = (*focusIndex + 1) % len(focusables)
	a.app.SetFocus(focusables[*focusIndex])
}

func (a *App) focusPrev(focusables []tview.Primitive, focusIndex *int) {
	if len(focusables) == 0 {
		return
	}

	*focusIndex--
	if *focusIndex < 0 {
		*focusIndex = len(focusables) - 1
	}

	a.app.SetFocus(focusables[*focusIndex])
}

func (a *App) initInputCapture() {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if !a.isModalActive {
			switch event.Key() {
			case tcell.KeyTab:
				a.focusNext(a.focusables, &a.focusIndex)
				return nil

			case tcell.KeyBacktab:
				a.focusPrev(a.focusables, &a.focusIndex)
				return nil
			}
		} else {
			switch event.Key() {
			case tcell.KeyTab:
				a.focusNext(a.modalFocusables, &a.modalFocusIndex)
				return nil

			case tcell.KeyBacktab:
				a.focusPrev(a.modalFocusables, &a.modalFocusIndex)
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
		modal := a.buildDownloadPage()
		a.app.SetRoot(modal, true).SetFocus(modal)
	})
}

func (a *App) bindEvents() {
	a.urlAction()
	a.downloadButtonAction()
}
