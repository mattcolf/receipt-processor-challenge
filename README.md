# Receipt Processor Challenge
This repository contains an implementation of the Receipt Processor Challenge prompt. Here's a brief guided tour of the 
codebase. 
- `main.go` Entrypoint for the API server. 
- `api/` 
  - `api.go` API routes and handlers.
  - `config.go` Basic environment variable based configuration
  - `receipt.go` Types and associated methods (Receipt & Receipt Item)
  - `database.go` An instance of MemDB for storing and querying Receipts
  - `receipt_test.go` A small suite of unit tests for the Receipts type
- `directions/` The original challenge prompt
- `receipts.postman_collection` A small Postman collection used for testing the API

## Configuration
This application uses environment variables for configuration. There's no need to change these values as they default to
reasonable values to enable testing. For more information, see the `api/config.go` file.
- `API_HOSTNAME` (defaults to `0.0.0.0`)
- `API_PORT` (defaults to `8080`)
- `API_ENV` (defaults to `production`)
- `API_WRITE_TIMEOUT` (defaults to `15s`)
- `API_READ_TIMEOUT` (defaults to `15s`)

## Build & Run API
This application can be built and run using either of the following options.

### Option 1: Docker
1) Install Docker
2) Build and run the included Dockerfile. If you wish to modify any of the configuration environmental values, you can do so
when you run the Docker image in this step. 
```sh
cd receipt-processor-challenge
docker build -t receipt-processor-challenge-api .
docker run -it --rm --name receipt-processor-challenge-api -p 8080:8080
```

### Option 2: Native
1) Install Go version 1.23 or higher
2) Install all Go module dependencies
```sh
go mod download
```
3) Run the unit tests
```sh
gp test -v ./...
```
4) Build the application
```sh
go build -o bin/api
```