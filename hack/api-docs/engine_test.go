package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func engineModelOperation(t *testing.T, model Object, id string) Object {
	t.Helper()
	for _, op := range model["operations"].([]Object) {
		if op["id"] == id {
			return op
		}
	}
	t.Fatalf("operation %q missing from presentation model", id)
	return nil
}

func TestEngine157Presentation(t *testing.T) {
	b, err := os.ReadFile("sources.json")
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		APIs []Source `json:"apis"`
	}
	if err := json.Unmarshal(b, &manifest); err != nil {
		t.Fatal(err)
	}
	var src Source
	for _, candidate := range manifest.APIs {
		if candidate.ID == "engine-1.57" {
			src = candidate
		}
	}
	if src.ID == "" {
		t.Fatal("Engine 1.57 must be registered")
	}
	d, err := loadDocument(filepath.Join("../..", src.Source), "validation/dialects")
	if err != nil {
		t.Fatal(err)
	}
	d.Source = src
	d.validate("")
	for _, diagnostic := range d.Diagnostics {
		if diagnostic.Severity != "warning" {
			t.Errorf("Engine contract failed: %+v", diagnostic)
		}
	}
	m := d.model()
	const base = "/reference/api/engine/version/v1.57/"
	if m["url"] != base || m["sourceURL"] != "/reference/api/engine/version/v1.57.yaml" || m["authoredOverview"] != true {
		t.Fatalf("incorrect Engine routes/overview: %v %v %v", m["url"], m["sourceURL"], m["authoredOverview"])
	}
	for _, op := range m["operations"].([]Object) {
		if !strings.HasPrefix(str(op["url"]), base+"operations/") {
			t.Errorf("unversioned operation route: %v", op["url"])
		}
		if command := str(op["curl"]); command != "" {
			if !strings.Contains(command, `--unix-socket "${DOCKER_SOCKET:-/var/run/docker.sock}"`) || !strings.Contains(command, "http://localhost/v1.57/") {
				t.Errorf("incorrect Unix request for %v: %s", op["id"], command)
			}
			if strings.Contains(command, "Authorization:") {
				t.Errorf("daemon authentication invented for %v: %s", op["id"], command)
			}
		}
	}
	if len(arr(m["schemas"])) == 0 {
		t.Fatal("Engine schemas missing")
	}
	for _, raw := range arr(m["schemas"]) {
		if u := str(obj(raw)["url"]); !strings.HasPrefix(u, base+"schemas/") {
			t.Errorf("unversioned schema route: %s", u)
		}
	}
	head := engineModelOperation(t, m, "SystemPingHead")
	if !strings.Contains(str(head["curl"]), "--head") || strings.Contains(str(head["curl"]), "--request HEAD") {
		t.Fatalf("HEAD must use curl's header-only mode: %v", head["curl"])
	}
	push := engineModelOperation(t, m, "ImagePush")
	if !strings.Contains(str(push["curl"]), "X-Registry-Auth:") {
		t.Fatal("required registry credential header omitted")
	}
}

