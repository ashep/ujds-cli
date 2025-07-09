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
ujds help export
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
