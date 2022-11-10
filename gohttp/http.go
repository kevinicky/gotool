package gohttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type HttpTools interface {
	JsonResponse(writer http.ResponseWriter, payload interface{}, httpStatusCode int)
	JsonValidatorError(writer http.ResponseWriter, error error)
	Request(context context.Context, url string, method string, payload interface{}) ([]byte, int, error)
	GetPagination(request *http.Request) (int, int, bool)
	CheckJsonHeader(request *http.Request) error
}

type ResponseCORSOptions struct {
	AllowOrigin  string `default:"*"`
	AllowMethods string `default:"*"`
	AllowHeader  string `default:"*"`
}

type PaginationOptions struct {
	DefaultLimit  int `default:"10"`
	DefaultOffset int `default:"0"`
}

type httpTools struct {
	responseCORSOptions ResponseCORSOptions
	paginationOptions   PaginationOptions
}

func NewHttpTools(responseCORSOptions ResponseCORSOptions, paginationOptions PaginationOptions) HttpTools {
	return &httpTools{
		responseCORSOptions,
		paginationOptions,
	}
}

func (httpTools *httpTools) JsonResponse(writer http.ResponseWriter, payload interface{}, httpStatusCode int) {
	jsonResp, _ := json.Marshal(payload)

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", httpTools.responseCORSOptions.AllowOrigin)
	writer.Header().Set("Access-Control-Allow-Methods", httpTools.responseCORSOptions.AllowMethods)
	writer.Header().Set("Access-Control-Allow-Headers", httpTools.responseCORSOptions.AllowHeader)
	writer.WriteHeader(httpStatusCode)
	_, _ = writer.Write(jsonResp)
}

func (httpTools *httpTools) JsonValidatorError(writer http.ResponseWriter, error error) {
	message := map[string]interface{}{}
	if castedObject, ok := error.(validator.ValidationErrors); ok {
		errObj := castedObject[0]

		switch errObj.Tag() {
		case "required":
			message = map[string]interface{}{"error": errObj.Field() + " is required"}
		case "android|ios":
			message = map[string]interface{}{"error": errObj.Field() + " must android, ios"}
		case "DANA|LINKAJA|OVO|SHOPEEPAY":
			message = map[string]interface{}{"error": errObj.Field() + " must DANA, LINKAJA, OVO, or SHOPEEPAY"}
		case "email":
			message = map[string]interface{}{"error": errObj.Field() + " is not valid email format"}
		case "gte":
			message = map[string]interface{}{"error": errObj.Field() + " value must be greater equal than " + errObj.Param()}
		case "gt":
			message = map[string]interface{}{"error": errObj.Field() + " value must be greater than " + errObj.Param()}
		default:
			message = map[string]interface{}{"error": "invalid input for " + errObj.Field()}
		}

		httpTools.JsonResponse(writer, message, http.StatusBadRequest)
	}
}

func (httpTools *httpTools) Request(context context.Context, url string, method string, payload interface{}) ([]byte, int, error) {
	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	req := &http.Request{}
	var err error

	if payload == nil {
		req, err = http.NewRequestWithContext(context, method, url, nil)
	} else {
		jsonData, _ := json.Marshal(payload)
		req, err = http.NewRequestWithContext(context, method, url, bytes.NewBuffer(jsonData))
	}

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	httpStatusCode := resp.StatusCode

	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, resp.Body); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return buf.Bytes(), httpStatusCode, nil
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

func (httpTools *httpTools) CheckJsonHeader(request *http.Request) error {
	headerContentType := request.Header.Get("Content-Type")

	if headerContentType != "application/json" {
		err := errors.New("invalid header content-type")

		return err
	}

	return nil
}