func TestEngineTagPresentation(t *testing.T) {
	for _, tc := range []struct {
		name, summary string
		tag           Object
		declared      bool
	}{
		{"display name", "Containers", Object{"name": "Container", "x-displayName": "Containers"}, true},
		{"explicit summary", "Container operations", Object{"name": "Container", "summary": "Container operations", "x-displayName": "Containers"}, true},
		{"name fallback", "Container", Object{"name": "Container"}, true},
		{"undeclared", "Container", Object{"name": "Container"}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := fixture(t, "valid.yaml")
			if tc.declared {
				d.Root["tags"] = append(arr(d.Root["tags"]), tc.tag)
			}
			op := obj(obj(obj(d.Root["paths"])["/public"])["get"])
			op["tags"] = []any{"Container"}
			delete(op, "description")
			d.validate("")
			if len(d.Diagnostics) != 0 {
				t.Fatalf("optional kind, undeclared tags and summary-only operation must be accepted: %+v", d.Diagnostics)
			}
			found := 0
			for _, raw := range arr(d.model()["tags"]) {
				tag := obj(raw)
				if tag["name"] == "Container" {
					found++
					if tag["kind"] != "nav" || tag["summary"] != tc.summary {
						t.Errorf("incorrect navigation tag: %v", tag)
					}
				}
			}
			if found != 1 {
				t.Errorf("expected one Container navigation group, got %d", found)
			}
			if _, ok := tc.tag["kind"]; ok {
				t.Fatal("presentation defaults mutated source tag")
			}
		})
	}
	for _, first := range []bool{false, true} {
		t.Run(map[bool]string{false: "secondary badge", true: "first badge"}[first], func(t *testing.T) {
			d := fixture(t, "valid.yaml")
			d.Root["tags"] = append(arr(d.Root["tags"]), Object{"name": "experimental", "kind": "badge", "summary": "Experimental"})
			op := obj(obj(obj(d.Root["paths"])["/public"])["get"])
			op["tags"] = []any{"operations", "experimental"}
			if first {
				op["tags"] = []any{"experimental", "operations"}
			}
			d.validate("")
			if first {
				if len(d.Diagnostics) != 1 || d.Diagnostics[0].Rule != "S5" || d.Diagnostics[0].Severity != "error" {
					t.Fatalf("badge-first navigation must remain a contract error: %+v", d.Diagnostics)
				}
			} else if len(d.Diagnostics) != 0 {
				t.Fatalf("secondary badge rejected: %+v", d.Diagnostics)
			}
			for _, raw := range arr(d.model()["tags"]) {
				if tag := obj(raw); tag["name"] == "experimental" && tag["kind"] != "badge" {
					t.Errorf("explicit badge kind overwritten: %v", tag)
				}
			}
		})
	}
}

func TestEngineInvalidSamplesAreWarningsAndOmitted(t *testing.T) {
	for _, mixed := range []bool{false, true} {
		t.Run(map[bool]string{false: "invalid only", true: "mixed valid and invalid"}[mixed], func(t *testing.T) {
			d := fixture(t, "valid.yaml")
			op := obj(obj(obj(d.Root["paths"])["/public"])["get"])
			media := func() Object {
				m := Object{"schema": Object{"type": "boolean"}, "example": "INVALID_SAMPLE"}
				if mixed {
					delete(m, "example")
					m["examples"] = Object{"bad": Object{"value": "INVALID_SAMPLE"}, "good": Object{"value": false}}
				}
				return m
			}
			op["requestBody"] = Object{"required": true, "content": Object{"application/json": media()}}
			obj(obj(op["responses"])["200"])["content"] = Object{"application/json": media()}
			op["parameters"] = []any{
				Object{"name": "optional", "in": "query", "description": "Optional count.", "schema": Object{"type": "integer"}, "example": "INVALID_OPTIONAL"},
				Object{"name": "required", "in": "query", "required": true, "description": "Required count.", "schema": Object{"type": "integer"}, "example": "INVALID_REQUIRED"},
				Object{"name": "zero", "in": "query", "description": "Zero count.", "schema": Object{"type": "integer"}, "example": jsonNumber("0")},
			}
			d.validate("")
			if len(d.Diagnostics) != 4 {
				t.Fatalf("expected request, response and two parameter warnings: %+v", d.Diagnostics)
			}
			for _, diagnostic := range d.Diagnostics {
				if diagnostic.Rule != "example" || diagnostic.Severity != "warning" {
					t.Errorf("invalid sample must be a warning: %+v", diagnostic)
				}
			}
			modelOp := engineModelOperation(t, d.model(), "getPublic")
			for _, raw := range arr(modelOp["variants"]) {
				variant := obj(raw)
				examples := arr(variant["examples"])
				if mixed {
					if len(examples) != 1 || obj(examples[0])["value"] != false {
						t.Errorf("valid false sample must survive filtering: %v", examples)
					}
				} else if len(examples) != 0 {
					t.Errorf("invalid media sample leaked into presentation: %v", examples)
				}
			}
			for _, raw := range arr(modelOp["parameters"]) {
				p := obj(raw)
				if p["name"] != "zero" {
					if _, usable := parameterValue(p); usable || p["exampleUsable"] != false {
						t.Errorf("invalid parameter sample remains usable: %v", p)
					}
				}
			}
			command := str(modelOp["curl"])
			if strings.Contains(command, "INVALID_") || strings.Contains(command, "optional=") || !strings.Contains(command, "zero=0") || !strings.Contains(command, "required=%3CREQUIRED%3E") {
				t.Errorf("incorrect parameter sample filtering: %s", command)
			}
			if mixed {
				if !strings.Contains(command, "--data-raw 'false'") {
					t.Errorf("valid false request sample lost: %s", command)
				}
			} else if !strings.Contains(command, "--data-binary @request-body") || strings.Contains(command, "--data-raw") {
				t.Errorf("invalid request sample must fall back to a body file: %s", command)
			}
		})
	}
}

