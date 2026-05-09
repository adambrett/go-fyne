package recent_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/adambrett/go-fyne/pkg/recent"
)

func TestNewItemsFromPaths(t *testing.T) {
	// Given
	paths := []string{"/tmp/one.db", "/tmp/two.db"}

	// When
	items := recent.NewItemsFromPaths(paths)

	// Then
	assert.Equal(t, recent.Items{
		{Path: "/tmp/one.db"},
		{Path: "/tmp/two.db"},
	}, items)
}

func TestItems_Clean(t *testing.T) {
	// Given
	first := touchItemFile(t, "first.db")
	second := touchItemFile(t, "second.db")
	missing := filepath.Join(t.TempDir(), "missing.db")

	// When
	items, changed := recent.Items{
		{Path: first},
		{Name: "First", Path: first},
		{Path: missing},
		{Path: second},
	}.Clean(2, false)

	// Then
	assert.True(t, changed)
	assert.Equal(t, recent.Items{
		{Path: first},
		{Path: second},
	}, items)
}

func TestItems_CleanKeepsMissing(t *testing.T) {
	// Given
	missing := filepath.Join(t.TempDir(), "missing.db")

	// When
	items, changed := recent.Items{
		{Path: missing},
	}.Clean(5, true)

	// Then
	assert.False(t, changed)
	assert.Equal(t, recent.Items{{Path: missing}}, items)
}

func TestItems_WithoutLocation(t *testing.T) {
	// Given
	first := recent.Item{Name: "First", Path: "/tmp/first.db"}
	second := recent.Item{Name: "Second", Path: "/tmp/second.db"}
	items := recent.Items{first, second}

	// When
	filtered := items.WithoutLocation(recent.Item{Path: "/tmp/first.db"})

	// Then
	assert.Equal(t, recent.Items{second}, filtered)
	assert.Equal(t, recent.Items{first, second}, items)
}

func TestItems_WithFirst(t *testing.T) {
	// Given
	first := recent.Item{Name: "First", Path: "/tmp/first.db"}
	second := recent.Item{Name: "Second", Path: "/tmp/second.db"}
	third := recent.Item{Name: "Third", Path: "/tmp/third.db"}
	items := recent.Items{first, second}

	// When
	ordered := items.WithFirst(third, 2)
	promoted := items.WithFirst(recent.Item{Name: "Updated", Path: "/tmp/second.db"}, 3)

	// Then
	assert.Equal(t, recent.Items{third, first}, ordered)
	assert.Equal(t, recent.Items{
		{Name: "Updated", Path: "/tmp/second.db"},
		first,
	}, promoted)
}

func TestItems_SnapshotPathsContainsAndSame(t *testing.T) {
	// Given
	first := recent.Item{Path: "/tmp/first.db"}
	second := recent.Item{Path: "/tmp/second.db"}
	items := recent.Items{first, second}

	// When
	snapshot := items.Snapshot()
	snapshot[0].Path = "/tmp/changed.db"

	// Then
	assert.Equal(t, "/tmp/first.db", items[0].Path)
	assert.Equal(t, []string{"/tmp/first.db", "/tmp/second.db"}, items.Paths())
	assert.True(t, items.Contains(recent.Item{Name: "First", Path: "/tmp/first.db"}))
	assert.True(t, items.Same(recent.Items{first, second}))
	assert.False(t, items.Same(recent.Items{second, first}))
}
