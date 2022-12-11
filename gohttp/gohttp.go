package gohttp

import (
	"context"
	"github.com/kevinicky/gotool/govalidate"
	"net/http"
)

type HttpTools interface {
	JsonResponse(writer http.ResponseWriter, payload interface{}, httpStatusCode int)
	JsonValidatorError(writer http.ResponseWriter, error error)
	RequestWithContext(context context.Context, url string, method string, payload interface{}) ([]byte, int, error)
	Request(url string, method string, payload interface{}) ([]byte, int, error)
	GetPagination(request *http.Request) (int, int, bool)
	CheckJsonHeader(request *http.Request) error
}

type httpTools struct {
	responseCORSOptions ResponseCORSOptions
	paginationOptions   PaginationOptions
	validateTools       govalidate.ValidateTools
}

// NewHttpTools create object for accessing http tool function
//
// Return http tool interface for accessing function
func NewHttpTools(responseCORSOptions ResponseCORSOptions, paginationOptions PaginationOptions) HttpTools {
	return &httpTools{
		defaultCORSOptions(responseCORSOptions),
		defaultPaginationOptions(paginationOptions),
		govalidate.NewValidateTools(),
	}
}