func TestEngineSampleWarningsDoNotWeakenGenerationContracts(t *testing.T) {
	for _, invalidContract := range []bool{false, true} {
		t.Run(map[bool]string{false: "warnings generate", true: "contract error blocks"}[invalidContract], func(t *testing.T) {
			d := fixture(t, "valid.yaml")
			response := obj(obj(obj(obj(obj(d.Root["paths"])["/public"])["get"])["responses"])["200"])
			obj(obj(response["content"])["application/json"])["example"] = "INVALID_SAMPLE"
			if invalidContract {
				delete(obj(d.Root["info"]), "version")
			}
			root := t.TempDir()
			write := func(name string, b []byte) {
				t.Helper()
				p := filepath.Join(root, name)
				if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(p, b, 0644); err != nil {
					t.Fatal(err)
				}
			}
			write("spec.json", encoded(d.Root))
			event, err := os.ReadFile("testdata/schemas/event.yaml")
			if err != nil {
				t.Fatal(err)
			}
			write("schemas/event.yaml", event)
			write("hack/api-docs/sources.json", encoded(Object{"apis": []Source{{ID: "engine-test", Source: "spec.json", Connection: "unix"}}}))
			for _, name := range []string{"lock.json", "oas-document.json", "oas-dialect.json", "oas-vocabulary.json"} {
				b, err := os.ReadFile(filepath.Join("validation/dialects", name))
				if err != nil {
					t.Fatal(err)
				}
				write("hack/api-docs/validation/dialects/"+name, b)
			}
			args := os.Args
			t.Cleanup(func() { os.Args = args })
			os.Args = []string{"api-docs", "generate", root}
			err = run()
			if invalidContract {
				if err == nil || !strings.Contains(err.Error(), "blocking diagnostics") {
					t.Fatalf("invalid contract must block generation: %v", err)
				}
				if _, err := os.Stat(filepath.Join(root, "data/api-reference.json")); !os.IsNotExist(err) {
					t.Fatalf("failed generation published data: %v", err)
				}
			} else if err != nil {
				t.Fatalf("sample warnings must permit generation: %v", err)
			} else if _, err := os.Stat(filepath.Join(root, "data/api-reference.json")); err != nil {
				t.Fatalf("generation did not publish presentation data: %v", err)
			}
			report, err := os.ReadFile(filepath.Join(root, "tmp/api-reference/validation.json"))
			if err != nil {
				t.Fatal(err)
			}
			var reports []struct {
				Diagnostics []Diagnostic `json:"diagnostics"`
			}
			if err := json.Unmarshal(report, &reports); err != nil {
				t.Fatal(err)
			}
			warning, contractError := false, false
			for _, diagnostic := range reports[0].Diagnostics {
				warning = warning || (diagnostic.Rule == "example" && diagnostic.Severity == "warning")
				contractError = contractError || (diagnostic.Rule == "structure" && diagnostic.Severity == "error")
			}
			if !warning || contractError != invalidContract {
				t.Fatalf("incorrect warning/error report: %+v", reports)
			}
		})
	}
}

