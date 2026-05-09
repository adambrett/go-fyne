package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/adambrett/go-fyne/pkg/browse"
	picker "github.com/adambrett/go-fyne/pkg/browse/native"
)

const (
	appID       = "browse-native"
	windowTitle = "Native Browse Example"
)

func main() {
	fyneApp := app.NewWithID(appID)
	window := newWindow(fyneApp)

	window.ShowAndRun()
}

func newWindow(fyneApp fyne.App) fyne.Window {
	window := fyneApp.NewWindow(windowTitle)
	files := picker.New(window)
	result := widget.NewLabel("No path selected")

	window.SetContent(content(files, result))
	window.Resize(fyne.NewSize(520, 220))

	return window
}

func content(picker browse.Picker, result *widget.Label) fyne.CanvasObject {
	openFile := widget.NewButtonWithIcon("Open File", theme.FileIcon(), nil)
	openFile.OnTapped = func() {
		openFile.Disable()

		picker.Open(browse.OpenOptions{
			Title:    "Open File",
			Filters:  documentFilters(),
			OnClosed: openFile.Enable,
		}, result.SetText)
	}

	openFolder := widget.NewButtonWithIcon("Open Folder", theme.FolderOpenIcon(), nil)
	openFolder.OnTapped = func() {
		openFolder.Disable()

		picker.Open(browse.OpenOptions{
			Title:    "Open Folder",
			Folder:   true,
			OnClosed: openFolder.Enable,
		}, result.SetText)
	}

	saveFile := widget.NewButtonWithIcon("Save File", theme.DocumentSaveIcon(), nil)
	saveFile.OnTapped = func() {
		saveFile.Disable()

		picker.Save(browse.SaveOptions{
			Title:            "Save File",
			Filename:         "document.fyne",
			Filters:          documentFilters(),
			ConfirmOverwrite: true,
			OnClosed:         saveFile.Enable,
		}, result.SetText)
	}

	return container.NewPadded(container.NewVBox(
		widget.NewLabelWithStyle("Native browse", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(openFile, openFolder, saveFile),
		widget.NewSeparator(),
		result,
	))
}

func documentFilters() browse.FileFilters {
	return browse.FileFilters{
		{Name: "Documents", Patterns: []string{"*.fyne", "*.json"}, CaseFold: true},
	}
}
