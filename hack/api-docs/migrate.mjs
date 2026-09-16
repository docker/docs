// Migration-only tooling. The publication compiler is Go; original sources are immutable.
import fs from "node:fs";
import path from "node:path";
import crypto from "node:crypto";
import { execFileSync } from "node:child_process";
import YAML from "yaml";
import { isDeepStrictEqual } from "node:util";
import converter from "swagger2openapi";
const root = path.resolve(import.meta.dirname, "../..");
const dir = path.join(root, "prototypes/api-docs");
const catalog = JSON.parse(fs.readFileSync(path.join(dir, "catalog.json")));
const digest = (b) => crypto.createHash("sha256").update(b).digest("hex");
const read = (p) => fs.readFileSync(path.join(dir, p), "utf8");
const json = (v) => JSON.stringify(v, null, 2) + "\n";
const clone = (v) => structuredClone(v);
const esc = (s) => s.replaceAll("~", "~0").replaceAll("/", "~1");
const own = (o, k) => Object.hasOwn(o, k);
const methods = [
  "get",
  "put",
  "post",
  "delete",
  "options",
  "head",
  "patch",
  "trace",
  "query",
];
function parse(s) {
  const d = YAML.parseDocument(s, { uniqueKeys: true });
  if (d.errors.length) throw d.errors[0];
  return d.toJS();
}
function output(file, value, verify) {
  const b = typeof value === "string" ? value : json(value);
  const p = path.join(dir, file);
  if (verify) {
    if (!fs.existsSync(p) || fs.readFileSync(p, "utf8") !== b)
      throw Error(`Reproduction differs: ${file}`);
  } else fs.writeFileSync(p, b);
}
function diff(a, b, p = "") {
  if (JSON.stringify(a) === JSON.stringify(b)) return [];
  if (
    a &&
    b &&
    !Array.isArray(a) &&
    !Array.isArray(b) &&
    typeof a === "object" &&
    typeof b === "object"
  )
    return [...new Set([...Object.keys(a), ...Object.keys(b)])]
      .sort()
      .flatMap((k) => diff(a[k], b[k], p + "/" + esc(k)));
  return [
    {
      pointer: p,
      beforePresent: a !== undefined,
      afterPresent: b !== undefined,
      ...(a === undefined ? {} : { before: a }),
      ...(b === undefined ? {} : { after: b }),
    },
  ];
}
function apply(tree, changes) {
  for (const c of changes) {
    const parts = c.pointer
      .slice(1)
      .split("/")
      .map((s) => s.replaceAll("~1", "/").replaceAll("~0", "~"));
    let o = tree;
    for (const p of parts.slice(0, -1)) o = o[p];
    const k = parts.at(-1);
    if (c.afterPresent) o[k] = clone(c.after);
    else delete o[k];
  }
  return tree;
}
function walk(v, fn, p = "") {
  if (!v || typeof v !== "object") return;
  fn(v, p);
  for (const [k, x] of Object.entries(v)) walk(x, fn, p + "/" + esc(k));
}
function operations(spec) {
  const out = [];
  for (const [p, item] of Object.entries(spec.paths ?? {}))
    for (const m of methods) if (item[m]) out.push([p, m, item[m], item]);
  return out;
}
const idPath = path.join(dir, "operation-ids.json");
let ids = fs.existsSync(idPath) ? JSON.parse(fs.readFileSync(idPath)) : {};
const command = process.argv[2] ?? "verify";
const verify = command === "verify";
if (command === "snapshot") {
  const revision = execFileSync("git", ["rev-parse", "HEAD"], {
    cwd: root,
    encoding: "utf8",
  }).trim();
  const lock = {
    revision,
    tools: {
      swagger2openapi: "7.0.8",
      yaml: "2.8.1",
      libopenapi: "0.38.7",
      vacuum: "0.30.3",
      jsonschema: "6.0.3",
      oasdiff: "1.31.0",
    },
    sources: [],
  };
  for (const api of catalog.apis) {
    const b = fs.readFileSync(path.join(root, api.source));
    const f = `original/${api.id}.yaml`;
    if (
      fs.existsSync(path.join(dir, f)) &&
      !fs.readFileSync(path.join(dir, f)).equals(b)
    )
      throw Error(`Refusing to overwrite snapshot: ${f}`);
    output(f, b.toString(), false);
    lock.sources.push({ ...api, sha256: digest(b) });
  }
  output("source-lock.json", lock, false);
  console.log("Captured six immutable source snapshots.");
  process.exit(0);
}
const lock = JSON.parse(read("source-lock.json"));
for (const api of catalog.apis) {
  const source = lock.sources.find((s) => s.id === api.id);
  const originalText = read(`original/${api.id}.yaml`);
  if (
    digest(originalText) !== source.sha256 ||
    digest(fs.readFileSync(path.join(root, api.source))) !== source.sha256
  )
    throw Error(`Source digest mismatch: ${api.id}`);
  const original = parse(originalText);
  let spec = clone(original);
  const ledger = [];
  const stage = (name, classification, rationale, evidence, fn) => {
    const before = clone(spec);
    fn();
    for (const c of diff(before, spec))
      ledger.push({
        ...c,
        id: `${api.id}-${digest(name + c.pointer).slice(0, 12)}`,
        stage: name,
        classification,
        rationale,
        evidence,
        owner: api.owner,
        sourcePointer: c.beforePresent ? c.pointer : null,
        destinationPointer: c.afterPresent ? c.pointer : null,
      });
  };
  if (spec.swagger) {
    const result = await converter.convertObj(clone(spec), {
      refSiblings: "allOf",
      patch: false,
      warnOnly: false,
    });
    stage(
      "swagger-to-oas3",
      "mechanical conversion",
      "Pinned converter with refSiblings=allOf; compare every moved field and review Engine generation semantics.",
      "swagger2openapi 7.0.8",
      () => {
        spec = result.openapi;
      },
    );
    output(`migrations/${api.id}.intermediate.json`, spec, verify);
  }
  stage(
    "oas32-header",
    "mechanical conversion",
    "Adopt the common OpenAPI version and its selected schema dialect.",
    "API-DOCUMENTATION-ARCHITECTURE.md S1",
    () => {
      spec.openapi = "3.2.0";
      spec.jsonSchemaDialect = catalog.dialect;
    },
  );
  // Schema migration is scoped to schema positions, excluding examples and application property names.
  function schemas(v, fn, p = "") {
    if (!v || typeof v !== "object") return;
    for (const [k, x] of Object.entries(v)) {
      const q = p + "/" + esc(k);
      if (k === "schemas" && p === "/components") {
        for (const [n, s] of Object.entries(x)) schema(s, fn, q + "/" + esc(n));
      } else if (k === "schema" || k === "itemSchema") schema(x, fn, q);
      else if (
        !["example", "examples", "value", "default", "enum"].includes(k) &&
        !k.startsWith("x-")
      )
        schemas(x, fn, q);
    }
  }
  function schema(s, fn, p) {
    if (!s || typeof s !== "object") return;
    fn(s, p);
    for (const [k, x] of Object.entries(s)) {
      if (
        [
          "properties",
          "patternProperties",
          "$defs",
          "definitions",
          "dependentSchemas",
        ].includes(k)
      )
        for (const [n, c] of Object.entries(x ?? {}))
          schema(c, fn, p + "/" + k + "/" + esc(n));
      else if (
        ["allOf", "anyOf", "oneOf", "prefixItems"].includes(k) &&
        Array.isArray(x)
      )
        x.forEach((c, i) => schema(c, fn, p + "/" + k + "/" + i));
      else if (
        [
          "items",
          "additionalProperties",
          "unevaluatedProperties",
          "unevaluatedItems",
          "contains",
          "not",
          "if",
          "then",
          "else",
          "contentSchema",
          "propertyNames",
        ].includes(k)
      )
        schema(x, fn, p + "/" + k);
    }
  }
  stage(
    "oas30-schema-semantics",
    "mechanical conversion",
    "Translate OpenAPI 3.0 nullable and exclusive bounds into JSON Schema constraints; preserve sibling assertions.",
    "OpenAPI 3.0 Schema Object and JSON Schema 2020-12",
    () =>
      schemas(spec, (s) => {
        if (own(s, "nullable")) {
          const nullable = s.nullable;
          delete s.nullable;
          if (nullable && s.type) {
            const type = Array.isArray(s.type) ? s.type : [s.type];
            s.type = [...new Set([...type, "null"])];
          }
        }
        for (const [ex, bound] of [
          ["exclusiveMinimum", "minimum"],
          ["exclusiveMaximum", "maximum"],
        ])
          if (typeof s[ex] === "boolean") {
            if (s[ex] && typeof s[bound] === "number") {
              s[ex] = s[bound];
              delete s[bound];
            } else delete s[ex];
          }
      }),
  );
  if (api.product === "engine")
    stage(
      "engine-nullability-assumption",
      "provisional assumption",
      "Interpret x-nullable as admitting null for this preview only. Moby must reconcile pointer/omission/generation semantics before upstream adoption.",
      "Original x-nullable declarations; inventory Engine migration risk",
      () =>
        schemas(spec, (s) => {
          if (s["x-nullable"] === true) {
            delete s["x-nullable"];
            const base = clone(s);
            for (const k of Object.keys(s)) delete s[k];
            s.anyOf = [base, { type: "null" }];
          }
        }),
    );
  stage(
    "editorial-metadata",
    "editorial completion",
    "Supply conservative operation descriptions and stable IDs; declare primary navigation tags without changing endpoint behavior.",
    "Existing operation summaries, paths, and tags",
    () => {
      spec.info.description ??= `HTTP reference for ${api.title}.`;
      const tags = new Map((spec.tags ?? []).map((t) => [t.name, t]));
      for (const [p, m, op] of operations(spec)) {
        const key = `${api.id}:${m.toUpperCase()} ${p}`;
        if (!op.operationId) {
          if (!ids[key]) {
            if (verify) throw Error(`Missing recorded operation ID: ${key}`);
            ids[key] =
              m +
              p
                .replace(/\{([^}]+)\}/g, " by $1 ")
                .replace(/[^a-zA-Z0-9]+/g, " ")
                .trim()
                .split(" ")
                .map((s) => s[0].toUpperCase() + s.slice(1))
                .join("");
          }
          op.operationId = ids[key];
        }
        if (!op.summary?.trim()) op.summary = `${m.toUpperCase()} ${p}`;
        if (!op.description?.trim())
          op.description = op.summary.endsWith(".")
            ? op.summary
            : op.summary + ".";
        if (!op.tags?.length) op.tags = ["Operations"];
        for (const name of op.tags)
          if (!tags.has(name))
            tags.set(name, { name, description: `${name} operations.` });
      }
      for (const t of tags.values()) {
        t.summary ??= t["x-displayName"] ?? t.name;
        t.kind = operations(spec).some((x) => x[2].tags?.[0] === t.name)
          ? "nav"
          : "info";
        t.description ??= `${t.summary} reference.`;
        delete t["x-displayName"];
      }
      spec.tags = [...tags.values()];
    },
  );
  stage(
    "property-required-correction",
    "evidence-backed correction",
    "Move misplaced boolean property required flags to the parent required array; preserve the declared requiredness.",
    "Original property-level required: true declarations",
    () =>
      schemas(spec, (s) => {
        for (const [name, prop] of Object.entries(s.properties ?? {})) {
          if (prop && typeof prop.required === "boolean") {
            if (prop.required)
              s.required = [...new Set([...(s.required ?? []), name])];
            delete prop.required;
          }
        }
      }),
  );
  stage(
    "portable-descriptions",
    "editorial completion",
    "Remove renderer-only badges/HTML presentation, retaining their text and Markdown links.",
    "Docker profile S12",
    () =>
      walk(spec, (v) => {
        for (const k of ["description", "summary"])
          if (typeof v[k] === "string")
            v[k] = v[k]
              .replace(
                /<SchemaDefinition schemaRef="#\/components\/schemas\/([^"]+)"\s*\/>/g,
                "[$1](#schema-$1)",
              )
              .replace(/\[<a href="([^"]+)">([^<]+)<\/a>\]/g, "[$2]($1)")
              .replace(/<a href="([^"]+)">([^<]+)<\/a>/g, "[$2]($1)")
              .replace(/<br\s*\/?\s*>/gi, "\n")
              .replace(/<\/?p>/gi, "\n")
              .replace(
                /<\/?(?:span|div|small|strong|em|b|i)(?:\s[^>]*)?>/gi,
                "",
              );
      }),
  );
  stage(
    "empty-description-completion",
    "editorial completion",
    "Replace descriptions that contained only presentation markup with the existing operation summary.",
    "Existing operation summary",
    () => {
      for (const [p, m, op] of operations(spec))
        if (!op.description?.trim()) op.description = op.summary + ".";
    },
  );
  if (api.product === "engine")
    stage(
      "engine-response-media",
      "evidence-backed correction",
      "Use JSON for documented ErrorResponse bodies and restore the archive success media type lost by Swagger conversion.",
      "Engine API introductory error contract and GET /containers/{id}/archive description",
      () => {
        for (const [p, m, op] of operations(spec))
          for (const [status, response] of Object.entries(op.responses ?? {})) {
            const content = response.content ?? {};
            const entries = Object.values(content);
            if (
              Number(status) >= 400 &&
              entries.length &&
              entries.every(
                (v) => v.schema?.$ref === "#/components/schemas/ErrorResponse",
              )
            )
              response.content = { "application/json": entries[0] };
            if (
              p === "/containers/{id}/archive" &&
              m === "get" &&
              status === "200"
            )
              response.content = { "application/x-tar": {} };
          }
      },
    );
  if (api.id === "dvp")
    stage(
      "dvp-structural-correction",
      "evidence-backed correction",
      "Use the standard http bearer scheme and move invalid path security to each operation without changing the documented requirement.",
      "Existing HubAuth bearer definition and path-level security declarations",
      () => {
        spec.components.securitySchemes.HubAuth.type = "http";
        for (const item of Object.values(spec.paths)) {
          if (own(item, "security")) {
            for (const m of methods)
              if (item[m] && !own(item[m], "security"))
                item[m].security = clone(item.security);
            delete item.security;
          }
        }
        if (own(spec, "features.openapi")) {
          spec["x-features-openapi"] = spec["features.openapi"];
          delete spec["features.openapi"];
        }
      },
    );
  if (api.id === "registry")
    stage(
      "registry-version-auth-assumption",
      "provisional assumption",
      "Use protocol version 2 as provisional info.version and declare challenge-acquired bearer access. Confirm anonymous/public access alternatives and version identity with Hub owners.",
      "Registry description and content/reference/api/registry/auth.md",
      () => {
        spec.info.version = "2";
        spec.components ??= {};
        spec.components.securitySchemes ??= {};
        spec.components.securitySchemes.registryToken = {
          type: "http",
          scheme: "bearer",
          description: api.auth,
        };
        spec.security = [{ registryToken: [] }];
      },
    );
  if (api.id === "hub") {
    stage(
      "hub-team-schema-assumption",
      "provisional assumption",
      "Replace the missing response reference with the existing repository schema, the likely base of a team repository. Product owner must verify returned fields.",
      "components.schemas.team_repo and existing repository response definition",
      () => {
        const s = spec.components.schemas.team_repo;
        if (s?.allOf?.[0]?.$ref === "#/components/responses/team_repo")
          s.allOf[0].$ref = "#/components/schemas/repository_info";
      },
    );
    stage(
      "scim-auth-context",
      "evidence-backed correction",
      "Declare the separate provisioning token for SCIM operations instead of inheriting Hub credential exchange guidance.",
      "content/manuals/security/provisioning/scim/provision-scim.md",
      () => {
        spec.components.securitySchemes.scimToken = {
          type: "http",
          scheme: "bearer",
          description:
            "Use the SCIM provisioning token configured for the organization.",
        };
        for (const [p, m, op] of operations(spec))
          if (p.startsWith("/v2/scim/")) op.security = [{ scimToken: [] }];
      },
    );
  }
  if (api.id === "governance")
    stage(
      "governance-token-prose",
      "evidence-backed correction",
      "Align copied token request field names with the Hub token operation; upstream owner must confirm credential support.",
      "Hub POST /v2/auth/token identifier and secret fields",
      () => {
        const s = spec.components.securitySchemes.bearerAuth;
        if (s?.description)
          s.description = s.description
            .replace(/`password`/g, "`secret`")
            .replace(/`username`/g, "`identifier`");
      },
    );
  stage(
    "head-response-bodies",
    "evidence-backed correction",
    "HEAD responses transfer headers without a response body. Preserve headers and status codes; remove declared content from HEAD responses.",
    "RFC 9110 section 9.3.2",
    () => {
      for (const [p, m, op] of operations(spec))
        if (m === "head")
          for (const response of Object.values(op.responses ?? {}))
            if (response.content) delete response.content;
    },
  );
  stage(
    "explicit-root-security",
    "editorial completion",
    "Make the existing absence of inherited HTTP authentication explicit. Local socket permissions remain connection metadata.",
    "OpenAPI root security inheritance; existing source operation policies",
    () => {
      if (!own(spec, "security")) spec.security = [];
    },
  );
  // Explicit hand-reviewed corrections are applied as JSON-pointer patches with their own rationale.
  const patchFile = path.join(dir, "corrections.json");
  const corrections = fs.existsSync(patchFile)
    ? JSON.parse(fs.readFileSync(patchFile))
    : [];
  for (const c of corrections.filter((x) => x.api === api.id))
    stage(c.id, c.classification, c.rationale, c.evidence, () =>
      apply(spec, [c]),
    );
  const converted = YAML.stringify(spec, { lineWidth: 0 });
  const replay = apply(clone(original), ledger);
  if (!isDeepStrictEqual(replay, spec))
    throw Error(`Ledger replay failed: ${api.id}`);
  output(`converted/${api.id}.yaml`, converted, verify);
  const report = {
    api: api.id,
    source: source,
    convertedSha256: digest(converted),
    profile: catalog.profileVersion,
    changes: ledger,
  };
  output(`migrations/${api.id}.json`, report, verify);
  let textual;
  try {
    textual = execFileSync(
      "diff",
      [
        "-u",
        "--label",
        `original/${api.id}.yaml`,
        "--label",
        `converted/${api.id}.yaml`,
        path.join(dir, `original/${api.id}.yaml`),
        "-",
      ],
      { input: converted, encoding: "utf8" },
    );
  } catch (e) {
    if (e.status !== 1) throw e;
    textual = e.stdout;
  }
  output(`migrations/${api.id}.patch`, textual, verify);
  console.log(
    `${api.id}: ${operations(spec).length} operations, ${ledger.length} recorded changes`,
  );
}
output("operation-ids.json", ids, verify);
