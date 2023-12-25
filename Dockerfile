FROM golang:alpine AS builder

WORKDIR /build

COPY go.* ./
RUN go mod download

ENV PROJECT_NAME $(grep -m1 "^module" go.mod | awk '{print $2}')

COPY . .

RUN go install github.com/go-task/task/v3/cmd/task@latest

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN task build

FROM alpine:latest

COPY --from=builder /build/${PROJECT_NAME} /server
COPY --from=builder /build/public /public

CMD ["/server"]
