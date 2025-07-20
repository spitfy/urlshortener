package handler

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_gzipMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		acceptEnc    string
		contentEnc   string
		requestBody  io.Reader
		wantGzipResp bool
		wantGzipReq  bool
		wantStatus   int
	}{
		{"both gzip", "gzip", "gzip", gzipCompress([]byte("hello")), true, true, http.StatusOK},
		{"no gzip", "", "", strings.NewReader("hello"), false, false, http.StatusOK},
		{"only request gzip", "", "gzip", gzipCompress([]byte("hello")), false, true, http.StatusOK},
		{"only response gzip", "gzip", "", strings.NewReader("hello"), true, false, http.StatusOK},
		{"accept encoding includes gzip", "gzip, deflate", "", strings.NewReader("hello"), true, false, http.StatusOK},
		{"content encoding unknown", "gzip", "deflate", strings.NewReader("hello"), true, false, http.StatusOK},
		{"invalid gzip body", "", "gzip", strings.NewReader("invalid data"), false, false, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := gzipMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Errorf("error reading request body: %v", err)
				}
				if tt.wantStatus == http.StatusInternalServerError {
					w.WriteHeader(tt.wantStatus)
					return
				}
				_, _ = w.Write(body)
			}))

			req := httptest.NewRequest("POST", "/", tt.requestBody)
			if tt.acceptEnc != "" {
				req.Header.Set("Accept-Encoding", tt.acceptEnc)
			}
			if tt.contentEnc != "" {
				req.Header.Set("Content-Encoding", tt.contentEnc)
			}

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, rr.Code)
			}

			ce := rr.Header().Get("Content-Encoding")
			if tt.wantGzipResp && ce != "gzip" {
				t.Errorf("expected response Content-Encoding: gzip, got %q", ce)
			}
			if !tt.wantGzipResp && ce != "" {
				t.Errorf("expected response Content-Encoding: <empty>, got %q", ce)
			}
		})
	}
}

func gzipCompress(data []byte) io.Reader {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, _ = gw.Write(data)
	_ = gw.Close()
	return bytes.NewReader(buf.Bytes())
}
