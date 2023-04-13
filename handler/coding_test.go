package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInt2(t *testing.T) {
	test1 := make([]byte, 0, 2)
	test1 = append(test1, 1, 0)
	res, pos := GetInt2(test1, 0)
	assert.Equal(t, uint16(1), res)
	assert.Equal(t, int32(2), pos)
	test2 := make([]byte, 0, 2)
	test2 = append(test2, 255, 255)
	res, pos = GetInt2(test2, 0)
	assert.Equal(t, uint16(65535), res)
	assert.Equal(t, int32(2), pos)
}

func TestGetInt4(t *testing.T) {
	test1 := make([]byte, 0, 4)
	test1 = append(test1, 1, 0, 0, 0)
	res, pos := GetInt4(test1, 0)
	assert.Equal(t, uint32(1), res)
	assert.Equal(t, int32(4), pos)
	test2 := make([]byte, 0, 4)
	test2 = append(test2, 255, 255, 255, 255)
	res, pos = GetInt4(test2, 0)
	assert.Equal(t, uint32(4294967295), res)
	assert.Equal(t, int32(4), pos)
}

func TestGetLengthString(t *testing.T) {
	// test0 test empty string
	test0 := make([]byte, 0, 1)
	test0 = append(test0, 0x00)
	res, pos := GetLengthString(test0, 0)
	t.Log(res)
	assert.Equal(t, "", res)
	assert.Equal(t, int32(1), pos)
	// test1 This is a test
	test1 := make([]byte, 0, 16)
	test1 = append(test1, 0x0e, 'T', 'h', 'i', 's', ' ', 'i', 's', ' ', 'a', ' ', 't', 'e', 's', 't')
	res, pos = GetLengthString(test1, 0)
	t.Log(res)
	assert.Equal(t, "This is a test", res)
	assert.Equal(t, int32(15), pos)

	// test2 MySQL is too hard to learn.
	test2 := make([]byte, 0, 32)
	test2 = append(test2, 0x1b, 'M', 'y', 'S', 'Q', 'L', ' ', 'i', 's', ' ', 't', 'o', 'o', ' ', 'h', 'a', 'r', 'd', ' ', 't', 'o', ' ', 'l', 'e', 'a', 'r', 'n', '.')
	res, pos = GetLengthString(test2, 0)
	t.Log(res)
	assert.Equal(t, "MySQL is too hard to learn.", res)
	assert.Equal(t, int32(28), pos)

	// test3 repeat t * 255
	test3 := make([]byte, 0, 300)
	test3 = append(test3, 0xfc, 0xff, 0x00)
	for i := 0; i < 255; i++ {
		test3 = append(test3, 't')
	}
	res, pos = GetLengthString(test3, 0)
	t.Log(res)
	var expect string
	for i := 0; i < 255; i++ {
		expect += "t"
	}
	assert.Equal(t, expect, res)
	assert.Equal(t, int32(258), pos)
}
