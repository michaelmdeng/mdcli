# Dockerfile

# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/mdcli .

# Final stage
FROM gcr.io/distroless/static
COPY --from=builder /bin/mdcli /bin/mdcli
ENTRYPOINT ["/bin/mdcli"]