func TestEngineHTTPUpgradeHasNoOrdinaryCurl(t *testing.T) {
	d := fixture(t, "valid.yaml")
	op := obj(obj(obj(d.Root["paths"])["/public"])["get"])
	obj(op["responses"])["101"] = Object{"description": "Switching protocols"}
	d.validate("")
	if len(d.Diagnostics) != 0 {
		t.Fatalf("upgrade response must be supported: %+v", d.Diagnostics)
	}
	modelOp := engineModelOperation(t, d.model(), "getPublic")
	notes, ok := modelOp["curlNotes"].([]string)
	if modelOp["curl"] != "" || !ok || len(notes) == 0 {
		t.Fatalf("HTTP 101 requires upgrade guidance instead of ordinary curl: %v", modelOp)
	}
	if command := str(engineModelOperation(t, d.model(), "getItems")["curl"]); command == "" {
		t.Fatal("ordinary HTTP operation lost its curl example")
	}
}

func TestReferencedInvalidSamplesAreOmitted(t *testing.T) {
	for _, mixed := range []bool{false, true} {
		t.Run(map[bool]string{false: "invalid only", true: "mixed samples"}[mixed], func(t *testing.T) {
			d := fixture(t, "valid.yaml")
			components := obj(d.Root["components"])
			obj(components["schemas"])["Flag"] = Object{"type": "boolean"}
			media := func() Object {
				m := Object{"schema": Object{"$ref": "#/components/schemas/Flag"}, "example": "INVALID_SAMPLE"}
				if mixed {
					delete(m, "example")
					m["examples"] = Object{"bad": Object{"value": "INVALID_SAMPLE"}, "good": Object{"value": false}}
				}
				return m
			}
			components["requestBodies"] = Object{
				"Flag":  Object{"required": true, "content": Object{"application/json": media()}},
				"Alias": Object{"$ref": "#/components/requestBodies/Flag"},
			}
			components["responses"] = Object{"Flag": Object{"description": "Flag.", "content": Object{"application/json": media()}}}
			components["parameters"] = Object{"Flag": Object{"name": "flag", "in": "query", "required": true, "description": "Flag.", "schema": Object{"$ref": "#/components/schemas/Flag"}, "example": "INVALID_PARAMETER"}}
			item := obj(obj(d.Root["paths"])["/public"])
			item["parameters"] = []any{Object{"$ref": "#/components/parameters/Flag"}}
			op := obj(item["get"])
			op["requestBody"] = Object{"$ref": "#/components/requestBodies/Alias"}
			obj(op["responses"])["200"] = Object{"$ref": "#/components/responses/Flag"}
			source := string(encoded(d.Root))
			d.validate("")
			if len(d.Diagnostics) != 3 {
				t.Fatalf("expected three invalid component sample warnings: %+v", d.Diagnostics)
			}
			for _, diag := range d.Diagnostics {
				if diag.Severity != "warning" || diag.Rule != "example" {
					t.Fatalf("unexpected contract error: %+v", diag)
				}
			}
			modelOp := engineModelOperation(t, d.model(), "getPublic")
			for _, raw := range arr(modelOp["variants"]) {
				examples := arr(obj(raw)["examples"])
				if mixed {
					if len(examples) != 1 || obj(examples[0])["value"] != false {
						t.Fatalf("valid false sample lost: %v", examples)
					}
				} else if len(examples) != 0 {
					t.Fatalf("invalid component sample retained: %v", examples)
				}
			}
			command := str(modelOp["curl"])
			if strings.Contains(command, "INVALID_") || !strings.Contains(command, "flag=%3CFLAG%3E") {
				t.Fatalf("invalid component sample leaked into request: %s", command)
			}
			if mixed {
				if !strings.Contains(command, "--data-raw 'false'") {
					t.Fatal(command)
				}
			} else if !strings.Contains(command, "--data-binary @request-body") {
				t.Fatal(command)
			}
			if string(encoded(d.Root)) != source {
				t.Fatal("source annotations mutated")
			}
		})
	}
}
