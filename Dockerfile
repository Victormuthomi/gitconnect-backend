# Step 1: Build the Go binary
FROM golang:1.23.2-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# Step 2: Create a smaller image to run the app
FROM alpine:3.18

WORKDIR /root/

COPY --from=build /app/main .

RUN chmod +x main

EXPOSE 8080

CMD ["./main"]

