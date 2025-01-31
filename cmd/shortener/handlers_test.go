package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRouter(t *testing.T) {
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
				body:   "The body should contain URL",
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
				body:   "URL is not found",
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
				body:   "Missing URL key",
			},
		},
	}

	client := resty.New()
	for _, tt := range tests {

		url := host + port + tt.target

		if tt.method == "GET" {
			res, err := client.R().Get(url)

			require.NoError(t, err)

			assert.Equal(t, tt.want.code, res.StatusCode())
			assert.Equal(t, tt.want.header, res.Header().Get("Location"))

		} else if tt.method == "POST" {
			res, err := client.R().SetBody(tt.body).Post(url)

			require.NoError(t, err)

			assert.Equal(t, tt.want.code, res.StatusCode())
			assert.Equal(t, tt.want.body, string(res.Body()))
		}
	}
}
