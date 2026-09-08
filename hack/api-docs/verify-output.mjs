import fs from "node:fs";
import path from "node:path";
import { isDeepStrictEqual } from "node:util";
const root = path.resolve(import.meta.dirname, "../..");
const base = path.join(root, "tmp/api-prototype/site");
const data = JSON.parse(
  fs.readFileSync(path.join(root, "tmp/api-prototype/data/api-prototype.json")),
);
const decode = (s) =>
  s
    .replace(/&#(\d+);/g, (_, n) => String.fromCodePoint(Number(n)))
    .replaceAll("&amp;", "&")
    .replaceAll("&#34;", '"')
    .replaceAll("&#39;", "'")
    .replaceAll("&lt;", "<")
    .replaceAll("&gt;", ">");
let count = 0;
const problems = [];
function check(url) {
  const htmlPath = path.join(base, url, "index.html");
  const mdPath = path.join(base, url.replace(/\/$/, "")) + ".md";
  if (!fs.existsSync(htmlPath) || !fs.existsSync(mdPath)) {
    problems.push(`Missing HTML/Markdown: ${url}`);
    return ["", ""];
  }
  const html = fs.readFileSync(htmlPath, "utf8"),
    md = fs.readFileSync(mdPath, "utf8");
  count++;
  return [decode(html), md];
}
check("/api-prototype/");
for (const api of data.apis) {
  check(api.url);
  for (const op of api.operations) {
    const [html, md] = check(op.url);
    for (const variant of op.variants)
      if (
        !html.includes(`data-api-variant="${variant.pointer}"`) ||
        !md.includes(`### ${variant.direction} ${variant.status}`)
      )
        problems.push(`Missing variant: ${op.url} ${variant.pointer}`);
    for (const p of op.parameters)
      if (
        !html.includes(`data-api-parameter="${p.name}"`) ||
        !md.includes(`### ${p.name}`)
      )
        problems.push(`Missing parameter: ${op.url} ${p.name}`);
    if (!html.includes(op.path) || !md.includes(op.method + " " + op.path))
      problems.push(`Missing operation signature: ${op.url}`);
    if (!md.includes(op.curl))
      problems.push(`Request example mismatch: ${op.url}`);
  }
  for (const schema of api.schemas) {
    const [html, md] = check(schema.url);
    if (
      ![...md.matchAll(/```json\n([\s\S]*?)\n```/g)].some((m) => {
        try {
          return isDeepStrictEqual(JSON.parse(m[1]), schema.schema);
        } catch {
          return false;
        }
      })
    )
      problems.push(`Incomplete Markdown schema: ${schema.url}`);
    if (!html.includes(schema.name))
      problems.push(`Missing schema heading: ${schema.url}`);
  }
}
// Check preview-owned links in every generated preview page, including copied specification prose.
function files(dir) {
  return fs
    .readdirSync(dir, { withFileTypes: true })
    .flatMap((e) =>
      e.isDirectory()
        ? files(path.join(dir, e.name))
        : [path.join(dir, e.name)],
    );
}
for (const file of files(path.join(base, "api-prototype")).filter((p) =>
  p.endsWith(".html"),
)) {
  const html = fs.readFileSync(file, "utf8");
  for (const m of html.matchAll(/href="([^"]+)"/g)) {
    let url;
    try {
      url = new URL(
        decode(m[1]),
        "http://localhost:1314/" + path.relative(base, file),
      );
    } catch {
      continue;
    }
    if (
      url.hostname !== "localhost" ||
      !url.pathname.startsWith("/api-prototype/")
    )
      continue;
    const target = path.join(
      base,
      decodeURIComponent(url.pathname),
      url.pathname.endsWith("/") ? "index.html" : "",
    );
    if (!fs.existsSync(target))
      problems.push(
        `Broken preview link: ${path.relative(base, file)} -> ${m[1]}`,
      );
  }
}
if (problems.length)
  throw Error(
    problems.slice(0, 25).join("\n") + `\n${problems.length} output failures`,
  );
const bytes = files(path.join(base, "api-prototype")).reduce(
  (n, p) => n + fs.statSync(p).size,
  0,
);
fs.writeFileSync(
  path.join(root, "tmp/api-prototype/reports/output.json"),
  JSON.stringify(
    { htmlMarkdownPagePairs: count, previewBytes: bytes, result: "pass" },
    null,
    2,
  ) + "\n",
);
console.log(
  `Verified ${count} HTML/Markdown page pairs and preview-owned links.`,
);
