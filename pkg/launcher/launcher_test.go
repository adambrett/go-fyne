package launcher_test

import (
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	fyneTheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adambrett/go-fyne/pkg/launcher"
	"github.com/adambrett/go-fyne/pkg/launcher/theme"
	"github.com/adambrett/go-fyne/pkg/recent"
)

func TestLauncher_RememberItem_UpdatesPreferences(t *testing.T) {
	// Given
	app := test.NewApp()
	first := touchItemFile(t, "first.db")
	second := touchItemFile(t, "second.db")
	l := launcher.New(app.Preferences(), app.NewWindow("test"), newTestItem, openTestItem, launcher.WithRecentLimit(2))

	// When
	l.RememberItem(recent.Item{Path: first})
	l.RememberItem(recent.Item{Path: second})
	l.RememberItem(recent.Item{Name: "First", Path: first})

	// Then
	assert.Equal(t, fyne.NewSize(720, 480), l.Size())
	assert.Equal(t, []string{first, second}, app.Preferences().StringList(recent.DefaultPreferencesKey))
}

func TestLauncher_RecentItems_LoadsPreferences(t *testing.T) {
	// Given
	app := test.NewApp()
	first := touchItemFile(t, "first.db")
	second := touchItemFile(t, "second.db")
	app.Preferences().SetStringList(recent.DefaultPreferencesKey, []string{first})

	// When
	l := launcher.New(
		app.Preferences(),
		app.NewWindow("test"),
		newTestItem,
		openTestItem,
		launcher.WithWindowSize(fyne.NewSize(800, 600)),
	)

	// Then
	assert.Equal(t, fyne.NewSize(800, 600), l.Size())
	assert.Equal(t, []string{first}, app.Preferences().StringList(recent.DefaultPreferencesKey))

	// When
	l.RememberItem(recent.Item{Path: second})

	// Then
	assert.Equal(t, []string{second, first}, app.Preferences().StringList(recent.DefaultPreferencesKey))
}

func TestLauncher_Container_WiresCallbacksAndOptions(t *testing.T) {
	// Given
	app := test.NewApp()
	first := touchItemFile(t, "first.db")
	second := touchItemFile(t, "second.db")
	third := touchItemFile(t, "third.db")
	app.Preferences().SetStringList(recent.DefaultPreferencesKey, []string{first, second})

	var newItem bool
	var opened recent.Item

	l := launcher.New(
		app.Preferences(),
		app.NewWindow("test"),
		func() (recent.Item, bool) {
			newItem = true

			return recent.Item{Path: third}, true
		},
		func(item recent.Item) {
			opened = item
		},
		launcher.WithTitle("Workspace"),
		launcher.WithCreateLabel("Create"),
		launcher.WithOpenLabel("Browse"),
		launcher.WithRecentTitle("Latest"),
		launcher.WithEmptyRecentText("No items"),
		launcher.WithLogo(testLogoResource()),
		launcher.WithLogoCanvas(canvas.NewCircle(color.Black)),
		launcher.WithLogoSize(fyne.NewSize(44, 44)),
		launcher.WithCreateIcon(fyneTheme.ContentAddIcon()),
		launcher.WithOpenIcon(fyneTheme.FolderOpenIcon()),
		launcher.WithSplitOffset(0.5),
		launcher.WithTheme(theme.Theme{PrimaryText: color.White}),
	)

	// When
	require.NotNil(t, l.CanvasObject())
	root := layOutLauncherForTap(l)

	tapButtonWithText(t, leftActionsVBox(root), "Create")

	recents := rightRecentVBox(root)
	test.Tap(recents.Objects[1].(fyne.Tappable))

	removeCtl := test.WidgetRenderer(recents.Objects[2].(fyne.Widget)).Objects()[3]
	test.Tap(removeCtl.(fyne.Tappable))

	// Then
	assert.True(t, newItem)
	assert.Equal(t, third, opened.Path)
	assert.Equal(t, []string{third, second}, app.Preferences().StringList(recent.DefaultPreferencesKey))
}

func newTestItem() (recent.Item, bool) {
	return recent.Item{}, false
}

func openTestItem(recent.Item) {}

func touchItemFile(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, []byte(""), 0o600))

	return path
}

func testLogoResource() fyne.Resource {
	return fyne.NewStaticResource("test-logo.svg", []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="1" height="1"></svg>`))
}

func layOutLauncherForTap(l *launcher.Launcher) fyne.CanvasObject {
	c := test.NewCanvas()
	c.SetContent(l.CanvasObject())
	c.Resize(fyne.NewSize(1024, 768))

	return l.CanvasObject()
}

func leftActionsVBox(root fyne.CanvasObject) *fyne.Container {
	stack := root.(*fyne.Container)
	split := stack.Objects[1].(*fyne.Container).Objects[0].(*container.Split)

	leftPad := split.Leading.(*fyne.Container)
	center := leftPad.Objects[0].(*fyne.Container)

	return center.Objects[0].(*fyne.Container)
}

func tapButtonWithText(t *testing.T, box *fyne.Container, wantText string) {
	t.Helper()

	var tap func(*fyne.Container)
	tap = func(c *fyne.Container) {
		for _, o := range c.Objects {
			switch x := o.(type) {
			case *widget.Button:
				if x.Text == wantText {
					test.Tap(x)
					return
				}
			case *fyne.Container:
				tap(x)
			}
		}
	}

	tap(box)
}

func rightRecentVBox(root fyne.CanvasObject) *fyne.Container {
	stack := root.(*fyne.Container)
	split := stack.Objects[1].(*fyne.Container).Objects[0].(*container.Split)

	rightPad := split.Trailing.(*fyne.Container)

	return rightPad.Objects[0].(*fyne.Container)
}
