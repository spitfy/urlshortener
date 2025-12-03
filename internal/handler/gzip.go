// Package handler предоставляет middleware для сжатия GZIP.
package handler

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет сжимать ответы в gzip.
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// compressReader реализует интерфейс io.ReadCloser и позволяет распаковывать gzip-сжатые запросы.
type compressReader struct {
	r  io.Reader
	zr *gzip.Reader
}

// Header возвращает HTTP-заголовки ответа.
func (cw *compressWriter) Header() http.Header {
	return cw.w.Header()
}

// Write записывает сжатые данные в ответ.
func (cw *compressWriter) Write(b []byte) (int, error) {
	return cw.zw.Write(b)
}

// WriteHeader устанавливает код статуса ответа и, если код < 300, добавляет заголовок Content-Encoding: gzip.
func (cw *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		cw.w.Header().Set("Content-Encoding", "gzip")
	}
	cw.w.WriteHeader(statusCode)
}

var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

// newCompressWriter создает новый compressWriter.
func (cw *compressWriter) Close() error {
	if cw.zw != nil {
		err := cw.zw.Close()
		gzipWriterPool.Put(cw.zw)
		cw.zw = nil
		return err
	}
	return nil
}

// newCompressReader создает новый compressReader.
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// newCompressReader создает новый compressReader.
func newCompressReader(r io.Reader) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{r: r, zr: zr}, nil
}

// Read читает распакованные данные из запроса.
func (cr *compressReader) Read(b []byte) (int, error) {
	return cr.zr.Read(b)
}

// Close закрывает gzip.Reader.
func (cr *compressReader) Close() error {
	return cr.zr.Close()
}

// gzipMiddleware проверяет поддержку gzip и применяет сжатие/распаковку к запросу/ответу.
func gzipMiddleware(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			cw := newCompressWriter(w)
			ow = cw
			defer func() {
				if err := cw.Close(); err != nil {
					log.Println(err)
				}
			}()
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer func() {
				if err := r.Body.Close(); err != nil {
					log.Println(err)
				}
			}()
		}

		h.ServeHTTP(ow, r)
	}
}
