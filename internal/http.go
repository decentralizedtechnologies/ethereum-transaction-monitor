package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/canya-com/canwork-database-client/model"
	appenginelog "google.golang.org/appengine/log"
)

const (
	// GET : string GET
	GET = "GET"
	// POST : string POST
	POST = "POST"
)

// BadRequest : describes the data structure for wrong requests
type BadRequest struct {
	Message     string              `json:"message"`
	Transaction model.Transaction   `json:"transaction"`
	Context     context.Context     `json:"context"`
	Writer      http.ResponseWriter `json:"writer"`
}

// OnBadRequest : writes a json error http response
func (response *BadRequest) OnBadRequest(httpStatusCode int) {
	appenginelog.Errorf(response.Context, response.Message)
	output, _ := json.Marshal(response)
	response.Writer.WriteHeader(httpStatusCode)
	response.Writer.Write(output)
}

// NewInvalidMethodMessage : returns invalid method message
func NewInvalidMethodMessage(method string, request *http.Request) string {
	return fmt.Sprintf("Incorrect HTTP method: %s, should be a %s request", request.Method, method)
}
