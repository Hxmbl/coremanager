# CoreManager

A simple CLI tool to manage CPU cores on Linux.

## Installation

### From source

```bash
go install github.com/hxmbl/coremanager@latest
```

### From release

Download the latest binary from [Releases](https://github.com/Hxmbl/coremanager/releases).

## Usage

```bash
coremanager --help
coremanager core-count
coremanager cpu-model
coremanager disable-cores 3
coremanager enable-cores all
coremanager debug-info
```

## Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `disable-cores [N\|all]` | `dc` | Disable N CPU cores or all secondary cores |
| `enable-cores [N\|all]` | `ec` | Enable N CPU cores or all secondary cores |
| `core-count` | `cc` | Display total and active core counts |
| `cpu-model` | `cm` | Display CPU model name |
| `debug-info` | `debug` | Show detailed CPU info |

All commands accept `-v`/`--verbose` for detailed output.

## License

No.
