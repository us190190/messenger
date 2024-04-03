# Use an official Go runtime as a base image
FROM golang:1.22-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# RUN apk --no-cache add netcat-openbsd

# Copy the rest of the source code to the working directory
COPY . .

# Build the Go application
RUN go build -o myapp .

# Start a new stage from scratch
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the built executable from the previous stage
COPY --from=builder /app/myapp .

# Copy SSL certificate files into the image
COPY server.crt server.key ./
RUN chmod +r /app/server.key
RUN chmod +r /app/server.crt

# Expose port 443
EXPOSE 443

# Command to run the executable
CMD ["./myapp"]
