// api-docs prepares a source-preserving reference model for Docker API documentation.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	js "github.com/santhosh-tekuri/jsonschema/v6"
	yaml "go.yaml.in/yaml/v4"
)

type Object = map[string]any

const dialect = "https://spec.openapis.org/oas/3.1/dialect/base"

var methods = []string{"get", "put", "post", "delete", "options", "head", "patch", "trace", "query"}

func obj(v any) Object {
	m, _ := v.(map[string]any)
	if m == nil {
		return Object{}
	}
	return m
}
func arr(v any) []any  { a, _ := v.([]any); return a }
func str(v any) string { s, _ := v.(string); return s }
func keys(m Object) []string {
	a := make([]string, 0, len(m))
	for k := range m {
		a = append(a, k)
	}
	sort.Strings(a)
	return a
}
func hash(b []byte) string { s := sha256.Sum256(b); return hex.EncodeToString(s[:]) }
func encoded(v any) []byte {
	b, e := json.MarshalIndent(v, "", "  ")
	if e != nil {
		panic(e)
	}
	return append(b, '\n')
}
func esc(s string) string { return strings.ReplaceAll(strings.ReplaceAll(s, "~", "~0"), "/", "~1") }
func pointer(root any, p string) (any, error) {
	if p == "" {
		return root, nil
	}
	if !strings.HasPrefix(p, "/") {
		return nil, fmt.Errorf("unsupported pointer %s", p)
	}
	for _, k := range strings.Split(p[1:], "/") {
		k = strings.ReplaceAll(strings.ReplaceAll(k, "~1", "/"), "~0", "~")
		switch v := root.(type) {
		case map[string]any:
			var ok bool
			root, ok = v[k]
			if !ok {
				return nil, fmt.Errorf("missing pointer %s", p)
			}
		case []any:
			var n int
			if _, e := fmt.Sscanf(k, "%d", &n); e != nil || n < 0 || n >= len(v) {
				return nil, fmt.Errorf("invalid array pointer %s", p)
			}
			root = v[n]
		default:
			return nil, fmt.Errorf("non-container pointer %s", p)
		}
	}
	return root, nil
}
func readJSON(p string) (any, error) {
	b, e := os.ReadFile(p)
	if e != nil {
		return nil, e
	}
	return js.UnmarshalJSON(bytes.NewReader(b))
}
func writeJSON(p string, v any) error {
	if e := os.MkdirAll(filepath.Dir(p), 0755); e != nil {
		return e
	}
	return os.WriteFile(p, encoded(v), 0644)
}

type Diagnostic struct {
	Rule     string `json:"rule"`
	Pointer  string `json:"pointer"`
	Message  string `json:"message"`
	Waivable bool   `json:"waivable"`
	Excepted bool   `json:"excepted"`
	Reason   string `json:"reason,omitempty"`
}
type Source struct {
	ID         string   `json:"id"`
	Product    string   `json:"product"`
	Title      string   `json:"title"`
	Source     string   `json:"source"`
	Owner      string   `json:"owner"`
	Manual     string   `json:"manual"`
	Connection string   `json:"connection"`
	Auth       string   `json:"auth"`
	Guides     []string `json:"guides"`
}
type Registry struct{ resources map[string]any }

func (r *Registry) Load(uri string) (any, error) {
	if v, ok := r.resources[uri]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("unlocked resource: %s", uri)
}

type Document struct {
	Root         Object
	URI          string
	Registry     *Registry
	Diagnostics  []Diagnostic
	Schemas      map[string]any
	Locations    map[string]int
	Compiler     *js.Compiler
	Compiled     map[string]*js.Schema
	Digest       string
	Source       Source
	SchemaCount  int
	ExampleCount int
}

