package launcher

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/adambrett/go-fyne/pkg/browse"
	fynepicker "github.com/adambrett/go-fyne/pkg/browse/fyne"
	"github.com/adambrett/go-fyne/pkg/recent"
)

type actionsPanel struct {
	options    options
	content    fyne.CanvasObject
	createItem func()
	openItem   func(recent.Item)
}

func newActionsPanel(options options, window fyne.Window, createItem func(), openItem func(recent.Item)) *actionsPanel {
	if options.filePicker == nil {
		options.filePicker = fynepicker.New(window)
	}

	actions := &actionsPanel{
		options:    options,
		createItem: createItem,
		openItem:   openItem,
	}
	actions.content = actions.buildContent()

	return actions
}

func (a *actionsPanel) canvasObject() fyne.CanvasObject {
	return a.content
}

func (a *actionsPanel) buildContent() fyne.CanvasObject {
	newBtn := widget.NewButtonWithIcon(a.options.createLabel, a.options.createIcon, a.createItem)
	newBtn.Importance = widget.HighImportance

	openBtn := widget.NewButtonWithIcon(a.options.openLabel, a.options.openIcon, nil)
	openBtn.OnTapped = func() {
		a.chooseItem(openBtn)
	}

	objects := make([]fyne.CanvasObject, 0, 5)
	if logo := a.logoObject(); logo != nil {
		objects = append(objects, container.NewCenter(logo))
	}

	objects = append(objects,
		container.NewCenter(a.titleText()),
		newBtn,
		openBtn,
	)

	return container.NewPadded(container.NewCenter(container.NewVBox(objects...)))
}

func (a *actionsPanel) logoObject() fyne.CanvasObject {
	if a.options.logoObject != nil {
		return a.options.logoObject
	}

	if a.options.logo == nil {
		return nil
	}

	logo := canvas.NewImageFromResource(a.options.logo)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(a.options.logoSize)

	return logo
}

func (a *actionsPanel) titleText() *canvas.Text {
	title := canvas.NewText(a.options.title, a.options.theme.PrimaryTextColor())
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	return title
}

func (a *actionsPanel) chooseItem(trigger *widget.Button) {
	trigger.Disable()

	a.options.filePicker.Open(browse.OpenOptions{
		Title:    a.options.openLabel,
		OnClosed: trigger.Enable,
	}, func(path string) {
		if a.openItem != nil {
			a.openItem(recent.Item{Path: path})
		}
	})
}
