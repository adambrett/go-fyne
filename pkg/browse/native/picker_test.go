package native_test

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"github.com/ncruces/zenity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adambrett/go-fyne/pkg/browse"
	nativepicker "github.com/adambrett/go-fyne/pkg/browse/native"
)

func TestNew(t *testing.T) {
	// Given
	app := test.NewTempApp(t)
	window := app.NewWindow("test")

	// When
	got := nativepicker.New(window)

	// Then
	require.NotNil(t, got)
	assert.Implements(t, (*browse.Picker)(nil), got)
}

func TestPicker_Open_ReportsSelectedPath(t *testing.T) {
	// Given
	app := test.NewTempApp(t)
	window := app.NewWindow("test")
	closed := make(chan struct{}, 1)
	optionsSeen := make(chan nativepicker.OpenOptions, 1)
	selected := make(chan string, 1)

	got := nativepicker.NewWithDialogs(window, nativepicker.Dialogs{
		Open: func(options nativepicker.OpenOptions) (string, error) {
			optionsSeen <- options

			return "/tmp/document.fyne", nil
		},
	})

	// When
	got.Open(browse.OpenOptions{
		Title:    "Open Document",
		Folder:   true,
		Filters:  browse.FileFilters{{Name: "Documents", Patterns: []string{"*.fyne"}}},
		OnClosed: func() { closed <- struct{}{} },
	}, func(path string) {
		selected <- path
	})

	// Then
	assert.Equal(t, nativepicker.OpenOptions{
		Title:   "Open Document",
		Filters: nativepicker.FileFilters{{Name: "Documents", Patterns: []string{"*.fyne"}}},
		Folder:  true,
	}, <-optionsSeen)
	assert.Equal(t, "/tmp/document.fyne", receivePath(t, selected))
	receiveClosed(t, closed)
}

func TestPicker_Save_ReportsSelectedPath(t *testing.T) {
	// Given
	app := test.NewTempApp(t)
	window := app.NewWindow("test")
	closed := make(chan struct{}, 1)
	optionsSeen := make(chan nativepicker.SaveOptions, 1)
	selected := make(chan string, 1)

	got := nativepicker.NewWithDialogs(window, nativepicker.Dialogs{
		Save: func(options nativepicker.SaveOptions) (string, error) {
			optionsSeen <- options

			return "/tmp/document.fyne", nil
		},
	})

	// When
	got.Save(browse.SaveOptions{
		Title:            "Save Document",
		Filename:         "document.fyne",
		Filters:          browse.FileFilters{{Name: "Documents", Patterns: []string{"*.fyne"}}},
		ConfirmOverwrite: true,
		OnClosed:         func() { closed <- struct{}{} },
	}, func(path string) {
		selected <- path
	})

	// Then
	assert.Equal(t, nativepicker.SaveOptions{
		Title:            "Save Document",
		Filename:         "document.fyne",
		Filters:          nativepicker.FileFilters{{Name: "Documents", Patterns: []string{"*.fyne"}}},
		ConfirmOverwrite: true,
	}, <-optionsSeen)
	assert.Equal(t, "/tmp/document.fyne", receivePath(t, selected))
	receiveClosed(t, closed)
}

func TestPicker_Open_IgnoresCancel(t *testing.T) {
	// Given
	app := test.NewTempApp(t)
	window := app.NewWindow("test")
	closed := make(chan struct{}, 1)
	selected := make(chan string, 1)

	got := nativepicker.NewWithDialogs(window, nativepicker.Dialogs{
		Open: func(nativepicker.OpenOptions) (string, error) {
			return "", zenity.ErrCanceled
		},
	})

	// When
	got.Open(browse.OpenOptions{
		OnClosed: func() { closed <- struct{}{} },
	}, func(path string) {
		selected <- path
	})

	// Then
	receiveClosed(t, closed)
	assertNoPath(t, selected)
}

func TestOptions_Zenity(t *testing.T) {
	// Given
	filters := nativepicker.FileFilters{
		{Name: "Documents", Patterns: []string{"*.fyne", "*.json"}, CaseFold: true},
	}

	// When
	zenityFilters := filters.Zenity()
	openOptions := nativepicker.OpenOptions{
		Title:   "Open",
		Filters: filters,
		Folder:  true,
	}.Zenity()
	saveOptions := nativepicker.SaveOptions{
		Title:            "Save",
		Filename:         "document.fyne",
		Filters:          filters,
		ConfirmOverwrite: true,
	}.Zenity()

	// Then
	require.Len(t, zenityFilters, 1)
	assert.Equal(t, "Documents", zenityFilters[0].Name)
	assert.Equal(t, []string{"*.fyne", "*.json"}, zenityFilters[0].Patterns)
	assert.True(t, zenityFilters[0].CaseFold)
	assert.Len(t, openOptions, 3)
	assert.Len(t, saveOptions, 4)
}

func receivePath(t *testing.T, paths <-chan string) string {
	t.Helper()

	select {
	case path := <-paths:
		return path
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for selected path")

		return ""
	}
}

func receiveClosed(t *testing.T, closed <-chan struct{}) {
	t.Helper()

	select {
	case <-closed:
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for close callback")
	}
}

func assertNoPath(t *testing.T, paths <-chan string) {
	t.Helper()

	select {
	case path := <-paths:
		require.Failf(t, "unexpected selected path", "got %q", path)
	case <-time.After(50 * time.Millisecond):
	}
}