func (d *Document) issue(rule, p, msg string, waivable bool) {
	d.Diagnostics = append(d.Diagnostics, Diagnostic{Rule: rule, Pointer: p, Message: strings.ReplaceAll(msg, d.URI, "source:"+d.Source.ID), Waivable: waivable})
}
func walk(v any, p string, fn func(Object, string)) {
	switch n := v.(type) {
	case map[string]any:
		fn(n, p)
		for _, k := range keys(n) {
			walk(n[k], p+"/"+esc(k), fn)
		}
	case []any:
		for i, x := range n {
			walk(x, fmt.Sprintf("%s/%d", p, i), fn)
		}
	}
}
func strictSource(file string) (Object, libopenapi.Document, error) {
	b, e := os.ReadFile(file)
	if e != nil {
		return nil, nil, e
	}
	cfg := datamodel.NewDocumentConfiguration()
	cfg.AllowRemoteReferences = false
	cfg.AllowFileReferences = false
	cfg.SkipExternalRefResolution = true
	doc, e := libopenapi.NewDocumentWithConfiguration(b, cfg)
	if e != nil {
		return nil, nil, e
	}
	v, e := strictFragment(file)
	return obj(v), doc, e
}

// YAML timestamp scalars are strings in the JSON data model. Decode source nodes
// explicitly: libopenapi's JSON convenience view normalizes their spelling.
func strictFragment(file string) (any, error) {
	b, e := os.ReadFile(file)
	if e != nil {
		return nil, e
	}
	var node yaml.Node
	if e = yaml.Unmarshal(b, &node); e != nil {
		return nil, e
	}
	v, e := sourceValue(&node, map[*yaml.Node]bool{})
	if e != nil {
		return nil, e
	}
	raw, e := json.Marshal(v)
	if e != nil {
		return nil, e
	}
	return js.UnmarshalJSON(bytes.NewReader(raw))
}
func sourceValue(n *yaml.Node, active map[*yaml.Node]bool) (any, error) {
	if active[n] {
		return nil, errors.New("cyclic YAML alias")
	}
	active[n] = true
	defer delete(active, n)
	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) != 1 {
			return nil, errors.New("expected one YAML document")
		}
		return sourceValue(n.Content[0], active)
	case yaml.AliasNode:
		return sourceValue(n.Alias, active)
	case yaml.MappingNode:
		m := Object{}
		for i := 0; i < len(n.Content); i += 2 {
			k := n.Content[i]
			if k.Kind != yaml.ScalarNode || k.Value == "<<" {
				return nil, errors.New("complex/merge YAML keys are outside the source profile")
			}
			if _, ok := m[k.Value]; ok {
				return nil, fmt.Errorf("duplicate key %s at line %d", k.Value, k.Line)
			}
			v, e := sourceValue(n.Content[i+1], active)
			if e != nil {
				return nil, e
			}
			m[k.Value] = v
		}
		return m, nil
	case yaml.SequenceNode:
		a := []any{}
		for _, c := range n.Content {
			v, e := sourceValue(c, active)
			if e != nil {
				return nil, e
			}
			a = append(a, v)
		}
		return a, nil
	case yaml.ScalarNode:
		if n.Tag == "!!str" || n.Tag == "!!timestamp" {
			return n.Value, nil
		}
		var v any
		e := n.Decode(&v)
		return v, e
	}
	return nil, errors.New("unsupported YAML node")
}

