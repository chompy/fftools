package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

const httpPort = 31596

func webListen() error {

	http.HandleFunc("/favicon.ico", webServeFavicon)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pathSplit := strings.Split(r.URL.Path, "/")
		uid := pathSplit[0]
		proxyUser := getProxyUser(uid)
		if proxyUser == nil {
			webServeB64(webNotFound, http.StatusNotFound, w)
			return
		}
		webServeProxy(proxyUser, w, r)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

func webServeProxy(u *proxyUser, w http.ResponseWriter, r *http.Request) {
	pathSplit := strings.Split(r.URL.Path, "/")
	if len(pathSplit) < 2 {
		webServeB64(webNotFound, http.StatusNotFound, w)
		return
	}
	sendReq := r.Clone(context.Background())
	sendReq.URL.Path = strings.Join(pathSplit[1:], "/")
	rawReq, err := httputil.DumpRequest(sendReq, true)
	if err != nil {
		log.Printf("[WARN] %s", err.Error())
		webServeB64(webError, http.StatusInternalServerError, w)
		return
	}
	if _, err := u.connection.Write(rawReq); err != nil {
		log.Printf("[WARN] %s", err.Error())
		webServeB64(webError, http.StatusInternalServerError, w)
		return
	}

}

func webServeB64(data string, status int, w http.ResponseWriter) {
	dataBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Printf("[WARN] %s", err.Error())
		webHandleError(w)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	w.Write(dataBytes)
}

func webServeFavicon(w http.ResponseWriter, r *http.Request) {
	dataBytes, err := base64.StdEncoding.DecodeString(webFavicon)
	if err != nil {
		log.Printf("[WARN] %s", err.Error())
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
