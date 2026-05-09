package fyne_test

import (
	"errors"
	"io"
	"testing"
	"time"

	fyneio "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adambrett/go-fyne/pkg/browse"
	fynepicker "github.com/adambrett/go-fyne/pkg/browse/fyne"
)

func TestNew(t *testing.T) {
	// Given
	app := test.NewTempApp(t)
	window := app.NewWindow("test")

	// When
	got := fynepicker.New(window)

	// Then
	require.NotNil(t, got)
	assert.Implements(t, (*browse.Picker)(nil), got)
}

func TestPicker_Open_ReportsSelectedFile(t *testing.T) {
	// Given
	app := test.NewTempApp(t)
	window := app.NewWindow("test")
	closed := false
	selected := make(chan string, 1)
	reader := &fakeReadCloser{uri: storage.NewFileURI("/tmp/document.fyne")}
	dialog := &fakeFileDialog{}

	got := fynepicker.NewWithDialogs(window, fynepicker.Dialogs{
		NewFileOpen: func(callback func(fyneio.URIReadCloser, error), parent fyneio.Window) fynepicker.FileDialog {
			assert.Same(t, window, parent)
			callback(reader, nil)

			return dialog
		},
	})

	// When
	got.Open(browse.OpenOptions{
		Title:    "Open Document",
		Filters:  browse.FileFilters{{Name: "Documents", Patterns: []string{"*.fyne"}}},
		OnClosed: func() { closed = true },
	}, func(path string) {
		selected <- path
	})
	fyneio.DoAndWait(func() {})

	// Then
	assert.Equal(t, "/tmp/document.fyne", receivePath(t, selected))
	assert.True(t, reader.closed)
	assert.Equal(t, "Open Document", dialog.title)
	assert.NotNil(t, dialog.filter)
	assert.True(t, dialog.shown)
	assert.True(t, closed)
}

func TestPicker_Open_ReportsSelectedFolder(t *testing.T) {
	// Given
	app := test.NewTempApp(t)
	window := app.NewWindow("test")
	selected := make(chan string, 1)
	dialog := &fakeFileDialog{}

	got := fynepicker.NewWithDialogs(window, fynepicker.Dialogs{
		NewFolderOpen: func(callback func(fyneio.ListableURI, error), parent fyneio.Window) fynepicker.FileDialog {
			assert.Same(t, window, parent)
			callback(fakeListableURI{URI: storage.NewFileURI("/tmp/document")}, nil)

			return dialog
		},
	})

	// When
	got.Open(browse.OpenOptions{
		Title:  "Open Folder",
		Folder: true,
	}, func(path string) {
		selected <- path
	})
	fyneio.DoAndWait(func() {})

	// Then
	assert.Equal(t, "/tmp/document", receivePath(t, selected))
	assert.Equal(t, "Open Folder", dialog.title)
	assert.Nil(t, dialog.filter)
	assert.True(t, dialog.shown)
}

func TestPicker_Save_ReportsSavedFile(t *testing.T) {
	// Given
	app := test.NewTempApp(t)
	window := app.NewWindow("test")
	closed := false
	selected := make(chan string, 1)
	writer := &fakeWriteCloser{uri: storage.NewFileURI("/tmp/document.fyne")}
	dialog := &fakeFileDialog{}

	got := fynepicker.NewWithDialogs(window, fynepicker.Dialogs{
		NewFileSave: func(callback func(fyneio.URIWriteCloser, error), parent fyneio.Window) fynepicker.FileDialog {
			assert.Same(t, window, parent)
			callback(writer, nil)

			return dialog
		},
	})

	// When
	got.Save(browse.SaveOptions{
		Title:    "Save Document",
		Filename: "document.fyne",
		Filters:  browse.FileFilters{{Name: "Documents", Patterns: []string{"*.fyne"}}},
		OnClosed: func() { closed = true },
	}, func(path string) {
		selected <- path
	})
	fyneio.DoAndWait(func() {})

	// Then
	assert.Equal(t, "/tmp/document.fyne", receivePath(t, selected))
	assert.True(t, writer.closed)
	assert.Equal(t, "document.fyne", dialog.filename)
	assert.Equal(t, "Save Document", dialog.title)
	assert.NotNil(t, dialog.filter)
	assert.True(t, dialog.shown)
	assert.True(t, closed)
}

func TestPicker_Open_IgnoresNilSelection(t *testing.T) {
	// Given
	app := test.NewTempApp(t)
	window := app.NewWindow("test")
	selected := false
	dialog := &fakeFileDialog{}

	got := fynepicker.NewWithDialogs(window, fynepicker.Dialogs{
		NewFileOpen: func(callback func(fyneio.URIReadCloser, error), _ fyneio.Window) fynepicker.FileDialog {
			callback(nil, nil)

			return dialog
		},
	})

	// When
	got.Open(browse.OpenOptions{}, func(string) {
		selected = true
	})
	fyneio.DoAndWait(func() {})

	// Then
	assert.False(t, selected)
	assert.True(t, dialog.shown)
}

