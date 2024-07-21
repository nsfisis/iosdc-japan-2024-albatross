FROM golang:1.22.3 AS backend-builder

WORKDIR /build
COPY ./backend /build
RUN go build -o /build/server .

################################################################################
FROM node:18.20.2 AS frontend-builder

WORKDIR /build
COPY ./frontend /build
RUN npm install
RUN npm run build

################################################################################
FROM golang:1.22.3

WORKDIR /app
COPY --from=backend-builder /build/server /app/server
COPY --from=frontend-builder /build/dist /app/public

CMD ["/app/server"]
