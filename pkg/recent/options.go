package recent

// Option configures a [Recent].
type Option func(*options)

type options struct {
	limit       int
	keepMissing bool
	key         string
}

func defaultOptions() options {
	return options{
		limit: DefaultLimit,
		key:   DefaultPreferencesKey,
	}
}

// WithLimit sets how many items are retained.
func WithLimit(limit int) Option {
	return func(options *options) {
		options.limit = limit
	}
}

// WithKeepMissing keeps missing item paths in the recents.
func WithKeepMissing(keep bool) Option {
	return func(options *options) {
		options.keepMissing = keep
	}
}

// WithPreferencesKey sets the Fyne preferences key used to persist recent paths.
func WithPreferencesKey(key string) Option {
	return func(options *options) {
		options.key = key
	}
}
