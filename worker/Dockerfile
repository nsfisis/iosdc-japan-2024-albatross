FROM golang:1.22.3 AS builder

WORKDIR /build

RUN apt-get update && apt-get install -y curl xz-utils
RUN curl https://wasmtime.dev/install.sh -sSf | bash -s -- --version v23.0.1

COPY go.mod go.sum ./
RUN go mod download
COPY *.go /build
RUN CGO_ENABLED=0 go build -o /build/server .

# ################################################################################
FROM ghcr.io/swiftwasm/swift:5.10-focal

WORKDIR /app

COPY --from=builder /root/.wasmtime/bin/wasmtime /usr/bin/wasmtime
COPY --from=builder /build/server /app/server

CMD ["/app/server"]
