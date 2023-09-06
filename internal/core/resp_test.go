package core_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"memkv/internal/core"
	"testing"
)

func TestSimpleStringDecode(t *testing.T) {
	cases := map[string]string{
		"+OK\r\n": "OK",
	}
	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		if v != value {
			t.Fail()
		}
	}
}

func TestError(t *testing.T) {
	cases := map[string]string{
		"-Error message\r\n": "Error message",
	}
	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		if v != value {
			t.Fail()
		}
	}
}

func TestInt64(t *testing.T) {
	cases := map[string]int64{
		":0\r\n":    0,
		":1000\r\n": 1000,
	}
	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		if v != value {
			t.Fail()
		}
	}
}

func TestBulkStringDecode(t *testing.T) {
	cases := map[string]string{
		"$5\r\nhello\r\n": "hello",
		"$0\r\n\r\n":      "",
	}
	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		if v != value {
			t.Fail()
		}
	}
}

func TestArrayDecode(t *testing.T) {
	cases := map[string][]interface{}{
		"*0\r\n":                                                   {},
		"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n":                     {"hello", "world"},
		"*3\r\n:1\r\n:2\r\n:3\r\n":                                 {int64(1), int64(2), int64(3)},
		"*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n":            {int64(1), int64(2), int64(3), int64(4), "hello"},
		"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n": {[]int64{int64(1), int64(2), int64(3)}, []interface{}{"Hello", "World"}},
	}
	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		array := value.([]interface{})
		if len(array) != len(v) {
			t.Fail()
		}
		for i := range array {
			if fmt.Sprintf("%v", v[i]) != fmt.Sprintf("%v", array[i]) {
				t.Fail()
			}
		}
	}
}

func TestEncodeString2DArray(t *testing.T) {
	var decode = [][]string{{"hello", "world"}, {"1", "2", "3"}, {"xyz"}}
	encode := core.Encode(decode, false)
	assert.EqualValues(t, "*3\r\n*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n*3\r\n$1\r\n1\r\n$1\r\n2\r\n$1\r\n3\r\n*1\r\n$3\r\nxyz\r\n", string(encode))
	decodeAgain, _ := core.Decode(encode)
	for i := 0; i < 3; i++ {
		for j := 0; j < len(decode[i]); j++ {
			assert.EqualValues(t, decode[i][j], decodeAgain.([]interface{})[i].([]interface{})[j])
		}
	}
}

func TestEncodeInterfaceArray(t *testing.T) {
	cases := map[string][]interface{}{
		"*0\r\n":                                        {},
		"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n":          {"hello", "world"},
		"*3\r\n:1\r\n:2\r\n:3\r\n":                      {int64(1), int64(2), int64(3)},
		"*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n": {int64(1), int64(2), int64(3), int64(4), "hello"},
	git 	"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n$5\r\nHello\r\n$5\r\nWorld\r\n": {[]interface{}{int64(1), int64(2), int64(3)}, []interface{}{"Hello", "World"}},
	}
	for k, v := range cases {
		encode := core.Encode(v, false)
		assert.EqualValues(t, k, string(encode))
	}
}

func TestParseCmd(t *testing.T) {
	cases := map[string]core.MemKVCmd{
		"*3\r\n$3\r\nput\r\n$5\r\nhello\r\n$5\r\nworld\r\n": core.MemKVCmd{
			Cmd:  "PUT",
			Args: []string{"hello", "world"},
		}}
	for k, v := range cases {
		cmd, _ := core.ParseCmd([]byte(k))
		if cmd.Cmd != v.Cmd {
			t.Fail()
		}
		if len(cmd.Args) != len(v.Args) {
			t.Fail()
		}
		for i := 0; i < len(cmd.Args); i++ {
			if cmd.Args[i] != v.Args[i] {
				t.Fail()
			}
		}
	}
}
