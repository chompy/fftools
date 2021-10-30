package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

const webListenPort = 31594

func initWeb() {
	http.HandleFunc("/favicon.ico", webServeFavicon)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pathSplit := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")
		scriptName := pathSplit[0]
		if scriptName == "" {
			webServeB64(webNotFound, http.StatusNotFound, w)
			return
		}
		var luaScript *luaScript = nil
		for _, _luaScript := range loadedScripts {
			if _luaScript.ScriptName == scriptName {
				luaScript = _luaScript
				break
			}
		}
		// lua script not found
		if luaScript == nil {
			webServeB64(webNotFound, http.StatusNotFound, w)
			return
		}
		// execute lua script end point
		if len(pathSplit) > 1 && pathSplit[1] == "_data" {
			webServeLua(luaScript.L, w, r)
			return
		}
		// serve up static file
		pathTo := filepath.Join(getScriptWebPath(scriptName), strings.Trim(strings.Join(pathSplit[1:], "/"), "/"))
		// TODO custom not found
		http.ServeFile(w, r, pathTo)
	})
	if err := http.ListenAndServe(fmt.Sprintf(":%d", webListenPort), nil); err != nil {
		logWarn(err.Error())
	}
}

func webServeB64(data string, status int, w http.ResponseWriter) {
	dataBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		logWarn(err.Error())
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
		logWarn(err.Error())
		webHandleError(w)
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "image/ico")
	w.WriteHeader(http.StatusOK)
	w.Write(dataBytes)
}

func webServeLua(L *lua.LState, w http.ResponseWriter, r *http.Request) {
	L.SetTop(0)
	luaRequest := &lua.LTable{}
	luaRequest.RawSetString("url", lua.LString(r.URL.String()))
	luaRequest.RawSetString("host", lua.LString(r.URL.Host))
	luaRequest.RawSetString("hostname", lua.LString(r.URL.Hostname()))
	luaRequest.RawSetString("port", lua.LString(r.URL.Port()))
	luaRequest.RawSetString("path", lua.LString(r.URL.Path))
	queryTable := &lua.LTable{}
	for k, v := range r.URL.Query() {
		queryTable.RawSetString(k, valueGoToLua(v))
	}
	luaRequest.RawSetString("query", queryTable)
	L.Push(L.GetGlobal("web"))
	L.Push(luaRequest)
	if err := L.PCall(1, 1, nil); err != nil {
		logLuaWarn(L, err.Error())
		webServeB64(webNotFound, http.StatusNotFound, w)
		return
	}
	res := valueLuaToGo(L.Get(1))
	resJson, err := json.Marshal(res)
	if err != nil {
		logLuaWarn(L, err.Error())
		webServeB64(webNotFound, http.StatusNotFound, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}

func webHandleError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("An unknown error occured."))
}
