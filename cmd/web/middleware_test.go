package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/noonacedia/sourcepaste/internal/assert"
)

func TestSecureHeadersMiddleware(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	testRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) })
	secureHeaders(mockHandler).ServeHTTP(responseRecorder, testRequest)
	res := responseRecorder.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, string(body), "OK")
	assert.Equal(t, res.Header.Get("Content-Security-Policy"), "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
	assert.Equal(t, res.Header.Get("Referrer-Policy"), "origin-when-cross-origin")
	assert.Equal(t, res.Header.Get("X-Content-Type-Options"), "nosniff")
	assert.Equal(t, res.Header.Get("X-Frame-Options"), "deny")
	assert.Equal(t, res.Header.Get("X-XSS-Protection"), "0")
}
