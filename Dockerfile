FROM golang:alpine AS builder

WORKDIR /build

COPY go.* ./
RUN go mod download


COPY . .

RUN go install github.com/go-task/task/v3/cmd/task@latest

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN task build

FROM alpine:latest

COPY --from=builder /build/server /server
COPY --from=builder /build/public /public

CMD ["/server"]
