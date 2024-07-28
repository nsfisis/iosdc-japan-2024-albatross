FROM golang:1.22.3 AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . /build
RUN go build -o /build/server .

################################################################################
FROM golang:1.22.3

WORKDIR /app
COPY --from=builder /build/server /app/server

CMD ["/app/server"]
