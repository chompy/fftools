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
	"encoding/binary"
	"time"
)

func writeUint16(data *[]byte, value uint16) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, value)
	*data = append(*data, buf...)
}

func writeInt32(data *[]byte, value int32) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(value))
	*data = append(*data, buf...)
}

func writeByte(data *[]byte, value byte) {
	*data = append(*data, value)
}

func writeBool(data *[]byte, value bool) {
	if value {
		*data = append(*data, 1)
		return
	}
	*data = append(*data, 0)
}

func writeString(data *[]byte, value string) {
	writeUint16(data, uint16(len(value)))
	*data = append(*data, []byte(value)...)
}

func writeTime(data *[]byte, value time.Time) {
	writeString(data, value.UTC().Format(time.RFC3339Nano))
}
