package main

import (
	"reflect"
	"testing"

	"github.com/wuranbo/confd/log"
)

func TestInitConfigDefaultConfig(t *testing.T) {
	log.SetQuiet(true)
	want := Config{
		Backend:      "etcd",
		BackendNodes: []string{"http://127.0.0.1:4001"},
		ClientCaKeys: "",
		ClientCert:   "",
		ClientKey:    "",
		ConfDir:      "/etc/confd",
		Debug:        false,
		Interval:     600,
		Noop:         false,
		Prefix:       "/",
		Quiet:        false,
		SRVDomain:    "",
		Scheme:       "http",
		Verbose:      false,
	}
	if err := initConfig(); err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(want, config) {
		t.Errorf("initConfig() = %v, want %v", config, want)
	}
}
