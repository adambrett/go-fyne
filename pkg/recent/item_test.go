package recent_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/adambrett/go-fyne/pkg/recent"
)

func TestItem_DisplayName(t *testing.T) {
	// Given
	path := filepath.Join(t.TempDir(), "item.db")

	tests := []struct {
		name string
		item recent.Item
		want string
	}{
		{name: "uses name", item: recent.Item{Name: "Named", Path: path}, want: "Named"},
		{name: "uses path base", item: recent.Item{Path: path}, want: "item.db"},
		{name: "empty item", item: recent.Item{}, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			got := tt.item.DisplayName()

			// Then
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestItem_SameLocation(t *testing.T) {
	// Given
	path := filepath.Join(t.TempDir(), "item.db")

	// When
	samePath := recent.Item{Path: path}.SameLocation(recent.Item{Name: "Other", Path: path})
	differentPath := recent.Item{Path: path}.SameLocation(recent.Item{Path: path + ".old"})
	sameName := recent.Item{Name: "Untitled"}.SameLocation(recent.Item{Name: "Untitled"})

	// Then
	assert.True(t, samePath)
	assert.False(t, differentPath)
	assert.True(t, sameName)
}
