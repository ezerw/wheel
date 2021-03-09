FROM node:14-alpine as frontend

# Container app dir
WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

COPY ./client ./

RUN yarn
RUN yarn global add react-scripts@4.0.3

RUN yarn build

# Start from golang base image
FROM golang:alpine as api

# ENV GO111MODULE=on

# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod sum and env files
COPY go.mod go.sum .env ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Add files
ADD db ./db
ADD handlers ./handlers
ADD models ./models
COPY main.go .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wheel .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add bash ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage.
COPY --from=api /app/wheel .
COPY --from=api /app/.env .
COPY --from=frontend /app/build ./client/build

# Expose port 8080 to the outside world
EXPOSE 8080

#Command to run the executable
CMD ["./wheel"]