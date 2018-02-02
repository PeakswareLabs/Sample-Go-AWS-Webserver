package main

import (
	"log"
	"net/http"
	"os"

	"github.com/PeakswareLabs/Go-Webserver/app"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
)

var reqAccessor *core.RequestAccessor

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqAccessor = &core.RequestAccessor{}
	handler := app.Create()
	req, err := reqAccessor.ProxyEventToHTTPRequest(request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusGatewayTimeout}, err
	}
	respWriter := core.NewProxyResponseWriter()

	// handle the HTTP request
	handler.ServeHTTP(http.ResponseWriter(respWriter), req)

	proxyResponse, err := respWriter.GetProxyResponse()
	if err != nil {
		log.Println("Error while generating proxy response")
		log.Println(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusGatewayTimeout}, err
	}
	return proxyResponse, nil
}

func main() {
	if os.Getenv("ENVIRONMENT") == "dev" {
		lambda.Start(Handler)
	} else {
		handler := app.Create()
		if err := http.ListenAndServe(":65010", handler); err != nil {
			panic(err)
		}
	}

}
