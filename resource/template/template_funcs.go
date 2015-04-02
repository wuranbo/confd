package template

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func newFuncMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["base"] = path.Base
	m["split"] = strings.Split
	m["json"] = UnmarshalJsonObject
	m["jsonArray"] = UnmarshalJsonArray
	m["dir"] = path.Dir
	m["getenv"] = os.Getenv
	m["join"] = strings.Join
	m["datetime"] = time.Now
	m["concat"] = Concat
	m["byteToM"] = ByteToM
	m["strsub"] = StringSub
	m["stradd"] = StringAdd
	m["strmul"] = StringMul
	m["strdiv"] = StringDiv
	return m
}

func addFuncs(out, in map[string]interface{}) {
	for name, fn := range in {
		out[name] = fn
	}
}

func UnmarshalJsonObject(data string) (map[string]interface{}, error) {
	var ret map[string]interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}

func UnmarshalJsonArray(data string) ([]interface{}, error) {
	var ret []interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}

func Concat(strs ...interface{}) string {
	return fmt.Sprint(strs...)
}

func ByteToM(data string) (string, error) {
	r, err := strconv.ParseUint(data, 10, 64)
	if err != nil || r < 0 {
		return "0m", err
	}
	if r < 1024*1024*1024 {
		return "1m", err
	} else {
		return fmt.Sprint(r/(1024*1024*1024), "m"), err
	}
}

func StringAdd(num1 int, num2 int) int {
	return num1 + num2
}

func StringSub(num1 int, num2 int) int {
	var rlt = num1 - num2
	if rlt < 0 {
		rlt = 0
	}
	return rlt
}

func StringMul(num1 int, num2 int) int {
	return num1 * num2
}

func StringDiv(num1 int, num2 int) int {
	var rlt = num1 / num2
	if rlt <= 0 {
		rlt = 1
	}
	return rlt
}
