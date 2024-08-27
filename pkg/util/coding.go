package util

import "encoding/binary"

// GetInt2 MySQL always uses little-endian byte order for integer values.
func GetInt2(data []byte, pos int32) (uint16, int32) {
	return binary.LittleEndian.Uint16(data[pos:]), pos + 2
}

// SetInt2 MySQL always uses little-endian byte order for integer values.
func SetInt2(data []byte, pos int32, val uint16) {
	binary.LittleEndian.PutUint16(data[pos:], val)
}

// GetInt4 MySQL always uses little-endian byte order for integer values.
func GetInt4(data []byte, pos int32) (uint32, int32) {
	return binary.LittleEndian.Uint32(data[pos:]), pos + 4
}

// SetInt4 MySQL always uses little-endian byte order for integer values.
func SetInt4(data []byte, pos int32, val uint32) {
	binary.LittleEndian.PutUint32(data[pos:], val)
}

// GetLengthString max length of string is 65535
// If the first byte is < 0xfc (252), use one byte to represent the length of the string.
// If the first byte is 0xfc (252), use the next two bytes to represent the length of the string.
func GetLengthString(data []byte, pos int32) (string, int32) {
	if data[pos] < 252 {
		strLen := data[pos]
		return string(data[pos+1 : pos+1+int32(strLen)]), pos + 1 + int32(strLen)
	}
	// assert data[pos] == 252
	strLen, pos := GetInt2(data, pos+1)
	return string(data[pos : pos+int32(strLen)]), pos + int32(strLen)
}
