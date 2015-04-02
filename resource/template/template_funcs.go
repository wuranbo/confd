package template

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
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

func toInt(data interface{}) (ret int64, err error) {
	switch v := data.(type) {
	case float32:
		ret = int64(v)
	case float64:
		ret = int64(v)
	case int:
		ret = int64(v)
	case int8:
		ret = int64(v)
	case int16:
		ret = int64(v)
	case int32:
		ret = int64(v)
	case uint:
		ret = int64(v)
	case uint8:
		ret = int64(v)
	case uint16:
		ret = int64(v)
	case uint32:
		ret = int64(v)
	case uint64:
		ret = int64(v)
	case string:
		ret, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			tmpfloat, terr := strconv.ParseFloat(v, 10)
			ret = int64(tmpfloat)
			err = terr // haha! go compile"s bug
			if err != nil {
				tmpuint, terr := strconv.ParseUint(v, 10, 63)
				ret = int64(tmpuint)
				err = terr
			}
		}
	case []byte:
		ret, err = strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			tmpfloat, terr := strconv.ParseFloat(string(v), 10)
			ret = int64(tmpfloat)
			err = terr
			if err != nil {
				tmpuint, terr := strconv.ParseUint(string(v), 10, 63)
				ret = int64(tmpuint)
				err = terr
			}
		}
	default:
		ret = 0
	}
	return ret, err
}

func StringAdd(num1, num2 interface{}) (ret string, reterr error) {
	ret = "0"
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			reterr = r.(error)
		}
	}()
	inta, reterr := toInt(num1)
	if reterr != nil {
		return
	}
	intb, reterr := toInt(num2)
	if reterr != nil {
		return
	}
	ret = strconv.FormatInt(inta+intb, 10)
	return
}

func StringSub(num1, num2 interface{}) (ret string, reterr error) {
	ret = "0"
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			reterr = r.(error)
		}
	}()
	inta, reterr := toInt(num1)
	if reterr != nil {
		return
	}
	intb, reterr := toInt(num2)
	if reterr != nil {
		return
	}
	ret = strconv.FormatInt(inta-intb, 10)
	return
}

func StringMul(num1, num2 interface{}) (ret string, reterr error) {
	ret = "0"
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			reterr = r.(error)
		}
	}()
	inta, reterr := toInt(num1)
	if reterr != nil {
		return
	}
	intb, reterr := toInt(num2)
	if reterr != nil {
		return
	}
	ret = strconv.FormatInt(inta*intb, 10)
	return
}

func StringDiv(num1, num2 interface{}) (ret string, reterr error) {
	ret = "0"
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			reterr = r.(error)
		}
	}()
	inta, reterr := toInt(num1)
	if reterr != nil {
		return
	}
	intb, reterr := toInt(num2)
	if reterr != nil {
		return
	}
	ret = strconv.FormatInt(inta/intb, 10)
	return
}
