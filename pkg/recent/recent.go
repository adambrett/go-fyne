package recent

import "fyne.io/fyne/v2"

const (
	DefaultLimit          = 5
	DefaultPreferencesKey = "recent-items"
)

// Recent owns recent item ordering, de-duplication, and pruning rules.
type Recent struct {
	limit       int
	keepMissing bool
	items       Items
	preferences fyne.Preferences
	key         string
}

// New builds recent item state from options.
func New(preferences fyne.Preferences, opts ...Option) *Recent {
	options := defaultOptions()

	for _, opt := range opts {
		opt(&options)
	}

	if options.limit <= 0 {
		options.limit = DefaultLimit
	}

	if options.key == "" {
		options.key = DefaultPreferencesKey
	}

	recent := &Recent{
		limit:       options.limit,
		keepMissing: options.keepMissing,
		preferences: preferences,
		key:         options.key,
	}

	items, changed := NewItemsFromPaths(recent.preferences.StringList(recent.key)).Clean(recent.limit, recent.keepMissing)
	recent.items = items

	if changed {
		recent.save()
	}

	return recent
}

// Items returns the current recent items.
func (r *Recent) Items() Items {
	if r == nil || len(r.items) == 0 {
		return nil
	}

	return r.items.Snapshot()
}

// Replace replaces the recents with cleaned item entries.
func (r *Recent) Replace(items Items) bool {
	if r == nil {
		return false
	}

	cleaned, _ := items.Clean(r.limit, r.keepMissing)

	return r.replace(cleaned)
}

// Add puts an item at the top of the recents.
func (r *Recent) Add(item Item) bool {
	if r == nil {
		return false
	}

	item, ok := item.normalize(true, r.keepMissing)
	if !ok {
		return false
	}

	return r.replace(r.items.WithFirst(item, r.limit))
}

// Remove deletes an item from the recents.
func (r *Recent) Remove(item Item) bool {
	if r == nil {
		return false
	}

	item, ok := item.normalize(false, r.keepMissing)
	if !ok {
		return false
	}

	return r.replace(r.items.WithoutLocation(item))
}

func (r *Recent) replace(items Items) bool {
	if r.items.Same(items) {
		return false
	}

	r.items = append(Items(nil), items...)
	r.save()

	return true
}

func (r *Recent) save() {
	r.preferences.SetStringList(r.key, r.items.Paths())
}
