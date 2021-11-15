/*
This file is part of FFTools.

FFTools is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

FFTools is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with FFTools.  If not, see <https://www.gnu.org/licenses/>.
*/

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