func loadDocument(file, metaDir string) (*Document, error) {
	abs, _ := filepath.Abs(file)
	root, parsed, e := strictSource(abs)
	if e != nil {
		return nil, e
	}
	uri := (&url.URL{Scheme: "file", Path: abs}).String()
	b, _ := os.ReadFile(file)
	d := &Document{Root: root, URI: uri, Diagnostics: []Diagnostic{}, Registry: &Registry{resources: map[string]any{uri: root}}, Schemas: map[string]any{}, Locations: map[string]int{}, Compiled: map[string]*js.Schema{}, Digest: hash(b)}
	// libopenapi's typed model supplies operation views/locations; schema serialization is never used.
	model, modelErr := parsed.BuildV3Model()
	if modelErr == nil && model != nil && model.Model.Paths != nil {
		for p, item := range model.Model.Paths.PathItems.FromOldest() {
			for m, op := range item.GetOperations().FromOldest() {
				if op.GoLow().RootNode != nil {
					d.Locations["/paths/"+esc(p)+"/"+m] = op.GoLow().RootNode.Line
				}
			}
		}
	}
	// All file resources are acquired before compilation, under a bounded source directory.
	var acquire func(any, string) error
	acquire = func(v any, base string) error {
		var failure error
		walk(v, "", func(n Object, _ string) {
			for _, k := range []string{"$ref", "$dynamicRef"} {
				ref, ok := n[k].(string)
				if !ok {
					return
				}
				u, er := url.Parse(ref)
				if er != nil {
					failure = er
					return
				}
				bu, _ := url.Parse(base)
				u = bu.ResolveReference(u)
				u.Fragment = ""
				id := u.String()
				if _, ok := d.Registry.resources[id]; ok {
					return
				}
				if u.Scheme != "file" {
					failure = fmt.Errorf("unlocked reference %s", id)
					return
				}
				rel, er := filepath.Rel(filepath.Dir(abs), u.Path)
				if er != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
					failure = fmt.Errorf("reference escapes package: %s", ref)
					return
				}
				r, er := strictFragment(u.Path)
				if er != nil {
					failure = er
					return
				}
				d.Registry.resources[id] = r
				if er = acquire(r, id); er != nil {
					failure = er
				}
			}
		})
		return failure
	}
	if e = acquire(root, uri); e != nil {
		return nil, e
	}
	lockRaw, e := os.ReadFile(filepath.Join(metaDir, "lock.json"))
	if e != nil {
		return nil, e
	}
	var lock []struct{ File, URI, SHA256 string }
	if e = json.Unmarshal(lockRaw, &lock); e != nil {
		return nil, e
	}
	for _, l := range lock {
		b, e := os.ReadFile(filepath.Join(metaDir, l.File))
		if e != nil {
			return nil, e
		}
		if hash(b) != l.SHA256 {
			return nil, fmt.Errorf("dialect digest mismatch: %s", l.File)
		}
		v, e := js.UnmarshalJSON(bytes.NewReader(b))
		if e != nil {
			return nil, e
		}
		d.Registry.resources[l.URI] = v
		if id := str(obj(v)["$id"]); id != "" {
			d.Registry.resources[id] = v
		}
	}
	d.Compiler = js.NewCompiler()
	d.Compiler.DefaultDraft(js.Draft2020)
	d.Compiler.UseLoader(d.Registry)
	for id, v := range d.Registry.resources {
		if e = d.Compiler.AddResource(id, v); e != nil {
			return nil, e
		}
	}
	return d, nil
}
func (d *Document) resolve(v any) any {
	seen := map[string]bool{}
	for {
		r := str(obj(v)["$ref"])
		if r == "" {
			return v
		}
		if seen[r] {
			return v
		}
		seen[r] = true
		u, e := url.Parse(r)
		if e != nil {
			return v
		}
		base, _ := url.Parse(d.URI)
		u = base.ResolveReference(u)
		frag := u.Fragment
		u.Fragment = ""
		target, ok := d.Registry.resources[u.String()]
		if !ok {
			return v
		}
		resolved, e := pointer(target, frag)
		if e != nil {
			return v
		}
		merged := Object{}
		for k, x := range obj(resolved) {
			merged[k] = x
		}
		for _, k := range []string{"summary", "description"} {
			if x, ok := obj(v)[k]; ok {
				merged[k] = x
			}
		}
		if len(merged) == 0 {
			return resolved
		}
		v = merged
	}
}
func (d *Document) schemaRoots(v any, p string) {
	switch n := v.(type) {
	case map[string]any:
		for _, k := range keys(n) {
			x := n[k]
			q := p + "/" + esc(k)
			if k == "schemas" && p == "/components" {
				for name, s := range obj(x) {
					d.Schemas[q+"/"+esc(name)] = s
				}
			} else if k == "schema" || k == "itemSchema" {
				d.Schemas[q] = x
			} else if !strings.HasPrefix(k, "x-") && k != "example" && k != "examples" && k != "value" {
				d.schemaRoots(x, q)
			}
		}
	case []any:
		for i, x := range n {
			d.schemaRoots(x, fmt.Sprintf("%s/%d", p, i))
		}
	}
}
func (d *Document) validate(metaDir string) {
	if d.Root["openapi"] != "3.2.0" {
		d.issue("S1", "/openapi", "Expected OpenAPI 3.2.0", false)
	}
	if d.Root["jsonSchemaDialect"] != dialect {
		d.issue("S1", "/jsonSchemaDialect", "Expected the selected OAS dialect", false)
	}
	documentSchema, e := d.Compiler.Compile("https://spec.openapis.org/oas/3.2/schema/2025-09-17")
	if e != nil {
		d.issue("structure", "", e.Error(), false)
	} else if e = documentSchema.Validate(d.Root); e != nil {
		d.issue("structure", "", e.Error(), false)
	}
	// OAS Reference Objects also need explicit validation: a document schema cannot prove their target exists.
	walk(d.Root, "", func(n Object, p string) {
		if r := str(n["$ref"]); r != "" {
			u, er := url.Parse(r)
			if er != nil {
				d.issue("reference", p, er.Error(), false)
				return
			}
			base, _ := url.Parse(d.URI)
			u = base.ResolveReference(u)
			frag := u.Fragment
			u.Fragment = ""
			v, ok := d.Registry.resources[u.String()]
			if !ok {
				d.issue("reference", p, "Unregistered resource "+r, false)
			} else if strings.HasPrefix(frag, "/") {
				if _, er = pointer(v, frag); er != nil {
					d.issue("reference", p, er.Error(), false)
				}
			}
		}
	})
	if len(obj(d.Root["webhooks"])) > 0 {
		d.issue("capability", "/webhooks", "Webhook navigation is not supported by the reference renderer", false)
	}
	walk(d.Root, "", func(n Object, p string) {
		if len(obj(n["callbacks"])) > 0 {
			d.issue("capability", p+"/callbacks", "Callback navigation is not supported by the reference renderer", false)
		}
	})
	d.schemaRoots(d.Root, "")
	schemaMeta, e := d.Compiler.Compile(dialect)
	if e != nil {
		d.issue("dialect", "", e.Error(), false)
		return
	}
	for _, p := range keys(d.Schemas) {
		s := d.Schemas[p]
		d.SchemaCount++
		if e = schemaMeta.Validate(s); e != nil {
			d.issue("schema", p, e.Error(), false)
			continue
		}
		compiled, e := d.Compiler.Compile(d.URI + "#" + p)
		if e != nil {
			d.issue("schema", p, e.Error(), false)
			continue
		}
		d.Compiled[p] = compiled
		d.schemaExamples(s, p)
	}
	// Media, parameter and header examples are separate from examples nested inside schemas.
	walk(d.Root, "", func(n Object, p string) {
		schemaKey := "schema"
		if _, ok := n["itemSchema"]; ok {
			schemaKey = "itemSchema"
		}
		c := d.Compiled[p+"/"+schemaKey]
		if c == nil {
			return
		}
		if v, ok := n["example"]; ok {
			d.example(c, v, p+"/example")
		}
		for _, name := range keys(obj(n["examples"])) {
			ex := obj(d.resolve(obj(n["examples"])[name]))
			if v, ok := ex["value"]; ok {
				d.example(c, v, p+"/examples/"+esc(name)+"/value")
			}
			if _, ok := ex["externalValue"]; ok {
				d.issue("example-external", p+"/examples/"+esc(name), "External example requires a locked media fixture", true)
			}
		}
	})
	ids := map[string]bool{}
	tags := map[string]string{}
	for _, raw := range arr(d.Root["tags"]) {
		t := obj(raw)
		tags[str(t["name"])] = str(t["kind"])
	}
	for _, op := range d.operations() {
		id := str(op["id"])
		p := str(op["pointer"])
		if id == "" || ids[id] {
			d.issue("S4", p, "Operation ID must be present and unique", false)
		}
		ids[id] = true
		opTags := arr(op["tags"])
		if len(opTags) == 0 || tags[str(opTags[0])] != "nav" {
			d.issue("S5", p+"/tags", "The first operation tag must identify a declared navigation group", false)
		}
		for _, tag := range opTags {
			if _, ok := tags[str(tag)]; !ok {
				d.issue("S5", p+"/tags", "Operation tag is undeclared: "+str(tag), false)
			}
		}
		if len(arr(op["servers"])) == 0 && d.Source.Connection != "unix" {
			d.issue("S6", p, "Operation needs effective servers or a local connection profile", false)
		}
		if strings.TrimSpace(str(op["description"])) == "" {
			d.issue("S4", p+"/description", "Operation description required", true)
		}
		for _, pr := range arr(op["parameters"]) {
			param := obj(pr)
			if str(param["description"]) == "" {
				d.issue("S8", str(param["pointer"])+"/description", "Parameter description requires editorial review", true)
			}
		}
		for _, variant := range arr(op["variants"]) {
			v := obj(variant)
			if str(v["media"]) != "" && len(arr(v["examples"])) == 0 {
				d.issue("S11", str(v["pointer"]), "Media variant needs a reviewed example or transfer fixture", true)
			}
		}
	}
}
func (d *Document) schemaExamples(s any, p string) {
	n := obj(s)
	if len(n) == 0 {
		return
	}
	if _, ok := n["$schema"]; ok && n["$schema"] != dialect {
		d.issue("S1", p+"/$schema", "Schema dialect override is outside the profile", false)
	}
	c, e := d.Compiler.Compile(d.URI + "#" + p)
	if e == nil {
		if v, ok := n["example"]; ok {
			d.example(c, v, p+"/example")
		}
		for i, v := range arr(n["examples"]) {
			d.example(c, v, fmt.Sprintf("%s/examples/%d", p, i))
		}
	}
	for k, x := range n {
		switch k {
		case "properties", "patternProperties", "$defs", "dependentSchemas":
			for _, name := range keys(obj(x)) {
				d.schemaExamples(obj(x)[name], p+"/"+k+"/"+esc(name))
			}
		case "allOf", "anyOf", "oneOf", "prefixItems":
			for i, v := range arr(x) {
				d.schemaExamples(v, fmt.Sprintf("%s/%s/%d", p, k, i))
			}
		case "items", "additionalProperties", "unevaluatedProperties", "contains", "not", "if", "then", "else", "contentSchema", "propertyNames":
			d.schemaExamples(x, p+"/"+k)
		}
	}
}
func (d *Document) example(c *js.Schema, v any, p string) {
	d.ExampleCount++
	if e := c.Validate(v); e != nil {
		d.issue("example", p, e.Error(), true)
	}
}
func (d *Document) effective(root, item, op Object, k string) any {
	if v, ok := op[k]; ok {
		return v
	}
	if k == "servers" {
		if v, ok := item[k]; ok {
			return v
		}
	}
	return root[k]
}
func (d *Document) parameters(item, op Object, p string) []any {
	a := []any{}
	positions := map[string]int{}
	for i, container := range []Object{item, op} {
		for j, raw := range arr(container["parameters"]) {
			v := obj(d.resolve(raw))
			c := Object{}
			for k, x := range v {
				c[k] = x
			}
			prefix := p
			if i == 0 {
				prefix = p[:strings.LastIndex(p, "/")]
			}
			c["pointer"] = fmt.Sprintf("%s/parameters/%d", prefix, j)
			key := str(v["in"]) + ":" + str(v["name"])
			if n, ok := positions[key]; ok {
				a[n] = c
			} else {
				positions[key] = len(a)
				a = append(a, c)
			}
		}
	}
	return a
}
func exampleList(v Object) []any {
	out := []any{}
	if x, ok := v["example"]; ok {
		out = append(out, Object{"name": "Example", "value": x})
	}
	for _, name := range keys(obj(v["examples"])) {
		x := obj(obj(v["examples"])[name])
		if value, ok := x["value"]; ok {
			out = append(out, Object{"name": name, "value": value})
		}
	}
	return out
}
func (d *Document) variants(op Object, p string) []any {
	out := []any{}
	add := func(direction, status string, raw any, ptr string) {
		container := obj(d.resolve(raw))
		content := obj(container["content"])
		if len(content) == 0 {
			out = append(out, Object{"direction": direction, "status": status, "description": container["description"], "headers": container["headers"], "media": "", "pointer": ptr, "examples": []any{}})
		}
		for _, media := range keys(content) {
			v := obj(content[media])
			a := Object{"direction": direction, "status": status, "description": container["description"], "headers": container["headers"], "media": media, "pointer": ptr + "/content/" + esc(media), "examples": d.mediaExamples(v), "required": container["required"]}
			for _, k := range []string{"schema", "itemSchema", "encoding"} {
				if s, ok := v[k]; ok {
					a[k] = s
				}
			}
			out = append(out, a)
		}
	}
	if v, ok := op["requestBody"]; ok {
		add("Request", "", v, p+"/requestBody")
	}
	for _, code := range keys(obj(op["responses"])) {
		add("Response", code, obj(op["responses"])[code], p+"/responses/"+esc(code))
	}
	return out
}
func (d *Document) operations() []Object {
	out := []Object{}
	for _, path := range keys(obj(d.Root["paths"])) {
		item := obj(d.resolve(obj(d.Root["paths"])[path]))
		ms := append([]string{}, methods...)
		for m := range obj(item["additionalOperations"]) {
			ms = append(ms, m)
		}
		for _, method := range ms {
			raw, ok := item[method]
			extra := !ok
			if !ok {
				raw, ok = obj(item["additionalOperations"])[method]
			}
			if !ok {
				continue
			}
			op := obj(raw)
			p := "/paths/" + esc(path) + "/" + method
			if extra {
				p = "/paths/" + esc(path) + "/additionalOperations/" + esc(method)
			}
			v := Object{"id": op["operationId"], "method": strings.ToUpper(method), "path": path, "pointer": p, "line": d.Locations[p], "summary": op["summary"], "description": op["description"], "tags": op["tags"], "deprecated": op["deprecated"], "servers": d.effective(d.Root, item, op, "servers"), "security": d.effective(d.Root, item, op, "security"), "parameters": d.parameters(item, op, p), "variants": d.variants(op, p), "raw": op}
			out = append(out, v)
		}
	}
	return out
}
func main() {
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
func run() error {
	if len(os.Args) < 3 {
		return errors.New("usage: api-docs check|generate|inspect ROOT [--allow-known-issues]")
	}
	command := os.Args[1]
	if command != "check" && command != "generate" && command != "inspect" {
		return fmt.Errorf("unknown command: %s", command)
	}
	if len(os.Args) > 4 || (len(os.Args) == 4 && os.Args[3] != "--allow-known-issues") {
		return errors.New("unexpected arguments")
	}
	root, _ := filepath.Abs(os.Args[2])
	dir := filepath.Join(root, "hack/api-docs")
	meta := filepath.Join(dir, "dialects")
	catalogRaw, e := os.ReadFile(filepath.Join(dir, "catalog.json"))
	if e != nil {
		return e
	}
	var catalog struct {
		APIs       []Source `json:"apis"`
		LegacyAPIs []Object `json:"legacyAPIs"`
	}
	if e = json.Unmarshal(catalogRaw, &catalog); e != nil {
		return e
	}
	allowKnown := len(os.Args) > 3 && os.Args[3] == "--allow-known-issues"
	exceptions, e := readJSON(filepath.Join(dir, "known-issues.json"))
	if e != nil {
		return e
	}
	reports := []any{}
	models := []any{}
	blocking := 0
	for _, src := range catalog.APIs {
		d, e := loadDocument(filepath.Join(root, src.Source), meta)
		if e != nil {
			return fmt.Errorf("%s: %w", src.ID, e)
		}
		d.Source = src
		d.validate(meta)
		for i := range d.Diagnostics {
			diag := &d.Diagnostics[i]
			for _, raw := range arr(exceptions) {
				x := obj(raw)
				if diag.Waivable && x["api"] == src.ID && x["sha256"] == d.Digest && x["rule"] == diag.Rule && x["pointer"] == diag.Pointer && x["messageSha256"] == hash([]byte(diag.Message)) && str(x["reason"]) != "" {
					diag.Excepted = true
					diag.Reason = str(x["reason"])
					break
				}
			}
			if !diag.Excepted || !allowKnown {
				blocking++
			}
		}
		reports = append(reports, Object{"api": src.ID, "sha256": d.Digest, "schemas": d.SchemaCount, "examples": d.ExampleCount, "operations": len(d.operations()), "diagnostics": d.Diagnostics})
		models = append(models, d.model())
		fmt.Printf("%s: %d operations, %d schema roots, %d examples, %d diagnostics\n", src.ID, len(d.operations()), d.SchemaCount, d.ExampleCount, len(d.Diagnostics))
	}
	out := filepath.Join(root, "tmp/api-reference")
	if e = writeJSON(filepath.Join(out, "validation.json"), reports); e != nil {
		return e
	}
	if command == "inspect" {
		return nil
	}
	if blocking > 0 {
		return fmt.Errorf("%d blocking diagnostics; see tmp/api-reference/validation.json", blocking)
	}
	if command == "generate" {
		if e = writeJSON(filepath.Join(out, "data/api-reference.json"), Object{"modelVersion": 1, "apis": models, "legacyAPIs": catalog.LegacyAPIs}); e != nil {
			return e
		}
	}
	return nil
}
