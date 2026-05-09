package launcher

import (
	"fyne.io/fyne/v2"
	fyneTheme "fyne.io/fyne/v2/theme"

	"github.com/adambrett/go-fyne/pkg/browse"
	"github.com/adambrett/go-fyne/pkg/launcher/theme"
	"github.com/adambrett/go-fyne/pkg/recent"
)

// Option configures a [Launcher].
type Option func(*options)

type options struct {
	recentLimit int
	keepMissing bool

	windowSize  fyne.Size
	splitOffset float64

	title           string
	createLabel     string
	openLabel       string
	recentTitle     string
	emptyRecentText string

	logo       fyne.Resource
	logoObject fyne.CanvasObject
	logoSize   fyne.Size

	createIcon fyne.Resource
	openIcon   fyne.Resource

	filePicker browse.Picker

	theme theme.Theme
}

func defaultOptions() options {
	return options{
		recentLimit:     recent.DefaultLimit,
		windowSize:      fyne.NewSize(720, 480),
		splitOffset:     defaultSplitOffset,
		title:           "Items",
		createLabel:     "New Item",
		openLabel:       "Open Item",
		recentTitle:     "Recent Items",
		emptyRecentText: "No recent items yet",
		logoSize:        fyne.NewSize(88, 88),
		createIcon:      fyneTheme.FolderNewIcon(),
		openIcon:        fyneTheme.FolderOpenIcon(),
	}
}

// WithTitle sets the headline shown above the action buttons.
func WithTitle(v string) Option {
	return func(options *options) {
		options.title = v
	}
}

// WithCreateLabel sets the primary action button label.
func WithCreateLabel(v string) Option {
	return func(options *options) {
		options.createLabel = v
	}
}

// WithOpenLabel sets the secondary action button label.
func WithOpenLabel(v string) Option {
	return func(options *options) {
		options.openLabel = v
	}
}

// WithRecentTitle sets the heading above the recent items list.
func WithRecentTitle(v string) Option {
	return func(options *options) {
		options.recentTitle = v
	}
}

// WithEmptyRecentText sets the placeholder when there are no recent items.
func WithEmptyRecentText(v string) Option {
	return func(options *options) {
		options.emptyRecentText = v
	}
}

// WithRecentLimit sets how many recent items are retained.
func WithRecentLimit(limit int) Option {
	return func(options *options) {
		options.recentLimit = limit
	}
}

// WithKeepMissingRecentItems keeps recent paths that no longer exist on disk.
func WithKeepMissingRecentItems(keep bool) Option {
	return func(options *options) {
		options.keepMissing = keep
	}
}

// WithLogo sets the image resource shown above the title.
func WithLogo(res fyne.Resource) Option {
	return func(options *options) {
		options.logo = res
	}
}

// WithLogoCanvas sets a custom canvas object instead of [WithLogo].
func WithLogoCanvas(co fyne.CanvasObject) Option {
	return func(options *options) {
		options.logoObject = co
	}
}

// WithLogoSize sets the minimum logo dimensions when using [WithLogo].
func WithLogoSize(sz fyne.Size) Option {
	return func(options *options) {
		options.logoSize = sz
	}
}

// WithCreateIcon sets the icon for the create-item button.
func WithCreateIcon(res fyne.Resource) Option {
	return func(options *options) {
		options.createIcon = res
	}
}

// WithOpenIcon sets the icon for the open-item button.
func WithOpenIcon(res fyne.Resource) Option {
	return func(options *options) {
		options.openIcon = res
	}
}

// WithFilePicker sets the file picker used by the open-item button.
func WithFilePicker(picker browse.Picker) Option {
	return func(options *options) {
		options.filePicker = picker
	}
}

// WithWindowSize sets the preferred launcher window size.
func WithWindowSize(sz fyne.Size) Option {
	return func(options *options) {
		options.windowSize = sz
	}
}

// WithSplitOffset sets the horizontal split position between actions and recents (0-1).
func WithSplitOffset(offset float64) Option {
	return func(options *options) {
		options.splitOffset = offset
	}
}

// WithTheme sets launcher-specific colours (nil fields fall back to the active app theme).
func WithTheme(t theme.Theme) Option {
	return func(options *options) {
		options.theme = t
	}
}
