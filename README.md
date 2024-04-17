# docker-echo-headers

A simple Docker image that echoes request headers.

Find Docker image at <https://github.com/aslafy-z/docker-echo-headers/pkgs/container/echo-headers>.

## Usage

```bash
docker run --rm -p 8080:8080 ghcr.io/aslafy-z/echo-headers:latest
```

## Environment variables

| Variable | Default | Description |
| --- | --- | --- |
| `ECHO_ADDR` | `:8080` | addr to bind to |
| `ECHO_CONTEXT` | `true` | show extra context request |
| `ECHO_DELAY` | `0s` | add extra delay to responses |
| `ECHO_RAND_BYTES` | `0` | add extra N random bytes to responses |

## Query string

| Variable | Default | Description |
| --- | --- | --- |
| `delay` | `0s` | add extra delay to responses or override default one |
