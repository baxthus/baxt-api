FROM golang:alpine AS builder

WORKDIR /build

COPY go.* ./
RUN go mod download

ENV PROJECT_NAME $(grep -m1 "^module" go.mod | awk '{print $2}')

COPY . .

RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN ./bin/task build

FROM alpine:latest

COPY --from=builder /build/${PROJECT_NAME} /server
COPY --from=builder /build/public /public

CMD ["/server"]
