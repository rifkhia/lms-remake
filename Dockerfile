# Use the official Golang image as a base image
FROM golang:1.21.3

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod .
COPY go.sum .

# Download and install Go module dependencies
RUN go mod download

# Copy the entire project to the working directory
COPY . .

# Build the Go application
RUN go build -o main ./cmd/server

# Expose the port on which the application will run
EXPOSE 8080

# Command to run the executable
CMD ["./main"]