package urlkey

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateSlug(t *testing.T) {
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
			name: "Valid Uppercase and lowercase URL to compress",
			url:  "http://WWW.Test.COM/",
			want: "4b23",
		},
		{
			name: "Non-Valid URL",
			url:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSlug(tt.url)
			t.Logf("Generated Key: %s", got) // Debugging output
			assert.Equal(t, tt.want, got)
		})
	}
}
