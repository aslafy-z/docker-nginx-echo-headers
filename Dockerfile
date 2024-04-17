FROM golang:1.22 as builder
COPY go.mod go.sum main.go /go/src/echo-server/
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /usr/local/bin/echo-server /go/src/echo-server/main.go

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder --chmod=0755 /usr/local/bin/echo-server /usr/local/bin/echo-server
ENTRYPOINT ["/usr/local/bin/echo-server"]
