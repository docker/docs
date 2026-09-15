package main

import (
	js "github.com/santhosh-tekuri/jsonschema/v6"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func fixture(t *testing.T, name string) *Document {
	t.Helper()
	d, e := loadDocument(filepath.Join("testdata", name), "validation/dialects")
	if e != nil {
		t.Fatal(e)
	}
	d.Source = Source{ID: "test"}
	return d
}
func TestActualDialectAndReferences(t *testing.T) {
	d := fixture(t, "valid.yaml")
	d.validate("")
	for _, x := range d.Diagnostics {
		t.Errorf("%s: %s", x.Rule, x.Message)
	}
	if len(d.operations()) != 4 {
		t.Fatal("operation omitted")
	}
	var query Object
	for _, op := range d.operations() {
		if op["method"] == "QUERY" {
			query = op
		}
	}
	if query == nil {
		t.Fatal("QUERY omitted")
	}
	if obj(arr(query["servers"])[0])["url"] != "https://search.example.test" {
		t.Fatal("path server override lost")
	}
	p := obj(arr(query["parameters"])[0])
	if p["example"] != jsonNumber("0") {
		t.Fatalf("parameter override/zero lost: %#v", p["example"])
	}
	if len(arr(query["parameters"])) != 1 {
		t.Fatal("duplicate overridden parameter")
	}
}
func jsonNumber(s string) any {
	v, e := js.UnmarshalJSON(strings.NewReader(s))
	if e != nil {
		panic(e)
	}
	return v
}
func TestNegativeFixtures(t *testing.T) {
	for _, tc := range []struct{ name, rule string }{{"invalid-info.yaml", "structure"}, {"invalid-item-type.yaml", "schema"}, {"invalid-example.yaml", "example"}, {"profile-missing-id.yaml", "S4"}} {
		t.Run(tc.name, func(t *testing.T) {
			d := fixture(t, tc.name)
			d.validate("")
			found := false
			for _, x := range d.Diagnostics {
				found = found || x.Rule == tc.rule
			}
			if !found {
				t.Fatalf("missing %s diagnostic: %+v", tc.rule, d.Diagnostics)
			}
		})
	}
	for _, name := range []string{"duplicate-key.yaml", "invalid-stream-ref.yaml"} {
		t.Run(name, func(t *testing.T) {
			d, e := loadDocument(filepath.Join("testdata", name), "validation/dialects")
			if e == nil {
				d.validate("")
				for _, x := range d.Diagnostics {
					if x.Rule == "reference" || x.Rule == "schema" {
						return
					}
				}
				t.Fatal("invalid source accepted")
			}
		})
	}
}
func TestSchemaSemantics(t *testing.T) {
	d := fixture(t, "valid.yaml")
	c := d.Compiler
	cases := []struct {
		name           string
		schema         any
		valid, invalid any
	}{{"boolean false", false, nil, "x"}, {"ref siblings", Object{"$defs": Object{"base": Object{"type": "object", "required": []any{"a"}}}, "$ref": "#/$defs/base", "required": []any{"b"}}, Object{"a": true, "b": true}, Object{"b": true}}, {"recursive", Object{"type": "object", "properties": Object{"next": Object{"$ref": "#"}}}, Object{"next": Object{}}, Object{"next": 1}}, {"anchor", Object{"$defs": Object{"x": Object{"$anchor": "x", "type": "boolean"}}, "$ref": "#x"}, false, "false"}, {"dynamic", Object{"$dynamicAnchor": "node", "type": "object", "properties": Object{"child": Object{"$dynamicRef": "#node"}}}, Object{"child": Object{}}, Object{"child": false}}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			uri := "https://fixtures.test/" + strings.ReplaceAll(tc.name, " ", "-")
			s := tc.schema
			if m, ok := s.(map[string]any); ok {
				m["$schema"] = dialect
			}
			if e := c.AddResource(uri, s); e != nil {
				t.Fatal(e)
			}
			compiled, e := c.Compile(uri)
			if e != nil {
				t.Fatal(e)
			}
			if tc.valid != nil {
				if e = compiled.Validate(tc.valid); e != nil {
					t.Fatal(e)
				}
			}
			if e = compiled.Validate(tc.invalid); e == nil {
				t.Fatal("invalid example accepted")
			}
		})
	}
}
func TestAuthAndCurl(t *testing.T) {
	schemes := Object{"bearer": Object{"type": "http", "scheme": "bearer"}, "key": Object{"type": "apiKey", "in": "header", "name": "X-Key"}}
	op := Object{"method": "GET", "path": "/flag/{id}", "servers": []any{Object{"url": "https://api.test"}}, "securitySchemes": schemes, "security": []any{Object{"bearer": []any{}, "key": []any{}}, Object{}}, "parameters": []any{Object{"name": "id", "in": "path", "required": true, "example": "a/b"}, Object{"name": "enabled", "in": "query", "example": false}, Object{"name": "limit", "in": "query", "example": 0}}}
	curl, _ := curlExample(Source{}, op)
	for _, want := range []string{"a%2Fb", "enabled=false", "limit=0", "${TOKEN}", "X-Key: ${API_KEY}"} {
		if !strings.Contains(curl, want) {
			t.Errorf("missing %s: %s", want, curl)
		}
	}
	op["security"] = []any{}
	curl, _ = curlExample(Source{}, op)
	if strings.Contains(curl, "Authorization") {
		t.Fatal("anonymous operation gained auth")
	}
	op["method"] = "HEAD"
	curl, _ = curlExample(Source{Connection: "unix", ID: "engine-1.56"}, op)
	if !strings.Contains(curl, "--head") || !strings.Contains(curl, "--unix-socket") {
		t.Fatal(curl)
	}
	d := &Document{}
	root := Object{"security": []any{Object{"bearer": []any{}}}}
	if !reflect.DeepEqual(d.effective(root, Object{}, Object{"security": []any{}}, "security"), []any{}) {
		t.Fatal("empty override lost")
	}
}
func TestFalseSchemaAndExamplePreserved(t *testing.T) {
	b, e := os.ReadFile("testdata/valid.yaml")
	if e != nil {
		t.Fatal(e)
	}
	tmp := t.TempDir()
	if e = os.WriteFile(filepath.Join(tmp, "spec.yaml"), []byte(strings.ReplaceAll(string(b), "schema: true", "schema: false")), 0600); e != nil {
		t.Fatal(e)
	}
	if e = os.Mkdir(filepath.Join(tmp, "schemas"), 0700); e != nil {
		t.Fatal(e)
	}
	event, _ := os.ReadFile("testdata/schemas/event.yaml")
	os.WriteFile(filepath.Join(tmp, "schemas/event.yaml"), event, 0600)
	d, e := loadDocument(filepath.Join(tmp, "spec.yaml"), "validation/dialects")
	if e != nil {
		t.Fatal(e)
	}
	d.validate("")
	v, e := pointer(d.Root, "/paths/~1search/query/responses/200/content/application~1json/schema")
	if e != nil || v != false {
		t.Fatal("false schema changed")
	}
	for _, op := range d.operations() {
		if op["path"] == "/public" {
			ex := obj(arr(obj(arr(op["variants"])[0])["examples"])[0])
			if ex["value"] != false {
				t.Fatal("false example omitted")
			}
		}
	}
}

