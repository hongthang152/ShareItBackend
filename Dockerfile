# Start from golang v1.8 base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Thang Nguyen <nguyenhongthang1998@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/hongthang152/ShareItBackend

# Copy this directory to the container direcotry
COPY . .

# Install the package
RUN go get -d -v ./...

# Build Go
RUN go build

EXPOSE 8000

ENTRYPOINT ./ShareItBackend

