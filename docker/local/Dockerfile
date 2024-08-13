# Use the Golang base image for building the application
FROM golang:1.21-alpine AS builder

# Install required packages
RUN apk add --no-cache make git

# Set the working directory
WORKDIR /blade

# Add go.mod and go.sum files to the working directory
# Ensure these files are present in the build context
ADD go.mod go.sum ./
RUN go mod download

# Copy all files to the working directory
COPY . .

# Build the application using the Makefile
# Added VERBOSE=1 to increase build verbosity
RUN make build VERBOSE=1

# Use the Alpine image for running the application
FROM alpine:latest AS runner

# Install necessary packages
RUN apk --no-cache add ca-certificates jq

# Set the working directory
WORKDIR /blade

# Print the current working directory
RUN ls -al
RUN pwd

# Copy the built binary from the builder stage
COPY --from=builder /blade/blade ./

# Copy the blade.sh script to the working directory
COPY docker/local/blade.sh ./

# Expose ports for JSON-RPC, libp2p, and gRPC
EXPOSE 8545 9632 1478 5001

# Set the entrypoint to run the blade.sh script
ENTRYPOINT ["./blade.sh"]

