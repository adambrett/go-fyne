package native

import (
	"errors"

	fyneio "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/ncruces/zenity"

	"github.com/adambrett/go-fyne/pkg/browse"
)

var _ browse.Picker = (*Picker)(nil)

// Picker uses the operating system's native file dialogs.
type Picker struct {
	parent fyneio.Window

	dialogs Dialogs
}

// New returns a picker backed by the operating system's native dialog.
func New(parent fyneio.Window) *Picker {
	return NewWithDialogs(parent, defaultDialogs())
}

// NewWithDialogs returns a picker backed by custom native dialogs.
func NewWithDialogs(parent fyneio.Window, dialogs Dialogs) *Picker {
	return &Picker{
		parent:  parent,
		dialogs: dialogs,
	}
}

// Open shows a native file or folder selection dialog.
func (p *Picker) Open(options browse.OpenOptions, onSelected func(path string)) {
	go func() {
		defer p.close(options.OnClosed)

		path, err := p.openDialog(OpenOptions{
			Title:   options.Title,
			Filters: FileFilters(options.Filters),
			Folder:  options.Folder,
		})
		if err != nil {
			p.showError(err)
			return
		}

		p.selectPath(path, onSelected)
	}()
}

// Save shows a native save-file dialog.
func (p *Picker) Save(options browse.SaveOptions, onSelected func(path string)) {
	go func() {
		defer p.close(options.OnClosed)

		path, err := p.saveDialog(SaveOptions{
			Title:            options.Title,
			Filename:         options.Filename,
			Filters:          FileFilters(options.Filters),
			ConfirmOverwrite: options.ConfirmOverwrite,
		})
		if err != nil {
			p.showError(err)
			return
		}

		p.selectPath(path, onSelected)
	}()
}

func (p *Picker) openDialog(options OpenOptions) (string, error) {
	if p.dialogs.Open != nil {
		return p.dialogs.Open(options)
	}

	return openNativeDialog(options)
}

func (p *Picker) saveDialog(options SaveOptions) (string, error) {
	if p.dialogs.Save != nil {
		return p.dialogs.Save(options)
	}

	return saveNativeDialog(options)
}

func (p *Picker) selectPath(path string, onSelected func(path string)) {
	if onSelected != nil {
		fyneio.Do(func() {
			onSelected(path)
		})
	}
}

func (p *Picker) close(onClosed func()) {
	if onClosed != nil {
		fyneio.Do(onClosed)
	}
}

func (p *Picker) showError(err error) {
	if errors.Is(err, zenity.ErrCanceled) {
		return
	}
	fyneio.Do(func() {
		dialog.ShowError(err, p.parent)
	})
}
