# Start with a base Golang image
FROM golang:1.23rc2-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go script into the container
COPY . .

# Build the Go binary
RUN go build -o github-content-sync .

# Set the entry point as CMD
CMD ["./github-content-sync"]
