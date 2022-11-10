package gohttp

import (
	"context"
	"net/http"
)

type HttpTools interface {
	JsonResponse(writer http.ResponseWriter, payload interface{}, httpStatusCode int)
	JsonValidatorError(writer http.ResponseWriter, error error)
	Request(context context.Context, url string, method string, payload interface{}) ([]byte, int, error)
	GetPagination(request *http.Request) (int, int, bool)
	CheckJsonHeader(request *http.Request) error
}

type httpTools struct {
	responseCORSOptions ResponseCORSOptions
	paginationOptions   PaginationOptions
}

func NewHttpTools(responseCORSOptions ResponseCORSOptions, paginationOptions PaginationOptions) HttpTools {
	return &httpTools{
		defaultCORSOptions(responseCORSOptions),
		defaultPaginationOptions(paginationOptions),
	}
}
