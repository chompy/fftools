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
	"encoding/json"
	"io/ioutil"
	"os"
)

const persistUsersFile = "./data/_proxy_users.json"

func persistUsers() error {
	proxyUsersJSON, err := json.Marshal(proxyUsers)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(persistUsersFile, proxyUsersJSON, 0600)
}

func readUsers() error {
	proxyUsersJSON, err := ioutil.ReadFile(persistUsersFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(proxyUsersJSON, &proxyUsers)
}
