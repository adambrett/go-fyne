package native

import (
	"github.com/ncruces/zenity"

	"github.com/adambrett/go-fyne/pkg/browse"
)

// OpenOptions describes an operating-system open dialog request.
type OpenOptions struct {
	Title   string
	Filters FileFilters
	Folder  bool
}

// SaveOptions describes an operating-system save dialog request.
type SaveOptions struct {
	Title            string
	Filename         string
	Filters          FileFilters
	ConfirmOverwrite bool
}

// FileFilters converts browse filters into the shape expected by this backend.
type FileFilters browse.FileFilters

// Zenity returns these filters in the shape expected by zenity.
func (filters FileFilters) Zenity() []zenity.FileFilter {
	zenityFilters := make([]zenity.FileFilter, 0, len(filters))
	for _, filter := range filters {
		zenityFilters = append(zenityFilters, zenity.FileFilter{
			Name:     filter.Name,
			Patterns: filter.Patterns,
			CaseFold: filter.CaseFold,
		})
	}

	return zenityFilters
}

// Zenity returns these options in the shape expected by zenity.
func (options OpenOptions) Zenity() []zenity.Option {
	zenityOptions := []zenity.Option{zenity.Title(options.Title)}
	if options.Folder {
		zenityOptions = append(zenityOptions, zenity.Directory())
	}
	if len(options.Filters) > 0 {
		zenityOptions = append(zenityOptions, zenity.FileFilters(options.Filters.Zenity()))
	}

	return zenityOptions
}

// Zenity returns these options in the shape expected by zenity.
func (options SaveOptions) Zenity() []zenity.Option {
	zenityOptions := []zenity.Option{
		zenity.Title(options.Title),
		zenity.Filename(options.Filename),
	}
	if len(options.Filters) > 0 {
		zenityOptions = append(zenityOptions, zenity.FileFilters(options.Filters.Zenity()))
	}
	if options.ConfirmOverwrite {
		zenityOptions = append(zenityOptions, zenity.ConfirmOverwrite())
	}

	return zenityOptions
}
