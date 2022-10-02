package gohttp

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type HTTPUtil interface {
	JsonResponse(w http.ResponseWriter, message interface{}, httpStatusCode int)
	JsonValidatorError(w http.ResponseWriter, err error)
	Request(ctx context.Context, url string, method string, payload interface{}) ([]byte, int, error)
	GetPagination(r *http.Request) (int, int)
}

type httputil struct{}

func NewHttpTool() HTTPUtil {
	return &httputil{}
}

func (h *httputil) JsonResponse(w http.ResponseWriter, message interface{}, httpStatusCode int) {
	jsonResp, _ := json.Marshal(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	_, _ = w.Write(jsonResp)
}

func (h *httputil) JsonValidatorError(w http.ResponseWriter, err error) {
	message := map[string]interface{}{}
	if castedObject, ok := err.(validator.ValidationErrors); ok {
		err := castedObject[0]

		switch err.Tag() {
		case "required":
			message = map[string]interface{}{"message": err.Field() + " is required"}
		case "android|ios":
			message = map[string]interface{}{"message": err.Field() + " must android, ios"}
		}

		h.JsonResponse(w, message, http.StatusBadRequest)
	}
}

func (h *httputil) Request(ctx context.Context, url string, method string, payload interface{}) ([]byte, int, error) {
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

func (h *httputil) GetPagination(r *http.Request) (int, int) {
	isWebsite := true
	offset := 0

	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		isWebsite = false
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		limit = 10
	}

	if isWebsite {
		offset = (page - 1) * limit
	} else {
		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))

		if err != nil {
			offset = 0
		}
	}

	return limit, offset
}
