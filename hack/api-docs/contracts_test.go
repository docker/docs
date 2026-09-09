package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func localDocument(t *testing.T, api string) *Document {
	t.Helper()
	d, err := loadDocument(filepath.Join("..", "..", "content", "reference", "api", api, "latest.yaml"), "validation/dialects")
	if err != nil {
		t.Fatal(err)
	}
	d.Source = Source{ID: api}
	d.validate("")
	if len(d.Diagnostics) != 0 {
		t.Fatalf("%s must pass strict validation: %+v", api, d.Diagnostics)
	}
	return d
}

func TestHubResponseContracts(t *testing.T) {
	d := localDocument(t, "hub")
	for _, tc := range []struct {
		pointer, value string
		valid          bool
	}{
		{"/components/schemas/error", `{"message":"not found","errinfo":null}`, true},
		{"/components/schemas/error", `{"errinfo":{"field":["invalid"],"limit":5}}`, true},
		{"/components/schemas/error", `{"errinfo":[]}`, false},
		{"/components/schemas/getAccessTokensResponse", `{"next":null,"previous":null,"results":[]}`, true},
		{"/components/schemas/getOrgAccessTokensResponse", `{"next":null,"previous":null,"results":[]}`, true},
		{"/components/schemas/org_member_paginated", `{"count":0,"next":null,"previous":null,"results":[]}`, true},
		{"/components/schemas/org_member_paginated", `[]`, false},
		{"/components/schemas/bulk_invite", `{"invitees":[]}`, true},
		{"/components/schemas/bulk_invite", `{"invitees":{"invitees":[]}}`, false},
		{"/components/schemas/tag", `{"v2":true,"images":[{"architecture":"amd64","variant":null,"features":null,"os_features":null,"os_version":null}]}`, true},
		{"/components/schemas/tag", `{"v2":"true","images":{}}`, false},
		{"/components/schemas/immutable_tags_verify_request", `{"regex":"v.*"}`, true},
		{"/components/schemas/immutable_tags_verify_request", `{"regex":"v1,v2"}`, false},
		{"/components/schemas/scim_service_provider_config", `{"authenticationSchemes":[{"type":"oauthbearertoken"}]}`, true},
		{"/components/schemas/scim_service_provider_config", `{"authenticationSchemes":{}}`, false},
		{"/components/requestBodies/scim_update_user_request/content/application~1scim+json/schema", `{"schemas":["urn:ietf:params:scim:schemas:core:2.0:User"],"userName":"user@example.com","active":false}`, true},
		{"/components/requestBodies/scim_update_user_request/content/application~1scim+json/schema", `{"schemas":["urn:ietf:params:scim:schemas:core:2.0:User"]}`, false},
		{"/paths/~1v2~1orgs~1{org_name}~1members~1export/get/responses/200/content/text~1csv/schema", `"Name,Username\nUser,dockeruser\n"`, true},
		{"/paths/~1v2~1orgs~1{org_name}~1members~1export/get/responses/200/content/text~1csv/schema", `[{"Name":"User"}]`, false},
	} {
		t.Run(tc.pointer+"/"+tc.value, func(t *testing.T) {
			s, err := d.Compiler.Compile(d.URI + "#" + tc.pointer)
			if err != nil {
				t.Fatal(err)
			}
			if err := s.Validate(jsonNumber(tc.value)); (err == nil) != tc.valid {
				t.Fatalf("valid = %t, want %t: %v", err == nil, tc.valid, err)
			}
		})
	}
	for _, response := range []string{"scim_get_resource_types_resp", "scim_get_schemas_resp", "scim_get_users_resp"} {
		v, err := pointer(d.Root, "/components/responses/"+response+"/content/application~1scim+json/schema/properties")
		if err != nil {
			t.Fatal(err)
		}
		if obj(v)["Resources"] == nil || obj(v)["resources"] != nil {
			t.Fatalf("incorrect SCIM list field casing: %s", response)
		}
	}
}

func TestRegistryRequests(t *testing.T) {
	d := localDocument(t, "registry")
	for _, raw := range arr(d.model()["operations"]) {
		op := obj(raw)
		curl := str(op["curl"])
		if strings.Count(curl, "Authorization:") != 1 || !strings.Contains(curl, "${REGISTRY_TOKEN}") {
			t.Errorf("incorrect registry authentication: %s", curl)
		}
		if strings.Count(curl, "Content-Type:") > 1 {
			t.Errorf("duplicate content type: %s", curl)
		}
		if (op["id"] == "GetImageManifest" || op["id"] == "HeadImageManifest") && !strings.Contains(curl, "Accept: application/vnd.docker.distribution.manifest.v2+json") {
			t.Errorf("missing manifest negotiation: %s", curl)
		}
		if op["id"] == "GetBlobUploadStatus" && !strings.Contains(curl, "--request GET") {
			t.Errorf("incorrect upload status method: %s", curl)
		}
	}
}

func TestReservedHeaderParameters(t *testing.T) {
	for _, name := range []string{"Authorization", "accept", "CONTENT-TYPE", "Content-Range"} {
		t.Run(name, func(t *testing.T) {
			d := fixture(t, "valid.yaml")
			op := obj(obj(obj(d.Root["paths"])["/public"])["get"])
			op["parameters"] = []any{Object{"name": name, "in": "header", "description": "Header value.", "schema": Object{"type": "string"}}}
			d.validate("")
			found := false
			for _, diag := range d.Diagnostics {
				found = found || (diag.Rule == "S8" && strings.Contains(diag.Message, "ignores"))
			}
			if found != (name != "Content-Range") {
				t.Fatalf("unexpected header validation: %+v", d.Diagnostics)
			}
		})
	}
}

func TestExampleText(t *testing.T) {
	for _, tc := range []struct {
		media          string
		value          any
		text, language string
	}{
		{"text/csv", "Name,Username\nUser,dockeruser\n", "Name,Username\nUser,dockeruser\n", "text"},
		{"text/plain", "<html>&\n", "<html>&\n", "text"},
		{"application/json", "a\nb", `"a\nb"`, "json"},
		{"application/scim+json", false, "false", "json"},
		{"application/json", jsonNumber("0"), "0", "json"},
	} {
		text, language := exampleText(tc.media, tc.value)
		if text != tc.text || language != tc.language {
			t.Errorf("got %q (%s), want %q (%s)", text, language, tc.text, tc.language)
		}
	}
}
