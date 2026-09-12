# UJDS Command Line Interface

## Installation

Using pre-compiled binaries: https://github.com/ashep/ujds-cli/releases.

Using Go:

```shell
go install -v github.com/ashep/ujds-cli@latest
```

## Configuration

Copy `config.sample.yaml` to `config.yaml` into the directory you are running the command from and set `host`
and `token` values.

## Getting help

You can get a list of available command flags with their descriptions using the `help` command, for example:

```shell
ujds-cli help export
```

## `index` command

Index operations. Available commands:

- `list`

## `index list` command

List indices. Available flags:

- `-n`, `--names`: index name patterns to list. Allowed wildcard symbols are: `*`. Default: `*`.
  Example: `-n books.*.history,books.2024.fiction`.
- `-f`, `--format`: output format. Allowed variables: `{name}`, `{title}`. Default: `{name}`.
  Example: `-f '{name}: {title}'`.

### Examples

To get all index names:

```shell
ujds-cli index list
```

To get index names and titles, that have `book` word in the name:

```shell
ujds-cli index list -f '{name}: {title}' -n *book*
```

## `export` command

Export records. Available flags:

- `-i`, `--index`: index name patterns to scan. Allowed wildcard symbols are: `*`. Default: `*`.
  Example: `-n books.*.history,books.2024.fiction`.
- `-o`, `--out`: output file name, including the extension which is used to determine the output format. Currently, only
  CSV is supported. Default: `out.csv`.
- `--overwrite`: overwrite existing output file. Default: `false`.

### Examples

To export records from indices having prefix `books.2023.`, and the index named `books.2024.fivestars` to
the `books.csv` file, overwriting an existing one:

```shell
ujds-cli export --overwrite -i books.2023.*,books.2024.fivestars -o books.csv
```

## `record` command

Record operations. Available commands:

- `get`
- `find`
- `history`

## `record get` command

Get a single record by its index and ID and print it to stdout:

```shell
ujds-cli record get -i <index_id> <record_id>
```

Available flags:

- `-i`, `--index`: index name. Required.
- `-f`, `--format`: output format. It is a [Go template](https://pkg.go.dev/text/template) executed against the record.
  Available fields: `{{.Id}}`, `{{.Rev}}`, `{{.Index}}`, `{{.CreatedAt}}`, `{{.UpdatedAt}}`, `{{.TouchedAt}}`
  (UTC date and time), `{{.CreatedAtUnix}}`, `{{.UpdatedAtUnix}}`, `{{.TouchedAtUnix}}` (the same timestamps as raw
  unix seconds), `{{.Data}}` (raw JSON), `{{.DataTable}}` (data rendered as a table).
  Default: `ID: {{.Id}}\nCreated: {{.CreatedAt}}\nUpdated: {{.UpdatedAt}}\nTouched: {{.TouchedAt}}\nData:\n{{.DataTable}}`.

By default the record data is printed as a two-column table, where keys are flattened dotted paths into the JSON
document and array items are addressed by their index:

```
ID: 978-0451524935
Created: 2024-03-31 00:00:00 UTC
Updated: 2024-03-31 00:00:00 UTC
Touched: 2024-03-31 00:00:00 UTC
Data:
    author.born  1903
    author.name  George Orwell
    genres.0     dystopia
    genres.1     fiction
    pages        328
    title        1984
```

Keys are sorted segment by segment, with numeric segments compared as numbers, so `genres.2` comes before `genres.10`.
JSON `null` values are not listed. If the data is not valid JSON, it is printed as is.

Timestamps are always
printed in UTC, regardless of the local timezone.

### Examples

To get the record `978-0451524935` from the `books.2024.fiction` index:

```shell
ujds-cli record get -i books.2024.fiction 978-0451524935
```

To print only the record data as raw JSON:

```shell
ujds-cli record get -i books.2024.fiction 978-0451524935 -f '{{.Data}}'
```

## Debugging

User `APP_DEBUG` environment variable to get verbose logging:

```shell
APP_DEBUG=1 ujds-cli ...
```

## Changelog

### 0.3 (2024-03-31)

Record metadata added to the `export` command output.

### 0.2 (2024-03-31)

- The `index list` command added.
- The `api_key` configuration parameter renamed to `token`.

### 0.1 (2024-03-16)

Initial release.

## Authors

- [Oleksandr Shepetko](https://shepetko.com)
