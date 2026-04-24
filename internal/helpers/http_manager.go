package helpers

import (
	"net/http"
	"time"
)

const defaultTimeout = 10 * time.Second

// HTTPManager выполняет HTTP-запросы к API ASISA.
type HTTPManager struct{ client *http.Client }

// NewHTTPManager возвращает HTTPManager с разумным таймаутом по умолчанию.
func NewHTTPManager() *HTTPManager {
	return &HTTPManager{client: &http.Client{Timeout: defaultTimeout}}
}