func TestMediaExamplesFollowReferenceAnnotations(t *testing.T) {
	d := &Document{URI: "https://example.test/api.yaml", Registry: &Registry{resources: map[string]any{
		"https://example.test/api.yaml": Object{
			"alias": Object{"$ref": "#/base", "examples": []any{false, jsonNumber("0")}},
			"base":  Object{"example": "base example"},
			"cycle": Object{"$ref": "#/cycle"},
		},
		"https://example.test/schemas/alias.yaml": Object{"$ref": "base.yaml"},
		"https://example.test/schemas/base.yaml":  Object{"example": "external example"},
	}}}
	for _, tc := range []struct {
		name   string
		media  Object
		values []any
	}{
		{"intermediate annotations", Object{"schema": Object{"$ref": "#/alias"}}, []any{false, jsonNumber("0")}},
		{"nearest schema annotation", Object{"schema": Object{"$ref": "#/alias", "example": "local"}}, []any{"local"}},
		{"media annotation", Object{"example": false, "schema": Object{"$ref": "#/alias"}}, []any{false}},
		{"relative external reference", Object{"schema": Object{"$ref": "schemas/alias.yaml"}}, []any{"external example"}},
		{"cycle", Object{"schema": Object{"$ref": "#/cycle"}}, []any{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			values := []any{}
			for _, ex := range d.mediaExamples(tc.media) {
				values = append(values, obj(ex)["value"])
			}
			if !reflect.DeepEqual(values, tc.values) {
				t.Fatalf("examples = %#v, want %#v", values, tc.values)
			}
		})
	}
}

