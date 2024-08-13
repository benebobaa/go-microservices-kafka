# Build stage
FROM golang:1.22.6-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest  

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose port 5000
EXPOSE 5000

# Run the binary
CMD ["./main"]
