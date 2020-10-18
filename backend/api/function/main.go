package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "os"

  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)


func router(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  switch request.HTTPMethod {
  case "GET":
    return handleGet(request)
  case "PUT":
    return handlePost(request)
  default:
    return clientError(http.StatusMethodNotAllowed)
  }
}

func handleGet(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  //Get location based on name.
  locationName := request.QueryStringParameters["locationName"]
  fmt.Println(locationName)
  location, err := getDestination(locationName)
  if err != nil {
    return serverError(err)
  }
  if location == nil {
    return clientError(http.StatusNotFound)
  }
  js, err := json.Marshal(location)
  if err != nil {
    return serverError(err)
  }
  return events.APIGatewayProxyResponse{
    StatusCode: http.StatusOK,
    Body:       string(js),
  }, nil
}

func validateNewLocation(location Location) bool{
  return location.LocationName == "" || location.Description == ""
}

func handlePost(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  if request.Headers["content-type"] != "application/json" && request.Headers["Content-Type"] != "application/json" {
    return clientError(http.StatusNotAcceptable)
  }

  location := new(Location)
  err := json.Unmarshal([]byte(request.Body), location)
  fmt.Println(location)

  if err != nil {
    return clientError(http.StatusUnprocessableEntity)
  }

  if !validateNewLocation(*location) {
    return clientError(http.StatusBadRequest)
  }

  err = insertDestination(location)
  if err != nil {
    return serverError(err)
  }

  return events.APIGatewayProxyResponse{
    StatusCode: 201,
  }, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
  errorLogger.Println(err.Error())

  return events.APIGatewayProxyResponse{
    StatusCode: http.StatusInternalServerError,
    Body:       http.StatusText(http.StatusInternalServerError),
  }, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
  return events.APIGatewayProxyResponse{
    StatusCode: status,
    Body:       http.StatusText(status),
  }, nil
}

func main() {
  lambda.Start(router)
}