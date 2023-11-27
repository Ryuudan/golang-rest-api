# Use an official Golang runtime as a parent image
FROM golang:1.21.4-alpine3.18

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application
RUN go build -o build ./cmd

# Expose the port the application will run on
EXPOSE 8080

# Define the command to run your application
CMD ["./build"]
