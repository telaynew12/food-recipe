# Use a base Go image
FROM golang:1.23


# Set working directory
WORKDIR /app

# Copy Go files
COPY . .

# Install dependencies and build the Go app
RUN go mod tidy && go build -o main .

# Expose the application port
EXPOSE 8081

# Run the application
CMD ["./main"]
