package helpers

import "net/http"

type HttpManager struct {
	client *http.Client
}

func NewHttpManager() *HttpManager {
	return &HttpManager{client: &http.Client{}}
}
