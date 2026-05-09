package launcher

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/adambrett/go-fyne/pkg/recent"
)

const defaultSplitOffset = 0.375 // 1 - (1/1.6)

type screen struct {
	content        *fyne.Container
	refreshRecents func()
}

func newScreen(
	recents *recent.Recent,
	window fyne.Window,
	createItem func(),
	openItem func(recent.Item),
	options options,
) *screen {
	actions := newActionsPanel(options, window, createItem, openItem)
	recentItems := newRecentsPanel(options, recents, openItem)

	return &screen{
		content:        buildContent(options, actions, recentItems),
		refreshRecents: recentItems.refresh,
	}
}

func (s *screen) canvasObject() fyne.CanvasObject {
	return s.content
}

func buildContent(options options, actions *actionsPanel, recents *recentsPanel) *fyne.Container {
	split := container.NewHSplit(
		actions.canvasObject(),
		recents.canvasObject(),
	)
	split.Offset = options.splitOffset

	return container.NewStack(
		canvas.NewRectangle(options.theme.BackgroundColor()),
		container.NewPadded(split),
	)
}
