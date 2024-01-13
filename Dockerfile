# docs - https://docs.docker.com/language/golang/build-images/


# Base image
FROM golang:latest

# Docker metadata
LABEL maintainer = "mkotkov"
LABEL description = "Dockerized version of forum"
LABEL version = "1.0"

# Set destination for COPY
WORKDIR /app

# Copy Copy go.mod and go.sum files to the working directory.
COPY go.mod ./ 

# Download all dependencies in the image.
RUN go mod download 

# Copy the source from the current directory to image.
COPY . .

# Compile the application. Returns a binary called ascii-art-web-dockerize into the working directory.
RUN CGO_ENABLED=1 GOOS=linux go build -o /forum
EXPOSE 3000

# tell Docker what command to execute when our image is used to start a container
CMD ["/forum"]