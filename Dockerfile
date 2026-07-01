# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN go mod init app-dashboard && go get gopkg.in/yaml.v3
COPY main.go .
COPY index.html . 
RUN CGO_ENABLED=0 GOOS=linux go build -o dashboard main.go

# Stage 2: Minimal runtime image
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/dashboard .
RUN mkdir /config

EXPOSE 8080
CMD ["./dashboard"]
