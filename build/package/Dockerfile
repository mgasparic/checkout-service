FROM golang:1.19.3-alpine AS builder
WORKDIR /build
COPY . .
WORKDIR /build/cmd/checkout
RUN go mod tidy
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:3.17
WORKDIR /root/
COPY --from=builder /build/cmd/checkout/checkout service
CMD ["./service"]