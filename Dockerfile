# Step 1: Build the Go binary
FROM golang:1.23.2-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules file
COPY go.mod go.sum ./

# Download all the dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Step 2: Create a smaller image to run the app
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/main .

# Copy any necessary configuration files (for example, your .env file)
COPY .env ./

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the Go binary
CMD ["./main"]

