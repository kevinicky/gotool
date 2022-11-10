package gohttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io"
	"net/http"
	"reflect"
)

type ResponseCORSOptions struct {
	AllowOrigin  string `default:"*"`
	AllowMethods string `default:"*"`
	AllowHeader  string `default:"*"`
}

func defaultCORSOptions(options ResponseCORSOptions) ResponseCORSOptions {
	typ := reflect.TypeOf(options)

	if options.AllowOrigin == "" {
		f, _ := typ.FieldByName("AllowOrigin")
		options.AllowOrigin = f.Tag.Get("default")
	}

	if options.AllowMethods == "" {
		f, _ := typ.FieldByName("AllowMethods")
		options.AllowMethods = f.Tag.Get("default")
	}

	if options.AllowHeader == "" {
		f, _ := typ.FieldByName("AllowHeader")
		options.AllowHeader = f.Tag.Get("default")
	}

	return options
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

func (httpTools *httpTools) CheckJsonHeader(request *http.Request) error {
	headerContentType := request.Header.Get("Content-Type")

	if headerContentType != "application/json" {
		err := errors.New("invalid header content-type")

		return err
	}

	return nil
}
