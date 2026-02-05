FROM node:25-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.24-alpine AS go-builder
WORKDIR /app
COPY . .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
RUN go build -o lightarr

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=go-builder /app/lightarr /bin/lightarr
CMD ["/bin/lightarr"]
