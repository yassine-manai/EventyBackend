# Step 1: Use official Golang image as the base image (adjusted to the actual Go version)
FROM golang:1.23-alpine

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy go.mod and go.sum files to download dependencies
COPY go.mod . 
COPY go.sum .

# Step 4: Download the dependencies
RUN go mod download

RUN apk add --no-cache tzdata

# Step 6: ENV VARS
ENV TZ=UTC

# Step 7: Copy the source code to the container
COPY . .

# Step 8: Build the Go application
RUN go build -o bin/app .

# Step 9: Set the entry point for the container to run the application
ENTRYPOINT [ "/app/bin/app" ]
