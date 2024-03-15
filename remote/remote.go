package remote

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"../logging"
)

var logger = logging.Logger{}

func Server() *httptest.Server {
	//logging middleware captures any server requests first
	server := httptest.NewServer(logging.Middleware(http.HandlerFunc(remoteHandler)))
	return server
}

func remoteHandler(respWriter http.ResponseWriter, request *http.Request) {
	ctx := request.Context() //this ctx was given a traceId

	//this log should display the "upstream" traceId
	logger.Log(ctx, "This log should have a traceId prefix!")

	//we can invoke other "services", passing down the ctx
	result := appendService(ctx, "wow", "amazing!")

	respWriter.Write([]byte(result))
}

func appendService(ctx context.Context, str1, str2 string) string {
	//and the traceId will follow any logging
	logger.Log(ctx, fmt.Sprintf("Append service invoked with 2 params: %s - %s", str1, str2))
	return str1 + str2
}
