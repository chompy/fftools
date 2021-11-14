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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func initWeb() {
	http.HandleFunc("/favicon.ico", webServeFavicon)
	http.HandleFunc("/", webHandle)
	config := configAppLoad()
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.PortWeb), nil); err != nil {
		logWarn(err.Error())
	}
}

func webHandle(w http.ResponseWriter, r *http.Request) {
	pathSplit := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")
	scriptName := pathSplit[0]
	if scriptName == "" {
		webServeB64(assetError404General, http.StatusNotFound, w)
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
		webServeB64(assetError404General, http.StatusNotFound, w)
		return
	}
	// execute lua script end point
	if len(pathSplit) > 1 && pathSplit[1] == "_data" {
		luaScript.Lock.Lock()
		defer luaScript.Lock.Unlock()
		if luaScript.State != LuaScriptActive {
			webServeB64(assetError404General, http.StatusInternalServerError, w)
			return
		}
		webServeLua(luaScript, w, r)
		return
	}
	// serve up static file
	pathTo := filepath.Join(getScriptWebPath(scriptName), strings.Trim(strings.Join(pathSplit[1:], "/"), "/"))
	if _, err := os.Stat(pathTo); err != nil {
		if os.IsNotExist(err) {
			webServeB64(assetError404General, http.StatusNotFound, w)
			return
		}
		webServeB64(assetError500General, http.StatusInternalServerError, w)
		return
	}
	// TODO custom not found
	http.ServeFile(w, r, pathTo)
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
	dataBytes, err := base64.StdEncoding.DecodeString(assetFavicon)
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

func webServeLua(ls *luaScript, w http.ResponseWriter, r *http.Request) {
	res, err := ls.Web(r)
	if err != nil {
		webServeB64(assetError404General, http.StatusInternalServerError, w)
		return
	}
	if res == nil {
		webServeB64(assetError404General, http.StatusNotFound, w)
		return
	}
	resJson, err := json.Marshal(res)
	if err != nil {
		webServeB64(assetError404General, http.StatusNotFound, w)
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
