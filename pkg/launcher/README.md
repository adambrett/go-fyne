# Launcher

Launcher is a reusable welcome screen for Fyne applications that work with
path-backed items such as documents, workspaces, images, databases, saves, or
anything else your app can create and open.

It gives users the common starting points desktop apps tend to need:

- create a new item
- open an existing item
- reopen something from the recent list
- remove stale entries from that list

The launcher owns the welcome screen and recent-list UI. Your application owns
the actual create/open workflows and recent persistence setup.

# Installation

Launcher is part of `github.com/adambrett/go-fyne`:

```sh
go get github.com/adambrett/go-fyne
```

# Basic Usage

Create recents, build a launcher, then use the launcher's canvas object as the
window content:

```go
package main

import (
	"fyne.io/fyne/v2"

	"github.com/adambrett/go-fyne/pkg/launcher"
	"github.com/adambrett/go-fyne/pkg/recent"
)

func newLauncherWindow(app fyne.App) fyne.Window {
	window := app.NewWindow("Items")
	recents := recent.New(app.Preferences())

	var welcome *launcher.Launcher
	welcome = launcher.New(
		recents,
		window,
		func() (recent.Item, bool) {
			path, ok := runNewItemFlow()
			if !ok {
				return recent.Item{}, false
			}

			return recent.Item{Path: path}, true
		},
		func(item recent.Item) {
			if openItem(item.Path) == nil {
				welcome.RememberItem(item)
			}
		},
	)

	window.Resize(welcome.Size())
	window.SetContent(welcome.CanvasObject())

	return window
}
```

The callbacks are the contract between the launcher and your app:

```go
type CreateItem func() (recent.Item, bool)
type OpenItem func(recent.Item)
```

# Creating Items

`CreateItem` runs when the user taps the create button.

Return `false` when creation is cancelled or fails. Return `true` with a
`recent.Item` when creation succeeds:

```go
func createItem() (recent.Item, bool) {
	path, ok := runNewItemFlow()
	if !ok {
		return recent.Item{}, false
	}

	return recent.Item{Path: path}, true
}
```

Successful created items are remembered automatically.

# Opening Items

`OpenItem` runs when the user chooses a recent item or selects a path using the
open button:

```go
func openItem(item recent.Item) {
	if err := openItemFile(item.Path); err != nil {
		showError(err)
		return
	}

	welcome.RememberItem(item)
}
```

Calling `RememberItem` after a successful open keeps the recent list in sync
when items are opened from the file picker, another screen, a command line
argument, or any other app-specific flow.

# Recent Items

Recent items are owned by your app and passed to `launcher.New`.

```go
recents := recent.New(app.Preferences())

welcome := launcher.New(
	recents,
	window,
	createItem,
	openItem,
)
```

Each recent entry has a display name and a path:

```go
item := recent.Item{
	Name: "Example",
	Path: "/tmp/example.fyne",
}
```

If `Name` is empty, `DisplayName` uses the base name of `Path`.

Recent items are normalized to absolute paths, de-duplicated by location, and
pruned when the path no longer exists by default.

Change the limit with:

```go
recents := recent.New(app.Preferences(), recent.WithLimit(10))
```

Keep missing paths in the list with:

```go
recents := recent.New(app.Preferences(), recent.WithKeepMissing(true))
```

That is useful when items may live on removable drives, network shares, or
other locations that are temporarily unavailable.

# Opening with Native Dialogs

The open button uses Fyne's file dialog by default. To use native dialogs
instead, pass a browse picker:

```go
import nativepicker "github.com/adambrett/go-fyne/pkg/browse/native"

welcome := launcher.New(
	recents,
	window,
	createItem,
	openItem,
	launcher.WithFilePicker(nativepicker.New(window)),
)
```

Any value that implements `browse.Picker` can be supplied, which makes the
launcher easy to test or adapt to a custom search screen:

```go
type ItemPicker struct{}

func (ItemPicker) Open(options browse.OpenOptions, onSelected func(string)) {
	onSelected("/tmp/item.fyne")
}

func (ItemPicker) Save(browse.SaveOptions, func(string)) {}
```

# Customising the Screen

The default copy is intentionally generic:

- title: `Items`
- create button: `New Item`
- open button: `Open Item`
- recent title: `Recent Items`
- empty text: `No recent items yet`

Most apps should make that language specific to their own domain:

```go
welcome := launcher.New(
	recents,
	window,
	createItem,
	openItem,
	launcher.WithTitle("Workspaces"),
	launcher.WithCreateLabel("Create Workspace"),
	launcher.WithOpenLabel("Open Workspace"),
	launcher.WithRecentTitle("Recent Workspaces"),
	launcher.WithEmptyRecentText("No workspaces yet"),
)
```

Set action icons with:

```go
launcher.WithCreateIcon(myCreateIcon)
launcher.WithOpenIcon(myOpenIcon)
```

Add a logo with a Fyne resource:

```go
launcher.WithLogo(myLogoResource)
launcher.WithLogoSize(fyne.NewSize(96, 96))
```

Or provide a composed canvas object:

```go
launcher.WithLogoCanvas(canvas.NewCircle(color.White))
```

# Layout and Theme

Set the preferred window size:

```go
launcher.WithWindowSize(fyne.NewSize(900, 600))
```

Set the horizontal split between the actions and the recent list:

```go
launcher.WithSplitOffset(0.4)
```

The split offset is a value between `0` and `1`.

The launcher uses the active Fyne theme by default. Override only the launcher
colours you need:

```go
import launchertheme "github.com/adambrett/go-fyne/pkg/launcher/theme"

welcome := launcher.New(
	recents,
	window,
	createItem,
	openItem,
	launcher.WithTheme(launchertheme.Theme{
		Background:     color.NRGBA{R: 24, G: 24, B: 24, A: 255},
		PrimaryText:    color.White,
		SecondaryText:  color.NRGBA{R: 190, G: 190, B: 190, A: 255},
		CardBackground: color.NRGBA{R: 36, G: 36, B: 36, A: 255},
	}),
)
```

Nil colour fields fall back to Fyne theme colours.

# Example

Run the launcher example with:

```sh
make run-launcher
```
