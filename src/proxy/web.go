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
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const httpPort = 31596

func webListen() error {
	http.HandleFunc("/favicon.ico", webServeFavicon)
	http.HandleFunc("/header.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		dataBytes, err := base64.StdEncoding.DecodeString(assetHeader)
		if err != nil {
			log.Printf("[WARN] %s", err.Error())
			webHandleError(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(dataBytes)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pathSplit := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")
		// home
		if len(pathSplit) == 1 && pathSplit[0] == "" {
			webServeB64(assetIndex, http.StatusOK, w)
			return
		}
		// locate proxy user and serve proxy
		proxyUser := getProxyUser(pathSplit[0])
		if proxyUser == nil {
			webServeB64(assetError404UserNotFound, http.StatusNotFound, w)
			return
		}
		webServeProxy(proxyUser, w, r)
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

func webServeProxy(u *ProxyUser, w http.ResponseWriter, r *http.Request) {
	pathSplit := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")
	if len(pathSplit) < 2 {
		webServeB64(assetError404General, http.StatusNotFound, w)
		return
	}
	reqId, err := u.handleRequest(r)
	if err != nil {
		log.Printf("[WARN] #%d handleRequest :: %s", reqId, err.Error())
		webServeB64(assetError500General, http.StatusInternalServerError, w)
		return
	}
	rawResp, err := u.responseWait(reqId)
	if err != nil {
		log.Printf("[WARN] #%d responseWait :: %s", reqId, err.Error())
		webServeB64(assetError500General, http.StatusInternalServerError, w)
		return
	}
	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(rawResp)), r)
	if err != nil {
		log.Printf("[WARN] #%d response read :: %s", reqId, err.Error())
		webServeB64(assetError500General, http.StatusInternalServerError, w)
		return
	}
	for k, vs := range resp.Header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("[WARN] response write :: %s", err.Error())
		webServeB64(assetError500General, http.StatusInternalServerError, w)
		return
	}

}

func webServeB64(data string, status int, w http.ResponseWriter) {
	dataBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Printf("[WARN] webServeB64 :: %s", err.Error())
		webHandleError(w)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	w.Write(dataBytes)
}

func webServeFavicon(w http.ResponseWriter, r *http.Request) {
	dataBytes, err := base64.StdEncoding.DecodeString(assetFavicon)
	if err != nil {
		log.Printf("[WARN] webServeFavicon :: %s", err.Error())
		webHandleError(w)
		return
	}
	w.Header().Set("Content-Type", "image/ico")
	w.WriteHeader(http.StatusOK)
	w.Write(dataBytes)
}

func webHandleError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("An unknown error occured."))
}
