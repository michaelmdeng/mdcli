# mdcli

CLI for personal tools.

## Build

Build binary at `./bin/mdcli`

```bash
make build
```

Install/link binary to `~/.local/bin`

```bash
make install
```

## Test

```bash
make test
```

## Lint/fmt

Requires `golangci-lint` install

```bash
make fmt
```

# Usage

```bash
mdcli -h
```

```dockerfile
COPY --from ghcr.io/michaelmdeng/mdcli/mdcli:latest /bin/mdcli .
```
