package main

import (
	"archive/zip"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func route(id string) string {
	switch id {
	case "governance":
		return "/reference/api/ai-governance/"
	case "hub", "dvp", "registry":
		return "/reference/api/" + id + "/latest/"
	case "engine-1.56":
		return "/reference/api/engine/version/v1.56/"
	default:
		return "/api-prototype/" + id + "/"
	}
}
func slug(s string) string { return url.PathEscape(s) }
func (d *Document) model() Object {
	ops := d.operations()
	schemas := []any{}
	schemaURLs := Object{}
	for _, name := range keys(obj(obj(d.Root["components"])["schemas"])) {
		p := "/components/schemas/" + esc(name)
		u := route(d.Source.ID) + "schemas/" + slug(name) + "/"
		schemaURLs["#"+p] = u
		schemas = append(schemas, Object{"name": name, "pointer": p, "url": u, "schema": obj(obj(d.Root["components"])["schemas"])[name]})
	}
	for _, op := range ops {
		for _, raw := range arr(op["variants"]) {
			v := obj(raw)
			c := d.Compiled[str(v["pointer"])+"/schema"]
			if c != nil {
				for _, rawEx := range arr(v["examples"]) {
					ex := obj(rawEx)
					ex["valid"] = c.Validate(ex["value"]) == nil
				}
			}
		}
	}
	for _, op := range ops {
		op["url"] = route(d.Source.ID) + "operations/" + slug(str(op["id"])) + "/"
		op["securitySchemes"] = obj(d.Root["components"])["securitySchemes"]
		op["curl"], op["curlNotes"] = curlExample(d.Source, op)
		op["references"] = refs(op["raw"], schemaURLs)
		op["requestSchema"] = firstRequestSchema(op)
	}
	assumptions := []any{}
	for _, raw := range arr(mustJSON(filepath.Join("prototypes/api-docs/migrations", d.Source.ID+".json"))["changes"]) {
		c := obj(raw)
		if c["classification"] == "provisional assumption" {
			assumptions = append(assumptions, Object{"id": c["id"], "pointer": c["pointer"], "rationale": c["rationale"]})
		}
	}
	return Object{"id": d.Source.ID, "product": d.Source.Product, "title": d.Source.Title, "version": obj(d.Root["info"])["version"], "description": obj(d.Root["info"])["description"], "url": route(d.Source.ID), "manual": d.Source.Manual, "guides": d.Source.Guides, "connection": d.Source.Connection, "auth": d.Source.Auth, "servers": d.Root["servers"], "securitySchemes": obj(d.Root["components"])["securitySchemes"], "tags": d.Root["tags"], "operations": ops, "schemas": schemas, "schemaURLs": schemaURLs, "digest": d.Digest, "owner": d.Source.Owner, "source": d.Source.Source, "diagnostics": d.Diagnostics, "assumptions": assumptions, "schemaCount": d.SchemaCount, "exampleCount": d.ExampleCount}
}
func mustJSON(p string) Object { v, _ := readJSON(p); return obj(v) }
func refs(v any, urls Object) []any {
	found := map[string]bool{}
	walk(v, "", func(n Object, _ string) {
		if s := str(n["$ref"]); s != "" {
			found[s] = true
		}
	})
	out := []any{}
	names := []string{}
	for n := range found {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		out = append(out, Object{"ref": n, "url": urls[n]})
	}
	return out
}
func firstRequestSchema(op Object) any {
	for _, v := range arr(op["variants"]) {
		m := obj(v)
		if m["direction"] == "Request" {
			return m["schema"]
		}
	}
	return nil
}
func shell(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'" }
func scalar(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return strings.TrimSpace(string(encoded(v)))
}
func parameterValue(p Object) (any, bool) {
	if v, ok := p["example"]; ok {
		return v, true
	}
	s := obj(p["schema"])
	if v, ok := s["example"]; ok {
		return v, true
	}
	if a := arr(s["examples"]); len(a) > 0 {
		return a[0], true
	}
	return nil, false
}
func curlExample(src Source, op Object) (string, []string) {
	notes := []string{}
	args := []string{"curl"}
	if op["method"] == "HEAD" {
		args = append(args, "--head")
	} else {
		args = append(args, "--request "+str(op["method"]))
	}
	server := ""
	if a := arr(op["servers"]); len(a) > 0 {
		m := obj(a[0])
		server = str(m["url"])
		for name, v := range obj(m["variables"]) {
			server = strings.ReplaceAll(server, "{"+name+"}", str(obj(v)["default"]))
		}
	}
	if src.Connection == "unix" {
		args = append(args, "--unix-socket \"${DOCKER_SOCKET:-/var/run/docker.sock}\"")
		if !strings.HasPrefix(server, "/") {
			server = "/v" + strings.TrimPrefix(src.ID, "engine-")
		}
		server = "http://localhost" + server
	}
	if server == "" {
		server = "https://<SERVER>"
		notes = append(notes, "Set the API server address.")
	}
	target := strings.TrimRight(server, "/") + str(op["path"])
	query := []string{}
	for _, raw := range arr(op["parameters"]) {
		p := obj(raw)
		name := str(p["name"])
		v, has := parameterValue(p)
		required, _ := p["required"].(bool)
		if !has && !required {
			continue
		}
		value := "<" + strings.ToUpper(name) + ">"
		if has {
			value = scalar(v)
		}
		location := str(p["in"])
		if location == "path" {
			part := value
			if has {
				part = url.PathEscape(value)
			}
			target = strings.ReplaceAll(target, "{"+name+"}", part)
		} else if location == "header" {
			args = append(args, "--header "+shell(name+": "+value))
		} else if location == "query" {
			style := str(p["style"])
			if style == "" {
				style = "form"
			}
			explode := true
			if x, ok := p["explode"].(bool); ok {
				explode = x
			}
			if list, ok := v.([]any); ok && style == "form" {
				parts := []string{}
				for _, x := range list {
					parts = append(parts, scalar(x))
				}
				if explode {
					for _, x := range parts {
						query = append(query, url.QueryEscape(name)+"="+url.QueryEscape(x))
					}
				} else {
					query = append(query, url.QueryEscape(name)+"="+url.QueryEscape(strings.Join(parts, ",")))
				}
			} else if len(obj(v)) > 0 || style != "form" {
				notes = append(notes, "Serialize "+name+" using its documented "+style+" rules; this parameter is not generated.")
			} else {
				query = append(query, url.QueryEscape(name)+"="+url.QueryEscape(value))
			}
		} else {
			notes = append(notes, "Supply "+name+" using its documented "+location+" encoding.")
		}
	}
	if len(query) > 0 {
		target += "?" + strings.Join(query, "&")
	}
	security := arr(op["security"])
	if len(security) > 0 {
		first := obj(security[0])
		for _, scheme := range keys(first) {
			definition := obj(obj(op["securitySchemes"])[scheme])
			switch {
			case definition["type"] == "http" && definition["scheme"] == "bearer":
				token := "TOKEN"
				if scheme == "scimToken" {
					token = "SCIM_TOKEN"
				}
				if scheme == "registryToken" {
					token = "REGISTRY_TOKEN"
				}
				args = append(args, "--header \"Authorization: Bearer ${"+token+"}\"")
			case definition["type"] == "apiKey" && definition["in"] == "header":
				args = append(args, "--header \""+str(definition["name"])+": ${API_KEY}\"")
			default:
				notes = append(notes, "Configure authentication scheme "+scheme+" using its documented transport or credential format.")
			}
		}
		if len(security) > 1 {
			notes = append(notes, "This example uses the first authentication alternative. Review the complete requirements.")
		}
	}
	for _, raw := range arr(op["variants"]) {
		v := obj(raw)
		if v["direction"] != "Request" {
			continue
		}
		media := str(v["media"])
		if media == "" {
			continue
		}
		args = append(args, "--header "+shell("Content-Type: "+media))
		examples := arr(v["examples"])
		if len(examples) > 0 && obj(examples[0])["valid"] != false && strings.Contains(media, "json") {
			args = append(args, "--data-raw "+shell(strings.TrimSpace(string(encoded(obj(examples[0])["value"])))))
		} else {
			args = append(args, "--data-binary @request-body")
			notes = append(notes, "Prepare request-body using the selected media type and schema.")
		}
		break
	}
	args = append(args, shell(target))
	return strings.Join(args, " \\\n  "), notes
}
func publishSources(root string, sources []Source) error {
	out := filepath.Join(root, "tmp/api-prototype/static/api-prototype/sources")
	for _, s := range sources {
		dir := filepath.Join(out, s.ID)
		if e := os.MkdirAll(dir, 0755); e != nil {
			return e
		}
		f, e := os.Create(filepath.Join(dir, "source.zip"))
		if e != nil {
			return e
		}
		z := zip.NewWriter(f)
		for _, pair := range [][2]string{{"converted/" + s.ID + ".yaml", "openapi.yaml"}, {"original/" + s.ID + ".yaml", "original.yaml"}, {"migrations/" + s.ID + ".json", "migration.json"}, {"migrations/" + s.ID + ".patch", "migration.patch"}, {"source-lock.json", "source-lock.json"}, {"exceptions.json", "exceptions.json"}} {
			b, e := os.ReadFile(filepath.Join(root, "prototypes/api-docs", pair[0]))
			if e != nil {
				return e
			}
			if e = os.WriteFile(filepath.Join(dir, pair[1]), b, 0644); e != nil {
				return e
			}
			w, e := z.Create(pair[1])
			if e != nil {
				return e
			}
			if _, e = w.Write(b); e != nil {
				return e
			}
		}
		if e = z.Close(); e != nil {
			return e
		}
		if e = f.Close(); e != nil {
			return e
		}
	}
	return nil
}

// Examples are annotations: retrieve them without flattening schema constraints.
func (d *Document) mediaExamples(v Object) []any {
	out := exampleList(v)
	if len(out) > 0 {
		return out
	}
	s := v["schema"]
	seen := map[string]bool{}
	for s != nil {
		n := obj(s)
		if x, ok := n["example"]; ok {
			return []any{Object{"name": "Schema example", "value": x}}
		}
		if xs := arr(n["examples"]); len(xs) > 0 {
			for _, x := range xs {
				out = append(out, Object{"name": "Schema example", "value": x})
			}
			return out
		}
		ref := str(n["$ref"])
		if ref == "" || seen[ref] {
			break
		}
		seen[ref] = true
		s = d.resolve(n)
	}
	return out
}
