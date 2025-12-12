package entity

import (
	"net/http"
	"net/url"
)

type RequestModel struct {
	ID         string
	Method     string
	ClientIP   string
	URL        *url.URL
	Headers    http.Header
	Body       []byte
	ReceivedAt int64
}
