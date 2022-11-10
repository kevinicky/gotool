package gohttp

import (
	"net/http"
	"reflect"
	"strconv"
)

type PaginationOptions struct {
	DefaultLimit  int `default:"10"`
	DefaultOffset int `default:"0"`
}

func defaultPaginationOptions(options PaginationOptions) PaginationOptions {
	typ := reflect.TypeOf(options)
	if options.DefaultOffset == 0 {
		f, _ := typ.FieldByName("DefaultOffset")
		value, _ := strconv.Atoi(f.Tag.Get("default"))
		options.DefaultOffset = value
	}
	if options.DefaultLimit == 0 {
		f, _ := typ.FieldByName("DefaultLimit")
		value, _ := strconv.Atoi(f.Tag.Get("default"))
		options.DefaultOffset = value
	}

	return options
}

func (httpTools *httpTools) GetPagination(request *http.Request) (int, int, bool) {
	isPage := true
	offset := 0

	page, err := strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil {
		isPage = false
	}

	limit, err := strconv.Atoi(request.URL.Query().Get("limit"))
	if err != nil {
		limit = httpTools.paginationOptions.DefaultLimit
	}

	if isPage {
		offset = (page - 1) * limit
	} else {
		offset, err = strconv.Atoi(request.URL.Query().Get("offset"))
		if err != nil {
			offset = httpTools.paginationOptions.DefaultOffset
		}
	}

	return limit, offset, isPage
}
