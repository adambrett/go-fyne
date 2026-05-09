package browse

// OpenOptions configures an open-file or open-folder dialog.
type OpenOptions struct {
	// Title sets the dialog title.
	Title string

	// Filters limits selectable files by name pattern.
	Filters FileFilters

	// Folder opens a directory chooser instead of a file chooser.
	Folder bool

	// OnClosed runs after the picker closes.
	OnClosed func()
}

// SaveOptions configures a save-file dialog.
type SaveOptions struct {
	// Title sets the dialog title.
	Title string

	// Filename sets the initial filename.
	Filename string

	// Filters limits saved files by name pattern.
	Filters FileFilters

	// ConfirmOverwrite asks the user before replacing an existing file.
	ConfirmOverwrite bool

	// OnClosed runs after the picker closes.
	OnClosed func()
}
