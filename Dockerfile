FROM golang:1.23 AS base

WORKDIR /usr/src/app

# Install Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install Application
COPY . .

# Run Unit Tests
RUN go test ./...

# Build Application
RUN go build -v -o /usr/local/bin/receipt-processor-challenge-api ./main.go
RUN chmod +x /usr/local/bin/receipt-processor-challenge-api

# Run Application
EXPOSE 8080
CMD ["receipt-processor-challenge-api"]
