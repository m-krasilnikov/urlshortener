package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func WithGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.EqualFold(
			r.Header.Get("Content-Encoding"),
			"gzip",
		) {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(
					w,
					"Bad Request",
					http.StatusBadRequest,
				)
				return
			}

			defer reader.Close()

			r.Body = reader
		}

		if !strings.Contains(
			r.Header.Get("Accept-Encoding"),
			"gzip",
		) {
			next.ServeHTTP(w, r)
			return
		}

		gzipResponse := &gzipResponseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(gzipResponse, r)

		if gzipResponse.writer != nil {
			_ = gzipResponse.writer.Close()
		}
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	writer      *gzip.Writer
	wroteHeader bool
	compress    bool
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}

	w.wroteHeader = true

	contentType := w.Header().Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") ||
		strings.HasPrefix(contentType, "text/html") {

		w.compress = true

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")
	}

	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	if !w.compress {
		return w.ResponseWriter.Write(data)
	}

	if w.writer == nil {
		w.writer = gzip.NewWriter(w.ResponseWriter)
	}

	return w.writer.Write(data)
}
