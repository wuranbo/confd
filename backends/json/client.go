// Package main provides keys from jsonfiles,
// the format should follow:
// {
//   "prefix": "/myapp",
//   "z.sh": [
//      {
//        "key": "heapsize",
//        "value": "152m"
//      },
//    ],
//   "others": [
//      {
//        "fullkey": "/myapp/modules/frontend/env/LOG_DIR",
//        "value": "/var/log/myapp"
//      }
//    ]
// }

package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/wuranbo/confd/log"
)

// Client provides a wrapper around the json client
type Client struct {
	files []string
	kvs   map[string]string
}

// NewEnvClient returns a new client
func NewJsonClient(files []string) (*Client, error) {
	kvs := make(map[string]string, 0)
	if len(files) == 0 {
		return &Client{files, kvs}, errors.New("Please input the jsonfile in option -nodes.")
	}

	for _, f := range files {
		for k, v := range readKVsFromFile(f) {
			kvs[k] = v // later file override early files
		}
	}

	return &Client{files, kvs}, nil
}

type pair struct {
	FullKey string `json:"fullkey,omitempty"`
	Key     string `json:"key,omitempty"`
	Value   string `json:"value"`
}

// return not-empty map, or nil if error.
func readKVsFromFile(filename string) map[string]string {
	jsonstr, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("read jsonfile " + filename + " failed.")
		return nil
	}

	var prefix string = "/"
	prefixdec := json.NewDecoder(bytes.NewReader(jsonstr))
	for {
		var prefixmap map[string]interface{}
		err := prefixdec.Decode(&prefixmap)
		for k, v := range prefixmap {
			if k == "prefix" {
				prefix = filepath.Join("/", v.(string))
				log.Notice("file: " + filename + " decode prefix == " + prefix)
				delete(prefixmap, k)

				b := new(bytes.Buffer)
				enc := json.NewEncoder(b)
				enc.Encode(&prefixmap)
				jsonstr = b.Bytes()
				break
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			continue
		}
	}
	if prefix == "/" {
		log.Notice("file " + filename + " not has prefix")
	}

	ret := make(map[string]string, 0)
	outerdec := json.NewDecoder(bytes.NewReader(jsonstr))
	for {
		var outer map[string][]interface{}
		if err := outerdec.Decode(&outer); err == io.EOF {
			log.Notice("read file " + filename + " over.")
			return ret
		} else if err != nil {
			log.Error("read file " + filename + " error:" + err.Error())
			return nil
		}

		for _, v := range outer {
			for _, ele := range v {
				b := new(bytes.Buffer)
				enc := json.NewEncoder(b)
				if err := enc.Encode(&ele); err != nil {
					log.Error("decode file " + filename + " failed, around " + fmt.Sprintf("%v", ele) + " with error:" + err.Error())
					return nil
				}
				s := b.String()

				dec := json.NewDecoder(strings.NewReader(s))
				var p pair
				err := dec.Decode(&p)
				if p.Key != "" && p.Value != "" {
					k := filepath.Join(prefix, p.Key)
					ret[k] = p.Value
					log.Notice(fmt.Sprintf("read json key:%s, value:%s\n", k, p.Value))
					continue
				} else if p.FullKey != "" && p.Value != "" {
					k := filepath.Join("/", p.FullKey)
					ret[k] = p.Value
					log.Notice(fmt.Sprintf("read json key:%s, value:%s\n", k, p.Value))
					continue
				} else if err == io.EOF {
					break
				} else if err != nil {
					log.Error("decode file " + filename + " failed, around " + fmt.Sprintf("%v", ele) + " with error:" + err.Error())
					return nil
				}
			}
		}
	}
}

// GetValues queries the environment for keys
func (c *Client) GetValues(keys []string) (map[string]string, error) {
	vars := make(map[string]string)
	for _, key := range keys {
		for k, v := range c.kvs {
			if strings.HasPrefix(k, key) {
				vars[k] = v
			}
		}
	}
	return vars, nil
}

func (c *Client) WatchPrefix(prefix string, waitIndex uint64, stopChan chan bool) (uint64, error) {
	<-stopChan
	return 0, nil
}
