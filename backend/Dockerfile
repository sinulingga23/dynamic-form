FROM golang:1.19-alpine AS builder
RUN apk add --no-cache ca-certificates git

WORKDIR /opt/apps

ENV GOPROXY="https://proxy.golang.org,direct"
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /backend-service .

FROM alpine AS release
RUN apk add --no-cache ca-certificates

WORKDIR /backend-service
COPY --from=builder /backend-service ./main
ENTRYPOINT ["/backend-service/main"]