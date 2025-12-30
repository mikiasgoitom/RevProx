package entity

import (
	"net/http"
)

type ResponseModel struct {
	ID          string
	Status      int
	Headers      http.Header
	Body        []byte
	GeneratedAt int64
	Cacheable   bool
}
