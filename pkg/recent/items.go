package recent

// Items is a recent-item collection with domain-aware cleanup helpers.
type Items []Item

// NewItemsFromPaths builds item values from stored paths.
func NewItemsFromPaths(paths []string) Items {
	items := make(Items, 0, len(paths))

	for _, path := range paths {
		items = append(items, Item{Path: path})
	}

	return items
}

// Snapshot returns a detached item collection.
func (items Items) Snapshot() Items {
	if len(items) == 0 {
		return nil
	}

	return append(Items(nil), items...)
}

// Clean normalizes, de-duplicates, prunes missing entries, and enforces a limit.
func (items Items) Clean(limit int, keepMissing bool) (Items, bool) {
	cleaned := make(Items, 0, len(items))
	changed := false

	for _, item := range items {
		normalized, ok := item.normalize(true, keepMissing)
		if !ok || cleaned.Contains(normalized) {
			changed = true

			continue
		}

		if len(cleaned) >= limit {
			changed = true

			continue
		}

		if normalized != item {
			changed = true
		}

		cleaned = append(cleaned, normalized)
	}

	return cleaned, changed
}

// WithoutLocation returns items that do not match the target location.
func (items Items) WithoutLocation(target Item) Items {
	if len(items) == 0 {
		return nil
	}

	filtered := make(Items, 0, len(items))
	for _, item := range items {
		if item.SameLocation(target) {
			continue
		}

		filtered = append(filtered, item)
	}

	return filtered
}

// WithFirst returns a collection with item first and any duplicate location removed.
func (items Items) WithFirst(item Item, limit int) Items {
	if limit <= 0 {
		return nil
	}

	ordered := make(Items, 0, min(limit, len(items)+1))
	ordered = append(ordered, item)

	for _, existing := range items.WithoutLocation(item) {
		if len(ordered) >= limit {
			break
		}

		ordered = append(ordered, existing)
	}

	return ordered
}

// Paths returns the stored paths for the items.
func (items Items) Paths() []string {
	if len(items) == 0 {
		return nil
	}

	paths := make([]string, len(items))
	for i, item := range items {
		paths[i] = item.Path
	}

	return paths
}

// Contains reports whether the target item location exists in the collection.
func (items Items) Contains(target Item) bool {
	for _, item := range items {
		if item.SameLocation(target) {
			return true
		}
	}

	return false
}

// Same reports whether another collection has the same ordered item values.
func (items Items) Same(other Items) bool {
	if len(items) != len(other) {
		return false
	}

	for i := range items {
		if items[i] != other[i] {
			return false
		}
	}

	return true
}
