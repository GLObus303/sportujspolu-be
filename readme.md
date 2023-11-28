# Sportujspolu

This is a web application written in GO

## Installation

To use this project, you need to have Go installed on your machine. Then, you can follow these steps:

1. Install the dependencies: `go mod download`
2. Run the application: `go run main.go`

## Usage

Once the application is running, you can access the "events" route by navigating to http://localhost:3001/api/v1/events or even better - using our [postman](https://sportujspolu.postman.co/workspace/Team-Workspace~2f2621b5-b6ff-41f3-8472-28c07536fc3f/overview) .

## Local development with gow

There is a possibility to run the local server with tool called `gow` that does hot reload on any change.

Check their [documentation](https://github.com/mitranim/gow) for more info.

## Swaggo Yaml

To generate swaggo yaml use

```bash
swag init -g routes.go --ot yaml
```

It gets generate to docs/swagger.yaml
