package inmem

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/kelseyhightower/memkv"
	"github.com/wuranbo/confd/backends"
	confdtmpl "github.com/wuranbo/confd/resource/template"
)

type InmemConfig struct {
	Prefix      string // all template and toml has same  prefix
	StoreClient backends.StoreClient
}

type TomlContainer struct {
	TomlTemplateSection `toml:"template"`
}

type TomlTemplateSection struct {
	Src  string
	Dest string
	Keys []string
	Mode string
}

// InmemTemplateResource is the representation of a parsed template resource.
type InmemTemplateResource struct {
	Origin      string
	Stage       TextResource      // tmp save position, should be assigned to Dest
	Dest        InmemTemplateDest // result in memeory
	Src         InmemTemplateSrc  // template file in memory
	Keys        []string
	Prefix      string
	funcMap     map[string]interface{}
	lastIndex   uint64
	prefix      string
	store       memkv.Store
	storeClient backends.StoreClient
}

func (tr *InmemTemplateResource) Name() string {
	return tr.Origin
}

type TextResource struct {
	Data []byte
}

func (d *TextResource) Write(p []byte) (n int, err error) {
	d.Data = append(d.Data, p...)
	return len(p), nil
}
func (d *TextResource) String() string {
	return string(d.Data)
}

type InmemTemplateDest struct {
	Origin string
	Data   TextResource
}

func (td *InmemTemplateDest) Name() string {
	return td.Origin
}

type InmemTemplateSrc InmemTemplateDest

func (ts *InmemTemplateSrc) Name() string {
	return ts.Origin
}

// NewInmemTemplateResource creates a TemplateResource.
func NewInmemTemplateResource(
	tomltext string,
	tmpltext string,
	config InmemConfig) (*InmemTemplateResource, error) {

	if config.StoreClient == nil {
		return nil, errors.New("A valid StoreClient is required.")
	}
	var tc TomlContainer
	_, err := toml.Decode(tomltext, &tc)
	if err != nil {
		return nil, fmt.Errorf("Cannot process template resource, error:%s", err.Error())
	}
	data := make([]byte, len(tmpltext))
	copy(data[:], tmpltext)
	tr := InmemTemplateResource{
		Keys: tc.TomlTemplateSection.Keys,
		Dest: InmemTemplateDest{Origin: tc.TomlTemplateSection.Dest},
		Src:  InmemTemplateSrc{Origin: tc.TomlTemplateSection.Src, Data: TextResource{data}},
	}

	tr.storeClient = config.StoreClient
	tr.funcMap = confdtmpl.NewFuncMap()
	tr.store = memkv.New()
	confdtmpl.AddFuncs(tr.funcMap, tr.store.FuncMap)
	tr.prefix = filepath.Join("/", config.Prefix, tr.Prefix)
	if tr.Src.Origin == "" {
		return nil, ErrEmptySrc
	}
	return &tr, nil
}

// setVars sets the Vars for template resource.
func (t *InmemTemplateResource) setVars() error {
	var err error
	result, err := t.storeClient.GetValues(appendPrefix(t.prefix, t.Keys))
	if err != nil {
		return err
	}
	t.store.Purge()
	for k, v := range result {
		t.store.Set(filepath.Join("/", strings.TrimPrefix(k, t.prefix)), v)
	}
	return nil
}

// createStage stages the src configuration file by processing the src
// template and setting the desired owner, group, and mode. It also sets the
// StageFile for the template resource.
// It returns an error if any.
func (t *InmemTemplateResource) createStage() error {
	temp := TextResource{[]byte{}}
	tmpl := template.Must(template.New(t.Src.Name()).Funcs(t.funcMap).Parse(t.Src.Data.String()))
	if err := tmpl.Execute(&temp, nil); err != nil {
		return err
	}
	t.Stage = temp
	return nil
}

func (t *InmemTemplateResource) sync() error {
	t.Dest.Data = t.Stage
	return nil
}

// process is a convenience function that wraps calls to the three main tasks
// required to keep local configuration files in sync. First we gather vars
// from the store, then we stage a candidate configuration file, and finally sync
// things up.
// It returns an error if any.
func (t *InmemTemplateResource) process() error {
	if err := t.setVars(); err != nil {
		return err
	}
	if err := t.createStage(); err != nil {
		return err
	}
	t.Dest.Data = t.Stage
	return nil
}
func (t *InmemTemplateResource) Process() error {
	return t.process()
}

var ErrEmptySrc = errors.New("empty src template")

func appendPrefix(prefix string, keys []string) []string {
	s := make([]string, len(keys))
	for i, k := range keys {
		s[i] = path.Join(prefix, k)
	}
	return s
}
