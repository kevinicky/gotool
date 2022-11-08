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
	JsonResponse(w http.ResponseWriter, message interface{}, httpStatusCode int)
	JsonValidatorError(w http.ResponseWriter, err error)
	Request(ctx context.Context, url string, method string, payload interface{}) ([]byte, int, error)
	GetPagination(r *http.Request) (int, int, bool)
	CheckJsonHeader(request *http.Request) (err error)
}

type httpTools struct{}

func NewHttpTools() HttpTools {
	return &httpTools{}
}

func (h *httpTools) JsonResponse(w http.ResponseWriter, message interface{}, httpStatusCode int) {
	jsonResp, _ := json.Marshal(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	_, _ = w.Write(jsonResp)
}

func (h *httpTools) JsonValidatorError(w http.ResponseWriter, err error) {
	message := map[string]interface{}{}
	if castedObject, ok := err.(validator.ValidationErrors); ok {
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

		h.JsonResponse(w, message, http.StatusBadRequest)
	}
}

func (h *httpTools) Request(ctx context.Context, url string, method string, payload interface{}) ([]byte, int, error) {
	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	req := &http.Request{}
	var err error

	if payload == nil {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	} else {
		jsonData, _ := json.Marshal(payload)
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonData))
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

func (h *httpTools) GetPagination(r *http.Request) (int, int, bool) {
	isPage := true
	offset := 0

	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		isPage = false
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		limit = 10
	}

	if isPage {
		offset = (page - 1) * limit
	} else {
		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))

		if err != nil {
			offset = 0
		}
	}

	return limit, offset, isPage
}

func (h *httpTools) CheckJsonHeader(request *http.Request) (err error) {
	headerContentType := request.Header.Get("Content-Type")

	if headerContentType != "application/json" {
		err = errors.New("invalid header content-type")

		return
	}

	return
}
