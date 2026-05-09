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

type removeButton struct {
	widget.BaseWidget

	theme theme.Theme

	hovered  bool
	onTapped func()
}

func newRemoveButton(theme theme.Theme, onTapped func()) *removeButton {
	button := &removeButton{
		theme:    theme,
		onTapped: onTapped,
	}
	button.ExtendBaseWidget(button)

	return button
}

func (b *removeButton) Tapped(_ *fyne.PointEvent) {
	if b.onTapped != nil {
		b.onTapped()
	}
}

func (b *removeButton) TappedSecondary(_ *fyne.PointEvent) {}

func (b *removeButton) MouseIn(_ *desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

func (b *removeButton) MouseMoved(_ *desktop.MouseEvent) {}

func (b *removeButton) MouseOut() {
	b.hovered = false
	b.Refresh()
}

func (b *removeButton) CreateRenderer() fyne.WidgetRenderer {
	label := canvas.NewText("x", b.textColor())
	label.TextSize = 15
	label.TextStyle = fyne.TextStyle{Bold: true}

	renderer := &removeButtonRenderer{
		button:  b,
		label:   label,
		objects: []fyne.CanvasObject{label},
	}

	renderer.Refresh()

	return renderer
}

func (b *removeButton) textColor() color.Color {
	if b.hovered {
		return theme.ColorOr(b.theme.PrimaryText, fyneTheme.Color(fyneTheme.ColorNameForeground))
	}

	return theme.ColorOr(b.theme.MutedText, fyneTheme.Color(fyneTheme.ColorNamePlaceHolder))
}

type removeButtonRenderer struct {
	button *removeButton
	label  *canvas.Text

	objects []fyne.CanvasObject
}

func (r *removeButtonRenderer) Layout(size fyne.Size) {
	labelSize := r.label.MinSize()
	r.label.Move(fyne.NewPos((size.Width-labelSize.Width)/2, (size.Height-labelSize.Height)/2))
	r.label.Resize(labelSize)
}

func (r *removeButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(28, 28)
}

func (r *removeButtonRenderer) Refresh() {
	r.label.Color = r.button.textColor()

	canvas.Refresh(r.label)
}

func (r *removeButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *removeButtonRenderer) Destroy() {}
