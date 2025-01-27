package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "Valid URL with the trailing '/' to compress",
			url:  "http://www.test.com/",
			want: "4b23",
		},
		{
			name: "Non-Valid URL",
			url:  "/",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, generateKey(tt.url))
		})
	}
}
