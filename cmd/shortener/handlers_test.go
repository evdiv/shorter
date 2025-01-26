package main

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMainHandler(t *testing.T) {
	type want struct {
		code   int
		header string
		body   string
	}
	tests := []struct {
		name   string
		target string
		method string
		body   string
		want   want
	}{
		{
			name:   "POST: Positive. Valid Url for storing",
			target: "/",
			method: "POST",
			body:   "https://practicum.yandex.ru/",
			want: want{
				code:   201,
				header: "",
				body:   "http://localhost:8080/9740",
			},
		},
		{
			name:   "POST: Negative. Empty Url for storing",
			target: "/",
			method: "POST",
			body:   "",
			want: want{
				code:   400,
				header: "",
				body:   "",
			},
		},
		{
			name:   "GET: Positive. Extract Url by a key",
			target: "/9740",
			method: "GET",
			body:   "",
			want: want{
				code:   307,
				header: "https://practicum.yandex.ru/",
				body:   "",
			},
		},
		{
			name:   "GET: Negative. Extract Url by non-valid key",
			target: "/xxxxxxxxx",
			method: "GET",
			body:   "",
			want: want{
				code:   400,
				header: "",
				body:   "",
			},
		},
		{
			name:   "GET: Negative. Extract Url by empty key",
			target: "/",
			method: "GET",
			body:   "",
			want: want{
				code:   400,
				header: "",
				body:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.body))
			res := httptest.NewRecorder()

			var h MainHandler
			h.ServeHTTP(res, req)

			result := res.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.header, result.Header.Get("Location"))
			assert.Equal(t, tt.want.body, res.Body.String())
		})
	}
}
