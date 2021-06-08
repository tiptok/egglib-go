package cache

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mohae/deepcopy"
	"reflect"
	"strconv"
)

func GetKey(args ...interface{}) string {
	addBuf := func(i interface{}, bf *bytes.Buffer) {
		switch v := i.(type) {
		case int:
			bf.WriteString(strconv.Itoa(v))
		case int8:
			bf.WriteString(strconv.Itoa(int(v)))
		case int16:
			bf.WriteString(strconv.Itoa(int(v)))
		case int32:
			bf.WriteString(strconv.Itoa(int(v)))
		case int64:
			bf.WriteString(strconv.Itoa(int(v)))
		case uint8:
			bf.WriteString(strconv.Itoa(int(v)))
		case uint16:
			bf.WriteString(strconv.Itoa(int(v)))
		case uint32:
			bf.WriteString(strconv.Itoa(int(v)))
		case uint64:
			bf.WriteString(strconv.Itoa(int(v)))
		case string:
			bf.WriteString(v)
		}
	}

	var buf bytes.Buffer
	for i, k := range args {
		addBuf(k, &buf)
		if i < len(args)-1 {
			addBuf(":", &buf)
		}
	}
	return buf.String()
}

// clone object to return, to avoid dirty data
func Clone(src, dst interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
			return
		}
	}()

	v := deepcopy.Copy(src)
	if reflect.ValueOf(v).IsValid() {
		reflect.ValueOf(dst).Elem().Set(reflect.Indirect(reflect.ValueOf(v)))
	}
	return
}
