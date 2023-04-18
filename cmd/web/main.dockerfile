# Use an official Golang runtime as a parent image
FROM golang:1.16-alpine

# Set the working directory to /app
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./


# Download Go modules
RUN go mod download

# Copy the rest of the application code to the container
COPY . .

# Build the Go binary
RUN go build -o my-app ./cmd/web

# Expose port 4000 for the application to listen on
EXPOSE 4000

# Start the application when the container is started
CMD [ "./my-app" ]
