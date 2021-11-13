package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type webProxyResponseWriter struct {
	r http.Response
}

func (w *webProxyResponseWriter) Header() http.Header {
	return w.r.Header
}

func (w *webProxyResponseWriter) Write(data []byte) (int, error) {
	w.r.Body = ioutil.NopCloser(bytes.NewReader(data))
	return len(data), nil
}

func (w *webProxyResponseWriter) WriteHeader(statusCode int) {
	w.r.StatusCode = statusCode
}

func newWebProxyResponseWriter() *webProxyResponseWriter {
	resp := http.Response{
		StatusCode:    http.StatusOK,
		Status:        "200 OK",
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		ContentLength: 0,
		Body:          nil,
	}
	return &webProxyResponseWriter{
		r: resp,
	}
}
