# Browse

Browse gives Fyne applications a small, testable way to ask the user for file
and folder paths.

Fyne's built-in dialogs work with URI readers and writers. Native desktop
dialogs often return plain filesystem paths. Browse puts a simple path-based
interface in front of both approaches so application code can say "open a
document" or "save this file" without caring which dialog backend is currently
being used.

The picker is deliberately small:

```go
type Picker interface {
	Open(options OpenOptions, onSelected func(path string))
	Save(options SaveOptions, onSelected func(path string))
}
```

That shape makes browse flows easy to use in the app and easy to replace in
tests.

# Installation

Browse is part of `github.com/adambrett/go-fyne`:

```sh
go get github.com/adambrett/go-fyne
```

# Basic Usage

Choose a dialog backend at the edge of your application. The rest of your code
can depend on the shared `browse.Picker` interface.

Use native dialogs:

```go
import (
	"github.com/adambrett/go-fyne/pkg/browse"
	picker "github.com/adambrett/go-fyne/pkg/browse/native"
)
```

Or use Fyne's built-in dialogs:

```go
import (
	"github.com/adambrett/go-fyne/pkg/browse"
	picker "github.com/adambrett/go-fyne/pkg/browse/fyne"
)
```

Create the picker from your window:

```go
files := picker.New(window)
```

## Opening a File

```go
files.Open(browse.OpenOptions{
	Title: "Open Document",
	Filters: browse.FileFilters{
		{Name: "Documents", Patterns: []string{"*.fyne", "*.json"}, CaseFold: true},
	},
}, func(path string) {
	openDocument(path)
})
```

The callback receives a filesystem path. If the user cancels, the callback is
not called.

## Opening a Folder

Use the same `Open` method with `Folder` set to `true`:

```go
files.Open(browse.OpenOptions{
	Title:  "Open Folder",
	Folder: true,
}, func(path string) {
	openFolder(path)
})
```

## Saving a File

```go
files.Save(browse.SaveOptions{
	Title:            "Save Document",
	Filename:         "document.fyne",
	ConfirmOverwrite: true,
	Filters: browse.FileFilters{
		{Name: "Documents", Patterns: []string{"*.fyne", "*.json"}, CaseFold: true},
	},
}, func(path string) {
	saveDocument(path)
})
```

`Filename` sets the suggested filename. `ConfirmOverwrite` asks the backend to
confirm before replacing an existing file where that behavior is supported.

# Keeping UI State in Sync

Dialogs are asynchronous. If a button opens a dialog, you usually want to
disable it while the dialog is active and enable it again when the dialog
closes.

```go
button.Disable()

files.Open(browse.OpenOptions{
	Title:    "Open Document",
	OnClosed: button.Enable,
}, onDocument)
```

`OnClosed` is called whether the user selects a path or cancels the dialog.

# File Filters

Filters are shared across both picker backends:

```go
filters := browse.FileFilters{
	{Name: "Images", Patterns: []string{"*.png", "*.jpg"}, CaseFold: true},
}
```

`CaseFold` makes matching case-insensitive. The Fyne backend converts filters
to Fyne URI filters; the native backend converts them to zenity file filters.

# Switching Backends

Both backends expose the same constructor shape:

```go
files := picker.New(window)
```

That means switching from Fyne dialogs to native dialogs is usually just an
import change:

```go
picker "github.com/adambrett/go-fyne/pkg/browse/fyne"
```

to:

```go
picker "github.com/adambrett/go-fyne/pkg/browse/native"
```

The calling code can stay focused on `browse.Picker`.

# Testing and Custom Dialogs

Most applications should use `New(window)`. When you need full control over the
dialog behavior, both backends also expose `NewWithDialogs`.

For example, a native-backed picker can be supplied with custom functions:

```go
files := nativepicker.NewWithDialogs(window, nativepicker.Dialogs{
	Open: func(options nativepicker.OpenOptions) (string, error) {
		return "/tmp/document.fyne", nil
	},
})
```

For unit tests, you can usually avoid the real picker entirely and pass a small
fake that implements `browse.Picker`:

```go
type FakePicker struct {
	Path string
}

func (p FakePicker) Open(_ browse.OpenOptions, onSelected func(string)) {
	onSelected(p.Path)
}

func (p FakePicker) Save(_ browse.SaveOptions, onSelected func(string)) {
	onSelected(p.Path)
}
```

# Examples

Run the Fyne-backed picker example with:

```sh
make run-browse-fyne
```

Run the native picker example with:

```sh
make run-browse-native
```
