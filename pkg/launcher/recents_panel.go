package launcher

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/adambrett/go-fyne/pkg/recent"
)

type recentsPanel struct {
	options  options
	content  *fyne.Container
	recents  *recent.Recent
	openItem func(recent.Item)
}

func newRecentsPanel(
	options options,
	recentState *recent.Recent,
	openItem func(recent.Item),
) *recentsPanel {
	recents := &recentsPanel{
		options:  options,
		content:  container.NewVBox(),
		recents:  recentState,
		openItem: openItem,
	}
	recents.populateItems()

	return recents
}

func (r *recentsPanel) canvasObject() fyne.CanvasObject {
	return container.NewPadded(r.content)
}

func (r *recentsPanel) refresh() {
	r.populateItems()
}

func (r *recentsPanel) populateItems() {
	title := canvas.NewText(r.options.recentTitle, r.options.theme.PrimaryTextColor())
	title.TextSize = 22
	title.TextStyle = fyne.TextStyle{Bold: true}

	r.content.Objects = nil
	r.content.Add(title)

	items := r.recents.Items()
	if len(items) == 0 {
		empty := widget.NewLabel(r.options.emptyRecentText)
		empty.Wrapping = fyne.TextWrapWord

		r.content.Add(empty)
		r.content.Refresh()

		return
	}

	for _, item := range items {
		item := item
		r.content.Add(newItemButton(item.DisplayName(), item.Path, r.options.theme, func() {
			if r.openItem != nil {
				r.openItem(item)
			}
		}, func() {
			if r.recents.Remove(item) {
				r.populateItems()
			}
		}))
	}

	r.content.Refresh()
}
