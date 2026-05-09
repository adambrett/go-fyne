package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/adambrett/go-fyne/pkg/launcher"
	"github.com/adambrett/go-fyne/pkg/recent"
)

const (
	appID       = "launcher"
	windowTitle = "Launcher Example"
)

func main() {
	fyneApp := app.NewWithID(appID)
	window := newWindow(fyneApp)

	window.ShowAndRun()
}

func newWindow(fyneApp fyne.App) fyne.Window {
	launcherWindow := fyneApp.NewWindow(windowTitle)

	var welcome *launcher.Launcher
	var itemWindow fyne.Window

	closeItem := func() {
		if itemWindow != nil {
			itemWindow.SetCloseIntercept(nil)
			itemWindow.Close()
		}

		itemWindow = nil
		launcherWindow.Show()
	}

	showItem := func(item recent.Item) {
		closeItem()

		itemWindow = fyneApp.NewWindow(item.DisplayName())
		itemWindow.SetContent(itemContent(item))
		itemWindow.SetCloseIntercept(closeItem)
		itemWindow.Resize(fyne.NewSize(520, 220))
		itemWindow.CenterOnScreen()
		launcherWindow.Hide()
		itemWindow.Show()
	}

	createItem := func() (recent.Item, bool) {
		item, err := createExampleItem()
		if err != nil {
			fmt.Println(err)

			return recent.Item{}, false
		}

		showItem(item)

		return item, true
	}

	openItem := func(item recent.Item) {
		welcome.RememberItem(item)
		showItem(item)
	}

	welcome = launcher.New(
		fyneApp.Preferences(),
		launcherWindow,
		createItem,
		openItem,
	)

	launcherWindow.Resize(welcome.Size())
	launcherWindow.SetContent(welcome.CanvasObject())

	return launcherWindow
}

func itemContent(item recent.Item) fyne.CanvasObject {
	return container.NewPadded(container.NewVBox(
		widget.NewLabelWithStyle(item.DisplayName(), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(item.Path),
	))
}

func createExampleItem() (recent.Item, error) {
	dir, err := os.MkdirTemp("", "go-fyne-launcher-*")
	if err != nil {
		return recent.Item{}, fmt.Errorf("create example item: %w", err)
	}

	name := fmt.Sprintf("example-item-%d", time.Now().Unix())
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(path, 0o700); err != nil {
		return recent.Item{}, fmt.Errorf("create example item: %w", err)
	}

	return recent.Item{Name: name, Path: path}, nil
}
