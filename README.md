# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected listeners with configurable rules.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with a config file:

```bash
portwatch --config /etc/portwatch/config.yaml
```

Example `config.yaml`:

```yaml
interval: 30s
allowed_ports:
  - 22
  - 80
  - 443
alerts:
  - type: log
    path: /var/log/portwatch.log
  - type: webhook
    url: https://hooks.example.com/alert
```

Run a one-time scan without the daemon:

```bash
portwatch scan
```

View currently monitored ports:

```bash
portwatch status
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `./config.yaml` | Path to config file |
| `--interval` | `30s` | Polling interval |
| `--verbose` | `false` | Enable verbose logging |

## License

MIT © 2024 [yourusername](https://github.com/yourusername)