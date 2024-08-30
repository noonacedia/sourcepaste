package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/noonacedia/sourcepaste/internal/assert"
	"github.com/noonacedia/sourcepaste/internal/models/mocks"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping/")
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "It's OK")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  fmt.Sprintf("/snippets/%v/", mocks.MockSnippet.ID),
			wantCode: http.StatusOK,
			wantBody: mocks.MockSnippet.Content,
		},
		{
			name:     "Non-existent ID",
			urlPath:  fmt.Sprintf("/snippets/%v/", 0),
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  fmt.Sprintf("/snippets/%v/", -1),
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  fmt.Sprintf("/snippets/%v/", 1.5),
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  fmt.Sprintf("/snippets/%v/", "merde"),
			wantCode: http.StatusNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, _, body := ts.get(t, test.urlPath)
			assert.Equal(t, code, test.wantCode)
			if test.wantBody != "" {
				assert.Contains(t, body, test.wantBody)
			}
		})
	}
}
