# Use an official Golang image
FROM golang:1.20 as builder

# Set the working directory
WORKDIR /app

# Copy the Go modules and dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the source code
COPY . ./

# Build the application
RUN go build -o main .

# Use a lightweight image for production
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
