# Stage 1: Build the Go application
FROM golang:alpine as builder
RUN mkdir /build
WORKDIR /build
ADD . .
RUN go mod download
RUN go build -o main .

# Stage 2: Create a minimal image for deployment
FROM alpine:latest
RUN adduser -S -D -H -h /app appuser
USER appuser
WORKDIR /app
COPY --from=builder /build/main .
EXPOSE 3000
CMD ["./main"]
