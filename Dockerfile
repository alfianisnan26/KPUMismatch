### Step 1: Build stage
FROM golang:1.21-alpine as builder

WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main ./app/main.go

RUN apk add -U --no-cache ca-certificates

###
## Step 2: Runtime stage
FROM scratch

# Copy only the binary from the build stage to the final image
COPY --from=builder /app/main /

# Set the entry point for the container
ENTRYPOINT ["/main"]

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/