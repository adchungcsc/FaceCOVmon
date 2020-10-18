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
  fmt.Println(request.HTTPMethod)
  switch request.HTTPMethod {
  case "GET":
    return handleGet(request)
  case "PUT":
    return handlePut(request)
  default:
    return clientError(http.StatusMethodNotAllowed)
  }
}

func handleGet(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  cameraID := request.QueryStringParameters["cameraID"]
  date := request.QueryStringParameters["date"]
  fmt.Println(cameraID)
  fmt.Println(date)
  tenantData, _ := getTenantData(&cameraID, &date)
  //if err != nil {
  //  return serverError(err)
  //}
  if tenantData == nil {
    return clientError(http.StatusNotFound)
  }
  js, err := json.Marshal(tenantData)
  if err != nil {
    return serverError(err)
  }
  if len(js) == 0{
    //return empty json
    return events.APIGatewayProxyResponse{
      StatusCode: http.StatusOK,
      Body:       "",
    }, nil
  }

  return events.APIGatewayProxyResponse{
    StatusCode: http.StatusOK,
    Body:       string(js),
  }, nil
}

func handlePut(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  cameraID := request.QueryStringParameters["cameraID"]
  image := request.QueryStringParameters["base64Image"]
  fmt.Println(cameraID)
  fmt.Println(image)

  //Validate inputs
  if cameraID == "" || image == "" {
    return clientError(http.StatusUnprocessableEntity)
  }

  err := uploadData(cameraID, image)
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