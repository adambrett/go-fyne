package recent_test

import (
	"os"
	"path/filepath"
	"testing"

	fyneTest "fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adambrett/go-fyne/pkg/recent"
)

func TestRecent_Add_MovesItemToTopAndCollapsesDuplicates(t *testing.T) {
	// Given
	app := fyneTest.NewApp()
	first := touchItemFile(t, "first.db")
	second := touchItemFile(t, "second.db")
	app.Preferences().SetStringList(recent.DefaultPreferencesKey, []string{first, second})

	recents := recent.New(app.Preferences())

	// When
	require.True(t, recents.Add(recent.Item{Name: "Second", Path: second}))

	// Then
	assert.Equal(t, recent.Items{
		{Name: "Second", Path: second},
		{Path: first},
	}, recents.Items())
}

func TestRecent_Add_EnforcesLimit(t *testing.T) {
	// Given
	app := fyneTest.NewApp()
	recents := recent.New(app.Preferences(), recent.WithLimit(2))
	first := touchItemFile(t, "first.db")
	second := touchItemFile(t, "second.db")
	third := touchItemFile(t, "third.db")

	// When
	require.True(t, recents.Add(recent.Item{Path: first}))
	require.True(t, recents.Add(recent.Item{Path: second}))
	require.True(t, recents.Add(recent.Item{Path: third}))

	// Then
	assert.Equal(t, recent.Items{
		{Path: third},
		{Path: second},
	}, recents.Items())
}

func TestRecent_Remove_IsIdempotent(t *testing.T) {
	// Given
	app := fyneTest.NewApp()
	first := touchItemFile(t, "first.db")
	second := touchItemFile(t, "second.db")
	app.Preferences().SetStringList(recent.DefaultPreferencesKey, []string{first, second})

	recents := recent.New(app.Preferences())

	// When
	require.True(t, recents.Remove(recent.Item{Path: first}))
	assert.False(t, recents.Remove(recent.Item{Path: first}))

	// Then
	assert.Equal(t, recent.Items{{Path: second}}, recents.Items())
}

func TestRecent_Replace(t *testing.T) {
	// Given
	app := fyneTest.NewApp()
	first := touchItemFile(t, "first.db")
	second := touchItemFile(t, "second.db")
	recents := recent.New(app.Preferences(), recent.WithLimit(2))

	// When
	changed := recents.Replace(recent.Items{
		{Path: first},
		{Path: second},
	})
	unchanged := recents.Replace(recent.Items{
		{Path: first},
		{Path: second},
	})

	// Then
	assert.True(t, changed)
	assert.False(t, unchanged)
	assert.Equal(t, recent.Items{{Path: first}, {Path: second}}, recents.Items())
	assert.Equal(t, []string{first, second}, app.Preferences().StringList(recent.DefaultPreferencesKey))
}

func TestRecent_Load_PrunesMissingItemsByDefault(t *testing.T) {
	// Given
	app := fyneTest.NewApp()
	valid := touchItemFile(t, "valid.db")
	missing := filepath.Join(t.TempDir(), "missing.db")
	app.Preferences().SetStringList(recent.DefaultPreferencesKey, []string{valid, missing})

	// When
	recents := recent.New(app.Preferences())

	// Then
	assert.Equal(t, recent.Items{{Path: valid}}, recents.Items())
}

func TestRecent_Load_KeepsMissingItemsWhenConfigured(t *testing.T) {
	// Given
	app := fyneTest.NewApp()
	missing := filepath.Join(t.TempDir(), "missing.db")
	app.Preferences().SetStringList(recent.DefaultPreferencesKey, []string{missing})

	// When
	recents := recent.New(
		app.Preferences(),
		recent.WithKeepMissing(true),
	)

	// Then
	assert.Equal(t, recent.Items{{Path: missing}}, recents.Items())
}

func TestRecent_Load_CleansAndSavesStoredPaths(t *testing.T) {
	// Given
	app := fyneTest.NewApp()
	valid := touchItemFile(t, "valid.db")
	missing := filepath.Join(t.TempDir(), "missing.db")
	app.Preferences().SetStringList(recent.DefaultPreferencesKey, []string{valid, valid, missing})

	// When
	recents := recent.New(app.Preferences())

	// Then
	assert.Equal(t, recent.Items{{Path: valid}}, recents.Items())
	assert.Equal(t, []string{valid}, app.Preferences().StringList(recent.DefaultPreferencesKey))
}

func touchItemFile(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, []byte(""), 0o600))

	return path
}
