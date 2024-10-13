# Start by building the application.
FROM golang:1.22-alpine3.19 AS build
LABEL stage=dockerbuilder
WORKDIR /app

# Copy go.mod and go.sum first to take advantage of Docker caching.
COPY go.mod go.sum ./

# Download dependencies.
RUN go mod download

# copy code
COPY . .

# Build the binary with CGO disabled for compatibility with Alpine.
RUN CGO_ENABLED=0 GOOS=linux go build -o apps main.go

# Now copy it into our base image.
FROM alpine:3.18

# Copy bin file
WORKDIR /app
COPY --from=build /app/apps /app/apps
RUN mkdir /app/logs

EXPOSE 8080
ENTRYPOINT ["/app/apps"]