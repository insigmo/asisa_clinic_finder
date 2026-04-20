package helpers

import (
	"net/http"
	"time"
)

const defaultTimeout = 10 * time.Second

type HTTPManager struct {
	client *http.Client
}

func NewHTTPManager() *HTTPManager {
	return &HTTPManager{
		client: &http.Client{Timeout: defaultTimeout},
	}
}
