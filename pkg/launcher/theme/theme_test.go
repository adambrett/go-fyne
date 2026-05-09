package theme_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"

	"github.com/adambrett/go-fyne/pkg/launcher/theme"
)

func TestTheme_ColorAccessors(t *testing.T) {
	// Given
	test.NewTempApp(t)
	background := color.NRGBA{R: 10, G: 20, B: 30, A: 255}
	primary := color.NRGBA{R: 40, G: 50, B: 60, A: 255}
	customTheme := theme.Theme{
		Background:  background,
		PrimaryText: primary,
	}

	// When
	gotBackground := customTheme.BackgroundColor()
	gotPrimary := customTheme.PrimaryTextColor()

	// Then
	assert.Equal(t, background, gotBackground)
	assert.Equal(t, primary, gotPrimary)
}

func TestColorOr(t *testing.T) {
	// Given
	value := color.NRGBA{R: 10, A: 255}
	fallback := color.NRGBA{B: 20, A: 255}

	// When / Then
	assert.Equal(t, value, theme.ColorOr(value, fallback))
	assert.Equal(t, fallback, theme.ColorOr(nil, fallback))
}
