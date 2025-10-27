FROM golang:1.25.0-alpine3.22 AS builder

WORKDIR /app

COPY . .

RUN go build -o marketflow test.go


FROM alpine

WORKDIR /app

COPY --from=builder /app/marketflow .

ENTRYPOINT [ "./marketflow" ]