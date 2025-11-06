FROM golang:alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -trimpath -o marketflow cmd/main.go


FROM scratch

WORKDIR /app

COPY --from=builder /app/marketflow .

ENTRYPOINT [ "./marketflow" ]