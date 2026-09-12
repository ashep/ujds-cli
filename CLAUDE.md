# ujds-cli

Command line client for [UJDS](https://github.com/ashep/ujds) (Universal JSON Data Store). Go + Cobra, talks to the
server over Connect RPC using the `github.com/ashep/ujds/sdk` client.

## Build & check

```shell
go build ./...
go vet ./...
go run . <command> --help    # run locally
```

Dependencies are vendored (`/vendor`, gitignored) — run `go mod vendor` after changing `go.mod`.
`task init` clones the shared CI config into `.ci`; CI itself lives in `.github/workflows`.

## Architecture

```
main.go                  -> apprun.Run(app.New, app.Config{})
internal/app/            -> App wrapper; builds the ujds SDK client from Config{Host, Token}
internal/command/        -> Cobra wiring only: one file per top-level command group
internal/<group>/<cmd>/  -> one package per subcommand, holding the actual logic
pkg/jsontree/            -> helper for flattening JSON into column paths (used by export)
```

The layering is strict and worth preserving: `internal/command/*.go` contains **no business logic**. Each subcommand
constructor only declares flags/args and calls into its own package, e.g. `find.New(cli).Find(ctx, ..., out)`.

Current groups: `index` (`list`), `record` (`get`, `find`, `history`), plus the top-level `export`.

## Conventions

- Logic packages are named after the subcommand (`internal/record/get`) and export a type with the same name as the
  package's verb plus `New(cli)` constructor and a single method, e.g. `Get.Get(...)`.
- SDK access: `cli.I` is the index service, `cli.R` is the record service. Calls are wrapped in
  `connect.NewRequest(&<proto>.<X>Request{...})`; the proto packages are imported as `recordproto` / `indexproto`.
- Error wrapping is always `fmt.Errorf("<context>: %w", err)`; the server call site uses `"ujds response: %w"`.
- Output goes to `cmd.OutOrStdout()` (passed in as an `io.Writer`), never to `fmt.Println` — this keeps commands
  testable. The `index list` command is the odd one out; it takes a `printer` func.
- Most record commands accept a `-f/--format` flag that is a `text/template` executed against the result. The helper
  appends a trailing `\n` if the format lacks one. `index list` instead uses `{name}`-style placeholder substitution.
- `record get` and `history` map the proto record into a local struct before executing the template, so template fields
  are the struct's, not the proto's. `record get` exposes both `{{.Data}}` (raw JSON) and `{{.DataTable}}`, the latter
  rendered by `internal/record/get/table.go` via `jsontree` + `text/tabwriter`. Note `jsontree.Keys()` skips `nil`
  values, so JSON nulls never appear in the table (or in export's CSV columns).
- Timestamps come off the wire as unix seconds. `record get` renders them as `time.DateTime` in **UTC** (with the raw
  values kept on `{{.CreatedAtUnix}}` and friends); `history` still formats its `TimeStr` in **local** time. Worth
  unifying if it ever comes up.
- Record identity is `(index name, record id)`. The index always comes from the `-i/--index` flag (required);
  `record get` takes the record id as a positional arg, while `history` takes it as an `--id` flag.
- `zerolog.Logger` is threaded only into commands that need it (`export`, `index list`); record commands don't take one.

## Record proto fields

`recordproto.Record`: `Id`, `Rev`, `Index`, `CreatedAt`, `UpdatedAt`, `TouchedAt` (unix seconds), `Data` (JSON string).
Note the template field names are Go names (`{{.Id}}`, not `{{.ID}}`) except in `history`, which maps the proto into a
local `historyRecord{ID, Time, TimeStr, Data}` struct first.

## Docs

`README.md` documents every command group and flag and holds the changelog — update it when adding or changing a
command.
