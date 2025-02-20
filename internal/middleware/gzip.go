package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func WithGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Check if the request is gzip encoded
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// Unzip and read it
			zr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request", http.StatusInternalServerError)
				return
			}
			defer zr.Close()
			r.Body = zr
		}

		// Check if the client accept zipped response
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {

			contentType := r.Header.Get("Content-Type")
			if strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/html") {
				// Set encoding
				w.Header().Set("Content-Encoding", "gzip")
				cw := gzip.NewWriter(w)
				defer cw.Close()
				w = &gzipResponseWriter{ResponseWriter: w, Writer: cw}
			}
		}

		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
