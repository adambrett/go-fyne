package fyne

import (
	fyneio "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"github.com/adambrett/go-fyne/pkg/browse"
)

var _ browse.Picker = (*Picker)(nil)

// Picker uses Fyne's built-in file dialogs.
type Picker struct {
	parent fyneio.Window

	dialogs Dialogs
}

// New returns a picker backed by Fyne's built-in file dialogs.
func New(parent fyneio.Window) *Picker {
	return NewWithDialogs(parent, defaultDialogs())
}

// NewWithDialogs returns a picker backed by custom Fyne dialogs.
func NewWithDialogs(parent fyneio.Window, dialogs Dialogs) *Picker {
	return &Picker{
		parent:  parent,
		dialogs: dialogs,
	}
}

// Open shows a Fyne file or folder selection dialog.
func (p *Picker) Open(options browse.OpenOptions, onSelected func(path string)) {
	if options.Folder {
		fileDialog := p.folderOpenDialog(func(uri fyneio.ListableURI, err error) {
			if err != nil {
				p.showError(err)
				return
			}
			if uri == nil {
				return
			}
			p.selectPath(uri.Path(), onSelected)
		})

		p.show(fileDialog, options.Title, nil, options.OnClosed)
		return
	}

	fileDialog := p.fileOpenDialog(func(reader fyneio.URIReadCloser, err error) {
		if err != nil {
			p.showError(err)
			return
		}
		if reader == nil {
			return
		}
		defer func() {
			_ = reader.Close()
		}()

		p.selectPath(reader.URI().Path(), onSelected)
	})

	p.show(fileDialog, options.Title, options.Filters, options.OnClosed)
}

// Save shows a Fyne save-file dialog.
func (p *Picker) Save(options browse.SaveOptions, onSelected func(path string)) {
	fileDialog := p.fileSaveDialog(func(writer fyneio.URIWriteCloser, err error) {
		if err != nil {
			p.showError(err)
			return
		}
		if writer == nil {
			return
		}
		defer func() {
			_ = writer.Close()
		}()

		p.selectPath(writer.URI().Path(), onSelected)
	})
	if options.Filename != "" {
		fileDialog.SetFileName(options.Filename)
	}

	p.show(fileDialog, options.Title, options.Filters, options.OnClosed)
}

func (p *Picker) show(fileDialog FileDialog, title string, filters browse.FileFilters, onClosed func()) {
	if onClosed != nil {
		fileDialog.SetOnClosed(onClosed)
	}
	if title != "" {
		fileDialog.SetTitleText(title)
	}
	if len(filters) > 0 {
		fileDialog.SetFilter(filters)
	}
	fileDialog.Show()
}

func (p *Picker) fileOpenDialog(callback func(fyneio.URIReadCloser, error)) FileDialog {
	if p.dialogs.NewFileOpen != nil {
		return p.dialogs.NewFileOpen(callback, p.parent)
	}

	return newFileOpenDialog(callback, p.parent)
}

func (p *Picker) folderOpenDialog(callback func(fyneio.ListableURI, error)) FileDialog {
	if p.dialogs.NewFolderOpen != nil {
		return p.dialogs.NewFolderOpen(callback, p.parent)
	}

	return newFolderOpenDialog(callback, p.parent)
}

func (p *Picker) fileSaveDialog(callback func(fyneio.URIWriteCloser, error)) FileDialog {
	if p.dialogs.NewFileSave != nil {
		return p.dialogs.NewFileSave(callback, p.parent)
	}

	return newFileSaveDialog(callback, p.parent)
}

func (p *Picker) selectPath(path string, onSelected func(path string)) {
	if onSelected != nil {
		fyneio.Do(func() {
			onSelected(path)
		})
	}
}

func (p *Picker) showError(err error) {
	fyneio.Do(func() {
		dialog.ShowError(err, p.parent)
	})
}
