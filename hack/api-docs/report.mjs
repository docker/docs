import fs from "node:fs";
import path from "node:path";
import { isDeepStrictEqual } from "node:util";
import YAML from "yaml";
const root = path.resolve(import.meta.dirname, "../..");
const dir = path.join(root, "prototypes/api-docs");
const out = path.join(root, "tmp/api-prototype");
const read = (p) => JSON.parse(fs.readFileSync(p));
const model = read(path.join(out, "data/api-prototype.json"));
const validation = read(path.join(out, "validation.json"));
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
function operations(s) {
  return Object.entries(s.paths).flatMap(([p, item]) =>
    methods
      .filter((m) => item[m])
      .map((m) => ({ path: p, method: m.toUpperCase(), raw: item[m] })),
  );
}
function resolve(s, v) {
  const seen = new Set();
  while (v?.$ref) {
    if (seen.has(v.$ref)) throw Error("Reference object cycle");
    seen.add(v.$ref);
    if (!v.$ref.startsWith("#/"))
      throw Error("Unexpected external reference in the publication corpus");
    v = v.$ref
      .slice(2)
      .split("/")
      .map((x) => x.replaceAll("~1", "/").replaceAll("~0", "~"))
      .reduce((o, k) => o[k], s);
  }
  return v;
}
const rows = [];
for (const api of model.apis) {
  const spec = YAML.parse(
    fs.readFileSync(path.join(dir, `converted/${api.id}.yaml`), "utf8"),
  );
  const original = YAML.parse(
    fs.readFileSync(path.join(dir, `original/${api.id}.yaml`), "utf8"),
  );
  const ops = operations(spec),
    sourceOps = operations(original);
  if (
    !isDeepStrictEqual(
      sourceOps.map((o) => o.method + " " + o.path).sort(),
      ops.map((o) => o.method + " " + o.path).sort(),
    )
  )
    throw Error(`Changed operation set: ${api.id}`);
  if (ops.length !== api.operations.length)
    throw Error(`Missing operations: ${api.id}`);
  let variants = 0,
    headers = 0,
    references = 0;
  for (const op of ops) {
    const prepared = api.operations.find(
      (o) => o.path === op.path && o.method === op.method,
    );
    if (!prepared || !isDeepStrictEqual(prepared.raw, op.raw))
      throw Error(`Operation source loss: ${api.id} ${op.path}`);
    const expected = [];
    const add = (direction, status, value) => {
      const v = resolve(spec, value);
      for (const [media, body] of Object.entries(v.content ?? { "": {} }))
        expected.push({
          direction,
          status,
          media,
          body,
          headers: v.headers ?? null,
        });
    };
    if (op.raw.requestBody) add("Request", "", op.raw.requestBody);
    for (const [status, v] of Object.entries(op.raw.responses ?? {}))
      add("Response", status, v);
    if (expected.length !== prepared.variants.length)
      throw Error(`Variant count mismatch: ${api.id} ${op.path}`);
    for (const v of expected) {
      const actual = prepared.variants.find(
        (x) =>
          x.direction === v.direction &&
          x.status === v.status &&
          x.media === v.media,
      );
      if (!actual) throw Error("Missing media variant");
      for (const k of ["schema", "itemSchema", "encoding"])
        if (!isDeepStrictEqual(v.body[k], actual[k]))
          throw Error(`Lost ${k}: ${api.id} ${op.path}`);
      if (!isDeepStrictEqual(v.headers, actual.headers ?? null))
        throw Error("Lost response headers");
      headers += Object.keys(v.headers ?? {}).length;
      variants++;
    }
    references += (JSON.stringify(op.raw).match(/"\$ref":/g) ?? []).length;
  }
  for (const [name, schema] of Object.entries(spec.components?.schemas ?? {})) {
    const actual = api.schemas.find((s) => s.name === name);
    if (!actual || !isDeepStrictEqual(schema, actual.schema))
      throw Error(`Lost schema: ${api.id} ${name}`);
  }
  const ledger = read(path.join(dir, `migrations/${api.id}.json`));
  const classes = {};
  for (const c of ledger.changes)
    classes[c.classification] = (classes[c.classification] ?? 0) + 1;
  const diagnostics = validation.find((r) => r.api === api.id);
  rows.push({
    api: api.id,
    operations: ops.length,
    variants,
    headers,
    references,
    namedSchemas: api.schemas.length,
    schemaRoots: diagnostics.schemas,
    examples: diagnostics.examples,
    exceptions: diagnostics.diagnostics?.length ?? 0,
    changes: ledger.changes.length,
    classifications: classes,
    coverage: "pass",
  });
}
const engineA = model.apis.find((a) => a.id === "engine-1.55"),
  engineB = model.apis.find((a) => a.id === "engine-1.56");
const effective = [];
for (const op of engineB.operations) {
  const prior = engineA.operations.find((x) => x.id === op.id);
  if (!prior) {
    effective.push({ operation: op.id, change: "added" });
    continue;
  }
  for (const field of ["security", "servers", "parameters"])
    if (!isDeepStrictEqual(prior[field], op[field]))
      effective.push({
        operation: op.id,
        field,
        before: prior[field],
        after: op[field],
      });
}
fs.mkdirSync(path.join(out, "reports"), { recursive: true });
fs.writeFileSync(
  path.join(out, "reports/coverage.json"),
  JSON.stringify(
    {
      modelVersion: model.modelVersion,
      apis: rows,
      engineEffectiveChanges: effective,
    },
    null,
    2,
  ) + "\n",
);
const md = [
  "# Prototype migration report",
  "",
  "Generated from checked-in source snapshots, converted specifications, and the presentation model.",
  "",
  "| API | Operations | Media/body variants | Schemas | Examples checked | Exceptions | Ledger changes |",
  "| --- | ---: | ---: | ---: | ---: | ---: | ---: |",
  ...rows.map(
    (r) =>
      `| ${r.api} | ${r.operations} | ${r.variants} | ${r.namedSchemas} | ${r.examples} | ${r.exceptions} | ${r.changes} |`,
  ),
  "",
  "Every operation, response variant, header definition, and named schema matches the converted source. Original and converted method/path sets match. See the source ledger for intentional contract corrections.",
];
fs.writeFileSync(path.join(out, "reports/summary.md"), md.join("\n") + "\n");
console.log("Source-to-model coverage passed for all six specifications.");
