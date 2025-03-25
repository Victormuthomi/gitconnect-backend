# Step 1: Build the Go binary
FROM golang:1.23.2-alpine AS build

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go Modules files and download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Step 2: Create a smaller image to run the app
FROM alpine:3.18

# Set the working directory
WORKDIR /root/

# Copy the pre-built binary from the previous stage
COPY --from=build /app/main .

# Ensure the binary is executable
RUN chmod +x main

# Expose port 8080
EXPOSE 8080

# Run the Go binary
CMD ["./main"]

