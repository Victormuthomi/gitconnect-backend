# Step 1: Build the Go binary
FROM golang:1.23.2-alpine AS build

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o main .

# Step 2: Create a smaller image to run the app
FROM alpine:3.18

WORKDIR /root/

# Copy the pre-built binary from the build stage
COPY --from=build /app/main .

RUN chmod +x main

EXPOSE 8080

CMD ["./main"]

