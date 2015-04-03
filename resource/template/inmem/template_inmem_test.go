package template

import (
	"fmt"
	"testing"

	"github.com/wuranbo/confd/backends"
)

type inmemTemplateTest struct {
	desc        string                       // description of the test (for helpful errors)
	toml        string                       // toml file contents
	tmpl        string                       // template file contents
	expected    string                       // expected generated file contents
	updateStore func(*InmemTemplateResource) // function for setting values in store
}

// inmemTemplateTests is an array of inmemTemplateTest structs, each representing a test of
// some aspect of template processing. When the input tmpl and toml files are
// processed, they should produce a config file matching expected.
var inmemTemplateTests = []inmemTemplateTest{

	inmemTemplateTest{
		desc: "base, get test",
		toml: `
[template]
src = "test.conf.tmpl"
dest = "./tmp/test.conf"
keys = [
    "/test/key",
]
`,
		tmpl: `
{{with get "/test/key"}}
key: {{base .Key}}
val: {{.Value}}
{{end}}
`,
		expected: `

key: key
val: abc

`,
		updateStore: func(tr *InmemTemplateResource) {
			fmt.Println("pre store.set:", "dong")
			tr.store.Set("/test/key", "abc")
			fmt.Println("after store.set:", "dong")
		},
	},

	inmemTemplateTest{
		desc: "gets test",
		toml: `
[template]
src = "test.conf.tmpl"
dest = "./tmp/test.conf"
keys = [
    "/test/user",
    "/test/pass",
    "/nada/url",
]
`,
		tmpl: `
{{range gets "/test/*"}}
key: {{.Key}}
val: {{.Value}}
{{end}}
`,
		expected: `

key: /test/pass
val: abc

key: /test/user
val: mary

`,
		updateStore: func(tr *InmemTemplateResource) {
			tr.store.Set("/test/user", "mary")
			tr.store.Set("/test/pass", "abc")
			tr.store.Set("/nada/url", "url")
		},
	},

	inmemTemplateTest{
		desc: "getv test",
		toml: `
[template]
src = "test.conf.tmpl"
dest = "./tmp/test.conf"
keys = [
    "/test/url",
    "/test/user",
]
`,
		tmpl: `
url = {{getv "/test/url"}}
user = {{getv "/test/user"}}
`,
		expected: `
url = http://www.abc.com
user = bob
`,
		updateStore: func(tr *InmemTemplateResource) {
			tr.store.Set("/test/url", "http://www.abc.com")
			tr.store.Set("/test/user", "bob")
		},
	},

	inmemTemplateTest{
		desc: "getvs test",
		toml: `
[template]
src = "test.conf.tmpl"
dest = "./tmp/test.conf"
keys = [
    "/test/user",
    "/test/pass",
    "/nada/url",
]
`,
		tmpl: `
{{range getvs "/test/*"}}
val: {{.}}
{{end}}
`,
		expected: `

val: abc

val: mary

`,
		updateStore: func(tr *InmemTemplateResource) {
			tr.store.Set("/test/user", "mary")
			tr.store.Set("/test/pass", "abc")
			tr.store.Set("/nada/url", "url")
		},
	},

	inmemTemplateTest{
		desc: "split test",
		toml: `
[template]
src = "test.conf.tmpl"
dest = "./tmp/test.conf"
keys = [
    "/test/data",
]
`,
		tmpl: `
{{$data := split (getv "/test/data") ":"}}
f: {{index $data 0}}
br: {{index $data 1}}
bz: {{index $data 2}}
`,
		expected: `

f: foo
br: bar
bz: baz
`,
		updateStore: func(tr *InmemTemplateResource) {
			tr.store.Set("/test/data", "foo:bar:baz")
		},
	},

	inmemTemplateTest{
		desc: "json test",
		toml: `
[template]
src = "test.conf.tmpl"
dest = "./tmp/test.conf"
keys = [
    "/test/data/",
]
`,
		tmpl: `
{{range getvs "/test/data/*"}}
{{$data := json .}}
id: {{$data.Id}}
ip: {{$data.IP}}
{{end}}
`,
		expected: `


id: host1
ip: 192.168.10.11


id: host2
ip: 192.168.10.12

`,
		updateStore: func(tr *InmemTemplateResource) {
			tr.store.Set("/test/data/1", `{"Id":"host1", "IP":"192.168.10.11"}`)
			tr.store.Set("/test/data/2", `{"Id":"host2", "IP":"192.168.10.12"}`)
		},
	},

	inmemTemplateTest{
		desc: "jsonArray test",
		toml: `
[template]
src = "test.conf.tmpl"
dest = "./tmp/test.conf"
keys = [
    "/test/data/",
]
`,
		tmpl: `
{{range jsonArray (getv "/test/data/")}}
num: {{.}}
{{end}}
`,
		expected: `

num: 1

num: 2

num: 3

`,
		updateStore: func(tr *InmemTemplateResource) {
			tr.store.Set("/test/data/", `["1", "2", "3"]`)
		},
	},

	inmemTemplateTest{
		desc: "ls test",
		toml: `
[template]
src = "test.conf.tmpl"
dest = "./tmp/test.conf"
keys = [
    "/test/data/abc",
    "/test/data/def",
    "/test/data/ghi",
]
`,
		tmpl: `
{{range ls "/test/data"}}
value: {{.}}
{{end}}
`,
		expected: `

value: abc

value: def

value: ghi

`,
		updateStore: func(tr *InmemTemplateResource) {
			tr.store.Set("/test/data/abc", "123")
			tr.store.Set("/test/data/def", "456")
			tr.store.Set("/test/data/ghi", "789")
		},
	},

	inmemTemplateTest{
		desc: "lsdir test",
		toml: `
[template]
src = "test.conf.tmpl"
dest = "./tmp/test.conf"
keys = [
    "/test/data/abc",
    "/test/data/def/ghi",
    "/test/data/jkl/mno",
]
`,
		tmpl: `
{{range lsdir "/test/data"}}
value: {{.}}
{{end}}
`,
		expected: `

value: def

value: jkl

`,
		updateStore: func(tr *InmemTemplateResource) {
			tr.store.Set("/test/data/abc", "123")
			tr.store.Set("/test/data/def/ghi", "456")
			tr.store.Set("/test/data/jkl/mno", "789")
		},
	},
	inmemTemplateTest{
		desc: "dir test",
		toml: `
[template]
src = "test.conf.tmpl"
dest = "./tmp/test.conf"
keys = [
    "/test/data",
    "/test/data/abc",
]
`,
		tmpl: `
{{with dir "/test/data/abc"}}
dir: {{.}}
{{end}}
`,
		expected: `

dir: /test/data

`,
		updateStore: func(tr *InmemTemplateResource) {
			tr.store.Set("/test/data", "parent")
			tr.store.Set("/test/data/def", "child")
		},
	},
}

