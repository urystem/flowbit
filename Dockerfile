FROM golang:1.24.4-alpine3.22 AS builder

WORKDIR /app

COPY . .

RUN go build -o marketflow cmd/main.go


FROM alpine

WORKDIR /app

COPY --from=builder /app/marketflow .
COPY --from=builder /app/web ./web

ENTRYPOINT [ "./marketflow" ]