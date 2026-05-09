# go-fyne

go-fyne is a small collection of reusable helpers for [Fyne](https://fyne.io/)
applications.

## Installation

Install the module with `go get`:

```sh
go get github.com/adambrett/go-fyne
```

## What's Available

### Browse: Path-Based File Picking

Browse gives Fyne applications a path-based file picker interface.

The shared API lives in [`pkg/browse`](./pkg/browse):

- `Picker` opens and saves paths.
- `OpenOptions` configures open-file and open-folder dialogs.
- `SaveOptions` configures save-file dialogs.
- `FileFilter` and `FileFilters` provide path-based filtering.

Concrete pickers live in backend packages:

- [`pkg/browse/fyne`](./pkg/browse/fyne) uses Fyne's built-in dialogs.
- [`pkg/browse/native`](./pkg/browse/native) uses native dialogs.

Both backend packages expose `New(window)`, so an app can usually switch picker
backends by changing only the import.

Read the [Browse README](./pkg/browse) for usage examples.

### Launcher: Welcome Screen and Recents

[`pkg/launcher`](./pkg/launcher) provides a reusable welcome screen for Fyne
apps that create, open, and remember path-backed items such as documents,
workspaces, images, saves, databases, or any other app-specific file or folder.

It includes:

- create and open actions wired back to your app
- a persisted recent-item list
- stale-entry removal
- configurable labels, icons, logo, window size, split layout, and colours
- pluggable file picking through `browse.Picker`

Read the [Launcher README](./pkg/launcher) for usage examples.

### Recent Items: Persistence and Cleanup

[`pkg/recent`](./pkg/recent) contains the recent-item
storage and cleanup rules used by the launcher. It can also be used on its own
when an app wants recents without rendering the launcher UI.

It handles absolute-path normalization, duplicate removal, missing-path pruning,
ordering, limits, and persistence through Fyne preferences.

Read the [Recent README](./pkg/recent) for usage examples.

## Examples

Runnable examples live under [`examples`](./examples). They are intended to show
how the packages are wired in a small Fyne app.

Run them with:

```sh
make run-browse-fyne
make run-browse-native
make run-launcher
```

## Development

Run the full test suite, including example package compilation, with:

```sh
make test
```

Run the library coverage gate with:

```sh
make coverage
```

The coverage gate measures `./pkg/...`. Example packages are compiled by the
test suite but excluded from the coverage metric.

Format the project with:

```sh
make fmt
```

Run the linter with:

```sh
make lint
```

## Requirements

* Go 1.25 or newer.

## Contributing

### Pull Requests

1. Fork the go-fyne repository.
2. Create a new branch for each feature or improvement.
3. Send a pull request from each feature branch.

### Style Guide

This package is formatted with `gofmt` using `make fmt` and follows idiomatic Go conventions.

If you notice style or API consistency oversights, please send a patch via pull
request.

### Tests

The library is developed using test driven development. All pull requests should be accompanied by passing unit tests with at least 80% coverage. [testify](https://github.com/stretchr/testify) is used for assertions and [mockery](https://github.com/vektra/mockery) is used for mocks.
