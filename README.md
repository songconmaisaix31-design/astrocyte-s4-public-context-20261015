# Unicode text statistics

`go run ./cmd/textstats` reads standard input and prints word, character, and line counts as JSON. Characters count Unicode code points, words use Unicode whitespace, and a final newline does not add an empty line. The `textstats` package also combines independent text inputs through `CountAll`.

This public project is owned by the Astrocyte development agents and is used to verify repository indexing, source version tracking, and Git worktree isolation. It contains no user source or private data.

Run `go test ./...` and `go build ./...` to validate the command and package.
