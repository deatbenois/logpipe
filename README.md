# logpipe

Lightweight CLI for tailing and filtering structured JSON logs from multiple sources.

---

## Installation

```bash
go install github.com/yourusername/logpipe@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logpipe.git
cd logpipe
go build -o logpipe .
```

---

## Usage

Tail logs from a file and filter by log level:

```bash
logpipe tail --file /var/log/app.json --level error
```

Pipe logs from multiple sources and filter by a field value:

```bash
logpipe tail --file /var/log/app.json --file /var/log/worker.json --filter service=auth
```

Pretty-print raw JSON log output:

```bash
cat app.json | logpipe fmt
```

### Flags

| Flag | Description |
|------|-------------|
| `--file` | Path to a JSON log file (repeatable) |
| `--level` | Filter by log level (e.g. `info`, `warn`, `error`) |
| `--filter` | Filter by field value in `key=value` format |
| `--follow` | Continuously tail the file for new entries |
| `--output` | Output format: `pretty` or `json` (default: `pretty`) |

---

## Requirements

- Go 1.21 or later

---

## License

MIT © 2024 yourusername