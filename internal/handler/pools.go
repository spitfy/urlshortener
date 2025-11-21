// Package handler предоставляет пулы объектов для обработчиков.
package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
)

var jsonBufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func encodeJSONBuffered(w io.Writer, data interface{}) error {
	buf := jsonBufferPool.Get().(*bytes.Buffer)
	defer jsonBufferPool.Put(buf)
	buf.Reset()

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(data); err != nil {
		return err
	}

	_, err := buf.WriteTo(w)
	return err
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 4096))
	},
}

func readBodyLimited(r io.Reader, maxSize int64) ([]byte, error) {
	if r == nil {
		return nil, io.EOF
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)
	buf.Reset()

	lr := io.LimitReader(r, maxSize)

	_, err := io.Copy(buf, lr)
	if err != nil {
		return nil, err
	}

	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())

	return result, nil
}
