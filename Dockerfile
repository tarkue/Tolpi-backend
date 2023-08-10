# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang:1.12-alpine base image
FROM golang:1.18

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/tolpi

# Expose port 8080 to the outside world
EXPOSE 3000

# Run the executable
ENTRYPOINT ["./main"] 
