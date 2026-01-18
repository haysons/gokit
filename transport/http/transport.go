package http

import (
	"net/http"

	"github.com/haysons/gokit/transport"
)

type Transporter interface {
	transport.Transporter
	Request() *http.Request
	PathTemplate() string
	Response() http.ResponseWriter
}
