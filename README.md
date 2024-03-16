# UJDS Command Line Interface

## Configuration

Copy `config.sample.yaml` to `config.yaml` and set `host` and `api_key` values.

## Export

Get command help:

```shell
ujds-cli help export
```

### Examples

Export records from indices having prefix `books.2023.*`, and index named `books.2024.fivestars` to CSV file,
overwriting an existing one:

```shell
ujds-cli export --overwrite -i books.2023.*,books.2024.fivestars -o books.csv
```

## Debugging

User `APP_DEBUG` environment variable to get verbose logging:

```shell
APP_DEBUG=1 ujds-cli ...
```

## Changelog

### 0.1 (2024-03-16)

Initial release.

## Authors

- [Oleksandr Shepetko](https://shepetko.com)
