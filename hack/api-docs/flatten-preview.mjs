// The site postprocessor rewrites URL-like strings inside fenced JSON too.
// Keep generated reference Markdown intact; its links already target published URLs.
import fs from "node:fs";
import path from "node:path";
import { execFileSync } from "node:child_process";
const root = path.resolve(import.meta.dirname, "../..");
const site = path.join(root, "tmp/api-prototype/site");
function files(dir) {
  return fs
    .readdirSync(dir, { withFileTypes: true })
    .flatMap((e) =>
      e.isDirectory()
        ? files(path.join(dir, e.name))
        : [path.join(dir, e.name)],
    );
}
const saved = files(path.join(site, "api-prototype"))
  .filter((p) => p.endsWith(".md"))
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