func TestPicker_Open_IgnoresFileOpenError(t *testing.T) {
	// Given
	app := test.NewTempApp(t)
	window := app.NewWindow("test")
	selected := make(chan string, 1)
	dialog := &fakeFileDialog{}

	got := fynepicker.NewWithDialogs(window, fynepicker.Dialogs{
		NewFileOpen: func(callback func(fyneio.URIReadCloser, error), _ fyneio.Window) fynepicker.FileDialog {
			callback(nil, errors.New("open failed"))

			return dialog
		},
	})

	// When
	got.Open(browse.OpenOptions{}, func(path string) {
		selected <- path
	})
	fyneio.DoAndWait(func() {})

	// Then
	assertNoPath(t, selected)
	assert.True(t, dialog.shown)
}

func TestPicker_Open_IgnoresFolderOpenError(t *testing.T) {
	// Given
	app := test.NewTempApp(t)
	window := app.NewWindow("test")
	selected := make(chan string, 1)
	dialog := &fakeFileDialog{}

	got := fynepicker.NewWithDialogs(window, fynepicker.Dialogs{
		NewFolderOpen: func(callback func(fyneio.ListableURI, error), _ fyneio.Window) fynepicker.FileDialog {
			callback(nil, errors.New("open folder failed"))

			return dialog
		},
	})

	// When
	got.Open(browse.OpenOptions{Folder: true}, func(path string) {
		selected <- path
	})
	fyneio.DoAndWait(func() {})

	// Then
	assertNoPath(t, selected)
	assert.True(t, dialog.shown)
}

func TestPicker_Save_IgnoresNilSelection(t *testing.T) {
	// Given
	app := test.NewTempApp(t)
	window := app.NewWindow("test")
	selected := false
	dialog := &fakeFileDialog{}

	got := fynepicker.NewWithDialogs(window, fynepicker.Dialogs{
		NewFileSave: func(callback func(fyneio.URIWriteCloser, error), _ fyneio.Window) fynepicker.FileDialog {
			callback(nil, nil)

			return dialog
		},
	})

	// When
	got.Save(browse.SaveOptions{}, func(string) {
		selected = true
	})
	fyneio.DoAndWait(func() {})

	// Then
	assert.False(t, selected)
	assert.True(t, dialog.shown)
}

func TestPicker_Save_IgnoresSaveError(t *testing.T) {
	// Given
	app := test.NewTempApp(t)
	window := app.NewWindow("test")
	selected := make(chan string, 1)
	dialog := &fakeFileDialog{}

	got := fynepicker.NewWithDialogs(window, fynepicker.Dialogs{
		NewFileSave: func(callback func(fyneio.URIWriteCloser, error), _ fyneio.Window) fynepicker.FileDialog {
			callback(nil, errors.New("save failed"))

			return dialog
		},
	})

	// When
	got.Save(browse.SaveOptions{}, func(path string) {
		selected <- path
	})
	fyneio.DoAndWait(func() {})

	// Then
	assertNoPath(t, selected)
	assert.True(t, dialog.shown)
}

type fakeFileDialog struct {
	filename string
	filter   any
	onClosed func()
	shown    bool
	title    string
}

func (f *fakeFileDialog) SetFileName(filename string) {
	f.filename = filename
}

func (f *fakeFileDialog) SetFilter(filter storage.FileFilter) {
	f.filter = filter
}

func (f *fakeFileDialog) SetOnClosed(onClosed func()) {
	f.onClosed = onClosed
}

func (f *fakeFileDialog) SetTitleText(title string) {
	f.title = title
}

func (f *fakeFileDialog) Show() {
	f.shown = true
	if f.onClosed != nil {
		f.onClosed()
	}
}

type fakeListableURI struct {
	fyneio.URI
}

func (u fakeListableURI) List() ([]fyneio.URI, error) {
	return nil, nil
}

type fakeReadCloser struct {
	uri    fyneio.URI
	closed bool
}

func (r *fakeReadCloser) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (r *fakeReadCloser) Close() error {
	r.closed = true

	return nil
}

func (r *fakeReadCloser) URI() fyneio.URI {
	return r.uri
}

type fakeWriteCloser struct {
	uri    fyneio.URI
	closed bool
}

func (w *fakeWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

func (w *fakeWriteCloser) Close() error {
	w.closed = true

	return nil
}

func (w *fakeWriteCloser) URI() fyneio.URI {
	return w.uri
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

func assertNoPath(t *testing.T, paths <-chan string) {
	t.Helper()

	select {
	case path := <-paths:
		require.Failf(t, "unexpected selected path", "got %q", path)
	case <-time.After(50 * time.Millisecond):
	}
}
