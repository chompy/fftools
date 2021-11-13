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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pathSplit := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")
		proxyUser := getProxyUser(pathSplit[0])
		if proxyUser == nil {
			webServeB64(assetError404UserNotFound, http.StatusNotFound, w)
			return
		}
		webServeProxy(proxyUser, w, r)
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

func webServeProxy(u *proxyUser, w http.ResponseWriter, r *http.Request) {
	pathSplit := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")
	if len(pathSplit) < 2 {
		webServeB64(assetError404General, http.StatusNotFound, w)
		return
	}
	reqId, err := u.handleRequest(r)
	if err != nil {
		log.Printf("[WARN] handleRequest :: %s", err.Error())
		webServeB64(assetError500General, http.StatusInternalServerError, w)
		return
	}
	rawResp, err := u.responseWait(reqId)
	if err != nil {
		log.Printf("[WARN] responseWait :: %s", err.Error())
		webServeB64(assetError500General, http.StatusInternalServerError, w)
		return
	}

	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(rawResp)), r)
	if err != nil {
		log.Printf("[WARN] response read :: %s", err.Error())
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
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "image/ico")
	w.WriteHeader(http.StatusOK)
	w.Write(dataBytes)
}

func webHandleError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("An unknown error occured."))
}