// TestInmemTemplates runs all tests in inmemTemplateTests
func TestInmemTemplates(t *testing.T) {
	for _, tt := range inmemTemplateTests {
		ExecuteTestInmemTemplate(tt, t)
	}
}

// ExectureTestInmemTemplate builds a TemplateResource based on the toml and tmpl files described
// in the inmemTemplateTest, writes a config file, and compares the result against the expectation
// in the inmemTemplateTest.
func ExecuteTestInmemTemplate(tt inmemTemplateTest, t *testing.T) {
	tr, err := inmemTemplateResource(&tt)
	if err != nil {
		t.Errorf(tt.desc + ": failed to create InmemTemplateResource: " + err.Error())
	}

	tt.updateStore(tr)

	if err := tr.createStage(); err != nil {
		t.Errorf(tt.desc + ": failed createStage: " + err.Error())
	}

	actual := tr.Stage.String()
	if err != nil {
		t.Errorf(tt.desc + ": failed to read Stage: " + err.Error())
	}
	if actual != tt.expected {
		t.Errorf(fmt.Sprintf("%v: invalid Stage. Expected %v, actual %v", tt.desc, tt.expected, string(actual)))
	}
}

// inmemTemplateResource creates a inmemTemplateResource for creating a config file
func inmemTemplateResource(testdata *inmemTemplateTest) (*InmemTemplateResource, error) {
	backendConf := backends.Config{
		Backend: "env"}
	client, err := backends.New(backendConf)
	if err != nil {
		return nil, err
	}

	config := InmemConfig{
		StoreClient: client, // not used but must be set
	}

	fmt.Println("(((((((((((testdata:", testdata)
	tr, err := NewInmemTemplateResource(testdata.toml, testdata.tmpl, config)
	if err != nil {
		return nil, err
	}
	return tr, nil
}
