# grokir

`grokir` is a small command-line client for [Grokipedia](https://grokipedia.com).  
It lets you search articles and fetch a page by slug from the terminal.

## Requirements

- Go `1.25` (as defined in `go.mod`)
- `make` (optional but recommended)

## Build

Build the binary into `dist/`:

```bash
make build
```

This produces:

- `dist/grokir`

## Install

Install to `~/.local/bin`:

```bash
make install
```

Make sure `~/.local/bin` is in your `PATH`.

## Clean

Remove build artifacts:

```bash
make clean
```

## Run

From the project root:

```bash
./dist/grokir <command> [args]
```

Or, after install:

```bash
grokir <command> [args]
```

## Commands

### `search`

Search Grokipedia by query.

```bash
grokir search "kubernetes scheduler"
```

Options:

- `-l <num>`: max number of results (default: `10`)
- `-o <num>`: offset for pagination (default: `0`)

Example:

```bash
grokir search -l 5 -o 10 "distributed systems"
```

### `page`

Fetch a page by slug.

```bash
grokir page kubernetes-scheduler
```

### `version`

Show version and build date.

```bash
grokir version
```

## Output Modes

By default, output is human-readable text.

Use JSON output with `--json`:

```bash
grokir --json search "kubernetes"
grokir --json page kubernetes-scheduler
```

> Note: `--json` is a global flag, so place it **before** the command name.

## Development

Run checks:

```bash
make test
make lint
```
