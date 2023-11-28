# Sportujspolu

This is a web application written in GO

## Installation

To use this project, you need to have Go installed on your machine. Then, you can follow these steps:

1. Install the dependencies: `go mod download`
2. Run the application: `go run routes.go`

## Usage

Once the application is running, you can access the "events" route by navigating to http://localhost:3001/api/v1/events.

## Swaggo Yaml

To generate swaggo yaml use

```bash
swag init -g routes.go --ot yaml
```

It gets generate to docs/swagger.yaml
