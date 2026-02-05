FROM golang:1.23-alpine AS builder

COPY lightarr /bin/lightarr

RUN /bin/lightarr
