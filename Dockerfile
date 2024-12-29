# Start with a lightweight Go image
FROM golang:1.20 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . ./

# Build the Go CLI as a binary with a custom name
RUN go build -o cloudcents .

# Use a minimal base image for the final container
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/cloudcents .

# Command to run the CLI
ENTRYPOINT ["./cloudcents"]
