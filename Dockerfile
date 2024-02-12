# Start from the golang base image with Alpine
FROM golang:1.20-alpine as builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git gcc musl-dev sqlite-dev

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# CGO must be enabled for SQLite
ENV CGO_ENABLED=1
ENV GIN_MODE=release

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest  

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Declare the volume
VOLUME ["/opt/clodevo/data"]

EXPOSE 8080
EXPOSE 9090

# Command to run the executable
CMD ["./main"]
