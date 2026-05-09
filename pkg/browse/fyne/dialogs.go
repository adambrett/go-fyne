package fyne

import (
	fyneio "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

// FileDialog is the small part of Fyne's dialog API used by [Picker].
type FileDialog interface {
	SetFileName(string)
	SetFilter(storage.FileFilter)
	SetOnClosed(func())
	SetTitleText(string)
	Show()
}

// Dialogs creates the concrete dialogs used by [Picker].
type Dialogs struct {
	NewFileOpen   func(func(fyneio.URIReadCloser, error), fyneio.Window) FileDialog
	NewFolderOpen func(func(fyneio.ListableURI, error), fyneio.Window) FileDialog
	NewFileSave   func(func(fyneio.URIWriteCloser, error), fyneio.Window) FileDialog
}

func defaultDialogs() Dialogs {
	return Dialogs{
		NewFileOpen:   newFileOpenDialog,
		NewFolderOpen: newFolderOpenDialog,
		NewFileSave:   newFileSaveDialog,
	}
}

func newFileOpenDialog(callback func(fyneio.URIReadCloser, error), parent fyneio.Window) FileDialog {
	return dialog.NewFileOpen(callback, parent)
}

func newFolderOpenDialog(callback func(fyneio.ListableURI, error), parent fyneio.Window) FileDialog {
	return dialog.NewFolderOpen(callback, parent)
}

func newFileSaveDialog(callback func(fyneio.URIWriteCloser, error), parent fyneio.Window) FileDialog {
	return dialog.NewFileSave(callback, parent)
}
