// The site postprocessor rewrites URL-like strings inside fenced JSON too.
// Keep generated reference Markdown intact; its links already target published URLs.
import fs from "node:fs";
import path from "node:path";
import { execFileSync } from "node:child_process";
const root = path.resolve(import.meta.dirname, "../..");
const site = path.join(root, "tmp/api-prototype/site");
const data = JSON.parse(
  fs.readFileSync(path.join(root, "tmp/api-prototype/data/api-prototype.json")),
);
const urls = [
  "/reference/api/",
  ...data.apis.flatMap((api) => [
    api.url,
    ...api.operations.map((op) => op.url),
    ...api.schemas.map((schema) => schema.url),
  ]),
];
const saved = urls
  .map((url) => path.join(site, url, "index.md"))
  .map((p) => [p, fs.readFileSync(p)]);
for (const [p] of saved) fs.unlinkSync(p);
try {
  execFileSync(
    process.execPath,
    [path.join(root, "hack/flatten-and-resolve.js"), site],
    { stdio: "inherit" },
  );
} finally {
  for (const [p, b] of saved) {
    const dest = path.basename(p) === "index.md" ? path.dirname(p) + ".md" : p;
    fs.writeFileSync(dest, b);
  }
}
