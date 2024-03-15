package service

import (
	"context"
	"fmt"
	"net/http"
)

type Logger interface {
	Log(context.Context, string)
}

type RequestDecorator func(*http.Request) *http.Request

type SuperService struct {
	RequestDecorator RequestDecorator
	Logger           Logger
	RemoteUrl        string
}

func (svc SuperService) SomeServiceLogic(ctx context.Context, someParam string) error {
	svc.Logger.Log(ctx, "service call: "+someParam)

	//make an external server call
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, svc.RemoteUrl, nil)
	if err != nil {
		svc.Logger.Log(ctx, "error with new request with context:"+err.Error())
		return err
	}
	req = svc.RequestDecorator(req)

	//fire the remote request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	svc.Logger.Log(ctx, fmt.Sprintf("Request status: %s", response.Status))
	return nil
}
