# mdcli

CLI for personal tools.

## Build

```bash
make build
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