func TestDVPResponseContracts(t *testing.T) {
	d, err := loadDocument(filepath.Join("..", "..", "content", "reference", "api", "dvp", "latest.yaml"), "validation/dialects")
	if err != nil {
		t.Fatal(err)
	}
	d.Source = Source{ID: "dvp"}
	d.validate("")
	if len(d.Diagnostics) != 0 {
		t.Fatalf("DVP must pass without exceptions: %+v", d.Diagnostics)
	}
	for _, tc := range []struct {
		schema string
		value  string
		valid  bool
	}{
		{"TimespanModel", `{"month":7}`, true},
		{"TimespanModel", `{"week":31}`, true},
		{"TimespanModel", `{}`, false},
		{"TimespanModel", `{"month":7,"week":31}`, false},
		{"TimespanModel", `{"month":{"month":7}}`, false},
		{"TimespanModel", `7`, false},
		{"TimespanData", `{"months":[{"month":5},{"month":7}]}`, true},
		{"TimespanData", `{"weeks":[]}`, true},
		{"TimespanData", `{"months":[],"weeks":[]}`, false},
		{"TimespanData", `{"months":[{}]}`, false},
		{"TimespanData", `{"weeks":null}`, false},
		{"NamespaceMetadata", `{"namespace":"org1","extraRepos":null,"extensionPublisher":false}`, true},
		{"PullData", `{"pulls":null}`, true},
		{"PullData", `{"pulls":[{"start":"2022-08-01T00:00:00Z","pullCount":0}]}`, true},
	} {
		t.Run(tc.schema+"/"+tc.value, func(t *testing.T) {
			s, err := d.Compiler.Compile(d.URI + "#/components/schemas/" + tc.schema)
			if err != nil {
				t.Fatal(err)
			}
			value, err := js.UnmarshalJSON(strings.NewReader(tc.value))
			if err != nil {
				t.Fatal(err)
			}
			if err := s.Validate(value); (err == nil) != tc.valid {
				t.Fatalf("valid = %t, want %t: %v", err == nil, tc.valid, err)
			}
		})
	}
}

func TestTimestampStringsRetainSpelling(t *testing.T) {
	p := filepath.Join(t.TempDir(), "schema.yaml")
	if e := os.WriteFile(p, []byte("type: string\nexample: 2021-01-05T21:06:53.506400Z\n"), 0600); e != nil {
		t.Fatal(e)
	}
	v, e := strictFragment(p)
	if e != nil {
		t.Fatal(e)
	}
	if obj(v)["example"] != "2021-01-05T21:06:53.506400Z" {
		t.Fatalf("timestamp spelling changed: %v", v)
	}
}
func TestCaseSensitiveSchemaRoutes(t *testing.T) {
	if slug("Error") == slug("error") {
		t.Fatal("schema route collision")
	}
}
func TestUnlockedResourcesAndUnknownDialect(t *testing.T) {
	d := fixture(t, "valid.yaml")
	if _, e := d.Compiler.Compile("https://unlocked.test/schema"); e == nil {
		t.Fatal("unlocked schema fetched")
	}
	if e := d.Compiler.AddResource("https://fixtures.test/custom", Object{"$schema": "https://unlocked.test/dialect", "type": "object"}); e != nil {
		t.Fatal(e)
	}
	if _, e := d.Compiler.Compile("https://fixtures.test/custom"); e == nil {
		t.Fatal("unknown dialect accepted")
	}
}
