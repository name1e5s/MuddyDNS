// Package utils provides some utility functions for implementing
// a simple DNS server.
//
// In this package we provide the functions for converting binary
// formats and identifying valid IP/domain name strings.
package utils

import "encoding/binary"

func BytesToUInt16(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
}

func BytesToUInt32(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}

func UInt16ToBytes(data uint16) (byte, byte) {
	buffer := make([]byte, 2)
	binary.BigEndian.PutUint16(buffer, data)
	return buffer[0], buffer[1]
}

func BoolToInt(Bool bool) int {
	if Bool {
		return 1
	}
	return 0
}
