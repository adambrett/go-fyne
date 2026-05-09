package recent

import (
	"os"
	"path/filepath"
)

// Item describes an entry shown in the recent items list.
type Item struct {
	Name string
	Path string
}

// DisplayName returns Name when set, otherwise the base name of Path.
func (p Item) DisplayName() string {
	if p.Name != "" {
		return p.Name
	}

	if p.Path == "" {
		return ""
	}

	return filepath.Base(p.Path)
}

// SameLocation reports whether two items refer to the same item path.
func (p Item) SameLocation(other Item) bool {
	if p.Path != "" || other.Path != "" {
		return p.Path == other.Path
	}

	return p.Name == other.Name
}

func (p Item) normalize(checkExists bool, keepMissing bool) (Item, bool) {
	if p.Path == "" {
		return Item{}, false
	}

	abs, err := filepath.Abs(p.Path)
	if err != nil || abs == "" {
		return Item{}, false
	}

	if checkExists && !keepMissing {
		if _, err := os.Stat(abs); err != nil {
			return Item{}, false
		}
	}

	p.Path = abs

	return p, true
}
