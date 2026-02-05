FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o lightarr ./cmd/lightarr

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/lightarr /bin/lightarr
CMD ["/bin/lightarr"]
