package valueobject

import "net/http"

type HeaderSet  struct {
	Headers http.Header
}