package theme

import (
	"image/color"

	fyneTheme "fyne.io/fyne/v2/theme"
)

// Theme customises launcher visuals. Nil colours use the active Fyne theme.
type Theme struct {
	Background          color.Color
	CardBackground      color.Color
	CardHoverBackground color.Color
	CardBorder          color.Color
	CardHoverBorder     color.Color
	PrimaryText         color.Color
	SecondaryText       color.Color
	MutedText           color.Color
}

// BackgroundColor returns the configured background colour or the Fyne background colour.
func (t Theme) BackgroundColor() color.Color {
	return ColorOr(t.Background, fyneTheme.Color(fyneTheme.ColorNameBackground))
}

// PrimaryTextColor returns the configured primary text colour or the Fyne foreground colour.
func (t Theme) PrimaryTextColor() color.Color {
	return ColorOr(t.PrimaryText, fyneTheme.Color(fyneTheme.ColorNameForeground))
}

// ColorOr returns value when it is present, otherwise fallback.
func ColorOr(value color.Color, fallback color.Color) color.Color {
	if value != nil {
		return value
	}

	return fallback
}
