package browse

// Picker describes a file browser that reports selected filesystem paths.
type Picker interface {
	// Open shows an open-file or open-folder dialog.
	Open(options OpenOptions, onSelected func(path string))

	// Save shows a save-file dialog.
	Save(options SaveOptions, onSelected func(path string))
}
