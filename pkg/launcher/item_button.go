package launcher

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	fyneTheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/adambrett/go-fyne/pkg/launcher/theme"
)

type itemButton struct {
	widget.BaseWidget

	title  string
	detail string

	theme theme.Theme

	hovered bool

	onTapped func()
	onClose  func()
}

func newItemButton(title string, detail string, theme theme.Theme, onTapped func(), onClose func()) *itemButton {
	button := &itemButton{
		title:    title,
		detail:   detail,
		theme:    theme,
		onTapped: onTapped,
		onClose:  onClose,
	}
	button.ExtendBaseWidget(button)

	return button
}

func (b *itemButton) Tapped(_ *fyne.PointEvent) {
	if b.onTapped != nil {
		b.onTapped()
	}
}

func (b *itemButton) TappedSecondary(_ *fyne.PointEvent) {}

func (b *itemButton) MouseIn(_ *desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

func (b *itemButton) MouseMoved(_ *desktop.MouseEvent) {}

func (b *itemButton) MouseOut() {
	b.hovered = false
	b.Refresh()
}

func (b *itemButton) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(b.cardBackground())
	bg.CornerRadius = 8
	bg.StrokeWidth = 1

	title := canvas.NewText(b.title, b.primaryText())
	title.TextSize = 16
	title.TextStyle = fyne.TextStyle{Bold: true}

	detail := canvas.NewText(compactText(b.detail), b.secondaryText())
	detail.TextSize = 11
	detail.TextStyle = fyne.TextStyle{Monospace: true}

	close := newRemoveButton(b.theme, func() {
		if b.onClose != nil {
			b.onClose()
		}
	})

	r := &itemButtonRenderer{
		button:  b,
		bg:      bg,
		title:   title,
		detail:  detail,
		close:   close,
		objects: []fyne.CanvasObject{bg, title, detail, close},
	}
	r.Refresh()

	return r
}

func (b *itemButton) cardBackground() color.Color {
	if b.hovered {
		return theme.ColorOr(b.theme.CardHoverBackground, fyneTheme.Color(fyneTheme.ColorNameHover))
	}

	return theme.ColorOr(b.theme.CardBackground, fyneTheme.Color(fyneTheme.ColorNameButton))
}

func (b *itemButton) cardBorder() color.Color {
	if b.hovered {
		return theme.ColorOr(b.theme.CardHoverBorder, fyneTheme.Color(fyneTheme.ColorNameFocus))
	}

	return theme.ColorOr(b.theme.CardBorder, fyneTheme.Color(fyneTheme.ColorNamePlaceHolder))
}

func (b *itemButton) primaryText() color.Color {
	return theme.ColorOr(b.theme.PrimaryText, fyneTheme.Color(fyneTheme.ColorNameForeground))
}

func (b *itemButton) secondaryText() color.Color {
	return theme.ColorOr(b.theme.SecondaryText, fyneTheme.Color(fyneTheme.ColorNamePlaceHolder))
}

type itemButtonRenderer struct {
	button *itemButton

	bg     *canvas.Rectangle
	title  *canvas.Text
	detail *canvas.Text
	close  *removeButton

	objects []fyne.CanvasObject
}

func (r *itemButtonRenderer) Layout(size fyne.Size) {
	const (
		inset       = float32(14)
		buttonSize  = float32(28)
		buttonInset = float32(4)
	)

	r.bg.Move(fyne.NewPos(0, 0))
	r.bg.Resize(size)

	textWidth := size.Width - inset*2 - buttonSize - buttonInset
	r.title.Move(fyne.NewPos(inset, 11))
	r.title.Resize(fyne.NewSize(textWidth, r.title.MinSize().Height))

	r.detail.Move(fyne.NewPos(inset, 39))
	r.detail.Resize(fyne.NewSize(textWidth, r.detail.MinSize().Height))

	r.close.Move(fyne.NewPos(size.Width-buttonSize-buttonInset, buttonInset))
	r.close.Resize(fyne.NewSize(buttonSize, buttonSize))
}

func (r *itemButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(280, 68)
}

func (r *itemButtonRenderer) Refresh() {
	r.bg.FillColor = r.button.cardBackground()
	r.bg.StrokeColor = r.button.cardBorder()

	r.title.Text = r.button.title
	r.title.Color = r.button.primaryText()

	r.detail.Text = compactText(r.button.detail)
	r.detail.Color = r.button.secondaryText()

	canvas.Refresh(r.bg)
	canvas.Refresh(r.title)
	canvas.Refresh(r.detail)

	r.close.Refresh()
}

func (r *itemButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *itemButtonRenderer) Destroy() {}

func compactText(text string) string {
	const maxRunes = 52

	runes := []rune(text)
	if len(runes) <= maxRunes {
		return text
	}

	head := maxRunes/2 - 2
	tail := maxRunes - head - 3

	return string(runes[:head]) + "..." + string(runes[len(runes)-tail:])
}
