package main

import (
	"maps"
	"net/url"
	"sort"
	"strings"

	js "github.com/santhosh-tekuri/jsonschema/v6"
)

func route(src Source) string {
	if src.URL != "" {
		return strings.TrimRight(src.URL, "/") + "/"
	}
	return "/reference/api/" + src.ID + "/latest/"
}
func slug(s string) string { return url.PathEscape(s) }

// Presentation pointers describe operation views. Reference Objects can place
// their schemas elsewhere in the source, so follow references before compiling.
func (d *Document) presentationSchema(p string) *js.Schema {
	if c := d.Compiled[p]; c != nil {
		return c
	}
	location, _ := url.Parse(d.URI)
	var value any = d.Root
	for _, part := range strings.Split(strings.TrimPrefix(p, "/"), "/") {
		seen := map[string]bool{}
		for str(obj(value)["$ref"]) != "" {
			ref, err := url.Parse(str(obj(value)["$ref"]))
			if err != nil {
				return nil
			}
			location = location.ResolveReference(ref)
			if seen[location.String()] {
				return nil
			}
			seen[location.String()] = true
			resource := *location
			resource.Fragment = ""
			value, err = pointer(d.Registry.resources[resource.String()], location.Fragment)
			if err != nil {
				return nil
			}
		}
		var err error
		value, err = pointer(value, "/"+part)
		if err != nil {
			return nil
		}
		location.Fragment += "/" + part
	}
	c, _ := d.Compiler.Compile(location.String())
	return c
}

func (d *Document) model() Object {
	ops := d.operations()
	schemas := []any{}
	schemaURLs := Object{}
	for _, name := range keys(obj(obj(d.Root["components"])["schemas"])) {
		p := "/components/schemas/" + esc(name)
		u := route(d.Source) + "schemas/" + slug(name) + "/"
		schemaURLs["#"+p] = u
		schemas = append(schemas, Object{"name": name, "pointer": p, "url": u, "schema": obj(obj(d.Root["components"])["schemas"])[name]})
	}
	for _, op := range ops {
		for _, raw := range arr(op["variants"]) {
			v := obj(raw)
			examples := []any{}
			schemaKey := "schema"
			if _, ok := v["itemSchema"]; ok {
				schemaKey = "itemSchema"
			}
			schema := d.presentationSchema(str(v["pointer"]) + "/" + schemaKey)
			for _, rawEx := range arr(v["examples"]) {
				ex := obj(rawEx)
				if schema != nil && schema.Validate(ex["value"]) != nil {
					continue
				}
				ex["text"], ex["language"] = exampleText(str(v["media"]), ex["value"])
				examples = append(examples, ex)
			}
			v["examples"] = examples
		}
	}
	for _, op := range ops {
		for _, raw := range arr(op["parameters"]) {
			parameter := obj(raw)
			value, hasExample := parameterValue(parameter)
			if schema := d.presentationSchema(str(parameter["pointer"]) + "/schema"); hasExample && schema != nil && schema.Validate(value) != nil {
				parameter["exampleUsable"] = false
			}
		}
		op["url"] = route(d.Source) + "operations/" + slug(str(op["id"])) + "/"
		op["securitySchemes"] = obj(d.Root["components"])["securitySchemes"]
		op["acceptMedia"] = responseMedia(op)
		// HEAD negotiates the GET representation but has no response body.
		if op["method"] == "HEAD" && op["acceptMedia"] == "" {
			for _, get := range ops {
				if get["method"] == "GET" && get["path"] == op["path"] {
					op["acceptMedia"] = responseMedia(get)
					break
				}
			}
		}
		op["curl"], op["curlNotes"] = curlExample(d.Source, op)
		op["references"] = refs(op["raw"], schemaURLs)
		op["requestSchema"] = firstRequestSchema(op)
	}
	sourceURL := d.Source.SourceURL
	if sourceURL == "" {
		sourceURL = "/" + strings.TrimPrefix(d.Source.Source, "content/")
	}
	tags := []any{}
	seenTags := map[string]bool{}
	addTag := func(raw any) {
		tag := maps.Clone(obj(raw))
		name := str(tag["name"])
		if seenTags[name] {
			return
		}
		seenTags[name] = true
		if str(tag["kind"]) == "" {
			tag["kind"] = "nav"
		}
		if str(tag["summary"]) == "" {
			tag["summary"] = tag["x-displayName"]
			if str(tag["summary"]) == "" {
				tag["summary"] = tag["name"]
			}
		}
		tags = append(tags, tag)
	}
	for _, raw := range arr(d.Root["tags"]) {
		addTag(raw)
	}
	for _, op := range ops {
		for _, name := range arr(op["tags"]) {
			addTag(Object{"name": name})
		}
	}
	return Object{"experimental": d.Source.Experimental, "id": d.Source.ID, "product": d.Source.Product, "title": d.Source.Title, "version": obj(d.Root["info"])["version"], "description": obj(d.Root["info"])["description"], "url": route(d.Source), "manual": d.Source.Manual, "guides": d.Source.Guides, "connection": d.Source.Connection, "servers": d.Root["servers"], "securitySchemes": obj(d.Root["components"])["securitySchemes"], "tags": tags, "operations": ops, "schemas": schemas, "schemaURLs": schemaURLs, "digest": d.Digest, "owner": d.Source.Owner, "source": d.Source.Source, "sourceURL": sourceURL, "authoredOverview": d.Source.AuthoredOverview, "schemaCount": d.SchemaCount, "exampleCount": d.ExampleCount}
}
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
	if usable, ok := p["exampleUsable"].(bool); ok && !usable {
		return nil, false
	}
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

func exampleText(media string, value any) (string, string) {
	if s, ok := value.(string); ok && !strings.Contains(media, "json") {
		return s, "text"
	}
	return strings.TrimSpace(string(encoded(value))), "json"
}

func responseMedia(op Object) string {
	for _, raw := range arr(op["variants"]) {
		v := obj(raw)
		if v["direction"] == "Response" && strings.HasPrefix(str(v["status"]), "2") && str(v["media"]) != "" {
			return str(v["media"])
		}
	}
	return ""
}

func curlExample(src Source, op Object) (string, []string) {
	if len(obj(obj(op["raw"])["x-websocket"])) > 0 {
		return "", []string{"Use a WebSocket client with the subprotocol, credential transport, and frame schemas documented for this operation."}
	}
	if _, ok := obj(obj(op["raw"])["responses"])["101"]; ok {
		return "", []string{"This operation upgrades the HTTP connection. Use a client that supports the protocol described for this operation."}
	}
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
	if media := str(op["acceptMedia"]); media != "" {
		args = append(args, "--header "+shell("Accept: "+media))
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
		if len(examples) > 0 && strings.Contains(media, "json") {
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

// Examples are annotations: retrieve them without flattening schema constraints.
func (d *Document) mediaExamples(v Object) []any {
	out := d.exampleList(v)
	if len(out) > 0 {
		return out
	}
	s := v["schema"]
	seen := map[string]bool{}
	base, _ := url.Parse(d.URI)
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
		if ref == "" {
			break
		}
		u, err := url.Parse(ref)
		if err != nil {
			break
		}
		u = base.ResolveReference(u)
		if seen[u.String()] {
			break
		}
		seen[u.String()] = true
		fragment := u.Fragment
		u.Fragment = ""
		// Follow one reference at a time so intermediate schema annotations survive.
		s, err = pointer(d.Registry.resources[u.String()], fragment)
		if err != nil {
			break
		}
		base = u
	}
	return out
}
