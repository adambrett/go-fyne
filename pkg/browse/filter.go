package browse

import (
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

// FileFilter describes a named file filter.
type FileFilter struct {
	Name     string
	Patterns []string
	CaseFold bool
}

// Matches reports whether the URI is accepted by this filter.
func (filter FileFilter) Matches(uri fyne.URI) bool {
	for _, pattern := range filter.Patterns {
		if filter.matchesPattern(uri.Path(), pattern) {
			return true
		}
	}

	return false
}

// FileFilters is a collection of file filters.
type FileFilters []FileFilter

// Matches reports whether the URI is accepted by any filter in the collection.
func (filters FileFilters) Matches(uri fyne.URI) bool {
	for _, filter := range filters {
		if filter.Matches(uri) {
			return true
		}
	}

	return len(filters) == 0
}

func (filter FileFilter) matchesPattern(path string, pattern string) bool {
	base := filepath.Base(path)
	if filter.CaseFold {
		path = strings.ToLower(path)
		base = strings.ToLower(base)
		pattern = strings.ToLower(pattern)
	}

	if ok, _ := filepath.Match(pattern, base); ok {
		return true
	}

	ok, _ := filepath.Match(pattern, path)
	return ok
}

var (
	_ storage.FileFilter = FileFilter{}
	_ storage.FileFilter = FileFilters{}
)
