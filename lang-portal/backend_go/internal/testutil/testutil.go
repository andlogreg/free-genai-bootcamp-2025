package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// PerformRequest performs an HTTP request and returns the response
func PerformRequest(t *testing.T, r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		assert.NoError(t, err)
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, path, reqBody)
	assert.NoError(t, err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ParseResponse parses the response body into the given struct
func ParseResponse(t *testing.T, w *httptest.ResponseRecorder, v interface{}) {
	err := json.Unmarshal(w.Body.Bytes(), v)
	assert.NoError(t, err)
}

// AssertStatusCode asserts that the response has the expected status code
func AssertStatusCode(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	assert.Equal(t, expected, w.Code)
}

// AssertJSONResponse asserts that the response has the expected JSON body
func AssertJSONResponse(t *testing.T, w *httptest.ResponseRecorder, expected interface{}) {
	var actual interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actual)
	assert.NoError(t, err)

	expectedJSON, err := json.Marshal(expected)
	assert.NoError(t, err)

	actualJSON, err := json.Marshal(actual)
	assert.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(actualJSON))
}
