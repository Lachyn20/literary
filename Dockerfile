# Build stage
FROM golang:1.22 AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

# Final stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/server ./server
# create uploads dir
RUN mkdir -p /app/uploads
VOLUME ["/app/uploads"]
EXPOSE 8080
ENTRYPOINT ["/app/server"]
