# Recent

Recent stores and cleans path-backed items for Fyne applications.

It is useful when an app needs a persistent "recent files", "recent
workspaces", "recent images", or similar list without tying that behavior to a
specific screen.

# Installation

Recent is part of `github.com/adambrett/go-fyne`:

```sh
go get github.com/adambrett/go-fyne
```

# Basic Usage

Create a recent-item store from Fyne preferences:

```go
import "github.com/adambrett/go-fyne/pkg/recent"

recents := recent.New(
	app.Preferences(),
	recent.WithLimit(10),
)
```

Add an item after it has been created or opened successfully:

```go
recents.Add(recent.Item{
	Name: "Example",
	Path: "/tmp/example.fyne",
})
```

Read the current list:

```go
for _, item := range recents.Items() {
	fmt.Println(item.DisplayName(), item.Path)
}
```

# Items

An item has a display name and a filesystem path:

```go
item := recent.Item{
	Name: "Example",
	Path: "/tmp/example.fyne",
}
```

If `Name` is empty, `DisplayName` uses the base name of `Path`.

# Cleanup Rules

Recent normalizes paths to absolute paths, removes duplicate locations, prunes
missing paths by default, preserves ordering, and enforces a limit.

Keep missing paths when items may live on removable drives, network shares, or
other temporarily unavailable locations:

```go
recents := recent.New(
	app.Preferences(),
	recent.WithKeepMissing(true),
)
```

# Replacing the List

Use `Replace` when importing or restoring a complete recent list:

```go
recents.Replace(recent.Items{
	{Path: "/tmp/first.fyne"},
	{Path: "/tmp/second.fyne"},
})
```

The replacement list is cleaned using the same normalization, de-duplication,
missing-path, and limit rules.
