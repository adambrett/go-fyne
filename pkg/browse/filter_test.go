package browse_test

import (
	"testing"

	"fyne.io/fyne/v2/storage"
	"github.com/stretchr/testify/assert"

	"github.com/adambrett/go-fyne/pkg/browse"
)

func TestFileFilter_MatchesURIByPattern(t *testing.T) {
	// Given
	filter := browse.FileFilter{Name: "Documents", Patterns: []string{"*.DB"}, CaseFold: true}

	// When
	matchesDocument := filter.Matches(storage.NewFileURI("/tmp/document.db"))
	matchesText := filter.Matches(storage.NewFileURI("/tmp/document.txt"))

	// Then
	assert.True(t, matchesDocument)
	assert.False(t, matchesText)
}

func TestFileFilters_MatchesURIByPattern(t *testing.T) {
	// Given
	filters := browse.FileFilters{
		{Name: "Documents", Patterns: []string{"*.fyne"}, CaseFold: true},
	}

	// When
	matchesDocument := filters.Matches(storage.NewFileURI("/tmp/document.FYNE"))
	matchesText := filters.Matches(storage.NewFileURI("/tmp/document.txt"))

	// Then
	assert.True(t, matchesDocument)
	assert.False(t, matchesText)
}
