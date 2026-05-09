package native

import "github.com/ncruces/zenity"

// Dialogs runs the operating system dialogs used by [Picker].
type Dialogs struct {
	Open func(OpenOptions) (string, error)
	Save func(SaveOptions) (string, error)
}

func defaultDialogs() Dialogs {
	return Dialogs{
		Open: openNativeDialog,
		Save: saveNativeDialog,
	}
}

func openNativeDialog(options OpenOptions) (string, error) {
	return zenity.SelectFile(options.Zenity()...)
}

func saveNativeDialog(options SaveOptions) (string, error) {
	return zenity.SelectFileSave(options.Zenity()...)
}
