package launcher

import (
	"fyne.io/fyne/v2"

	"github.com/adambrett/go-fyne/pkg/recent"
)

// CreateItem starts the app's item creation flow.
// The boolean reports whether an item was created.
type CreateItem func() (recent.Item, bool)

// OpenItem opens an item chosen from the launcher.
type OpenItem func(recent.Item)

// Recents is the recent-item state used by the launcher.
type Recents interface {
	Items() recent.Items
	Add(recent.Item) bool
	Remove(recent.Item) bool
}

// Launcher is a reusable welcome screen for creating, opening, and selecting
// recent items.
type Launcher struct {
	options    options
	recent     Recents
	screen     *screen
	createItem CreateItem
	openItem   OpenItem
}

// New builds a launcher backed by caller-owned recents.
func New(
	recents Recents,
	window fyne.Window,
	createItem CreateItem,
	openItem OpenItem,
	opts ...Option,
) *Launcher {
	options := defaultOptions()

	for _, opt := range opts {
		opt(&options)
	}

	l := &Launcher{
		options:    options,
		recent:     recents,
		createItem: createItem,
		openItem:   openItem,
	}
	l.screen = newScreen(
		l.recent,
		window,
		l.createRecentItem,
		l.openRecentItem,
		options,
	)

	return l
}

// CanvasObject returns the Fyne object for this launcher.
func (l *Launcher) CanvasObject() fyne.CanvasObject {
	return l.screen.canvasObject()
}

// Size returns the preferred launcher window size.
func (l *Launcher) Size() fyne.Size {
	return l.options.windowSize
}

// RememberItem records an item so it appears in the launcher history.
func (l *Launcher) RememberItem(item recent.Item) {
	if l.recent.Add(item) {
		l.screen.refreshRecents()
	}
}

func (l *Launcher) createRecentItem() {
	if l.createItem == nil {
		return
	}

	item, ok := l.createItem()
	if ok {
		l.RememberItem(item)
	}
}

func (l *Launcher) openRecentItem(item recent.Item) {
	if l.openItem != nil {
		l.openItem(item)
	}
}
