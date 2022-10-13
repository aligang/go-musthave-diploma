package compress

import (
	"compress/gzip"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Incorrect Encoding", http.StatusBadRequest)
				return
			} else {
				r.Body = gz
				logging.Debug("Decompression was applied")
			}
		} else {
			logging.Debug("No Compression in request header")
		}

		var writer http.ResponseWriter
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			writer = w
			logging.Debug("No Response compression will be provided")
		} else {
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			writer = gzipWriter{ResponseWriter: w, Writer: gz}
			logging.Debug("Response will be compressed")
			writer.Header().Set("Content-Encoding", "gzip")
			defer gz.Close()
		}
		next.ServeHTTP(writer, r)
	})
}
