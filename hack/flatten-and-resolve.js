#!/usr/bin/env node

/**
 * Resolve Markdown links relative to each original file, then flatten index.md.
 * Only link destinations change; examples and other text stay byte-identical.
 * Usage: node hack/flatten-and-resolve.js [public-dir]
 */
const fs = require("node:fs");
const path = require("node:path");

const publicDir = path.resolve(process.argv[2] || "public");

function markdownFiles(dir) {
  return fs.readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    const file = path.join(dir, entry.name);
    return entry.isDirectory()
      ? markdownFiles(file)
      : file.endsWith(".md")
        ? [file]
        : [];
  });
}

function resolveLink(link, file) {
  // External URLs, protocol-relative URLs, and same-page anchors are complete.
  if (/^(?:[a-z][a-z0-9+.-]*:|\/\/|#)/i.test(link)) return link;
  const [, pathname, suffix] = link.match(/^([^?#]*)(.*)$/s);
  if (!pathname) return link;
  let url = pathname.startsWith("/")
    ? pathname
    : "/" +
      path
        .relative(publicDir, path.resolve(path.dirname(file), pathname))
        .split(path.sep)
        .join("/");
  url = url.replace(/^\/manuals\//, "/");
  url = url.replace(/\/_?index\.md$/, "/").replace(/\.md$/, "/");
  return url + suffix;
}

// Micromark provides exact destination spans, excluding labels, titles, and code.
// Apply edits backwards so offsets stay valid and formatting stays untouched.
function rewriteLinks(content, events, file) {
  const destinations = events
    .filter(
      ([event, token]) =>
        event === "enter" &&
        ["resourceDestinationString", "definitionDestinationString"].includes(
          token.type,
        ),
    )
    .map(([, token]) => token)
    .sort((a, b) => b.start.offset - a.start.offset);
  for (const token of destinations) {
    const start = token.start.offset;
    const end = token.end.offset;
    const destination = content.slice(start, end);
    content =
      content.slice(0, start) +
      resolveLink(destination, file) +
      content.slice(end);
  }
  return content;
}

async function main() {
  const { parse, preprocess, postprocess } = await import("micromark");
  const files = markdownFiles(publicDir);
  let rewritten = 0;
  let flattened = 0;
  for (const file of files) {
    const original = fs.readFileSync(file, "utf8");
    const events = postprocess(
      parse()
        .document()
        .write(preprocess()(original, "utf8", true)),
    );
    const content = rewriteLinks(original, events, file);
    if (content !== original) {
      fs.writeFileSync(file, content);
      rewritten++;
    }
    if (
      path.basename(file) === "index.md" &&
      file !== path.join(publicDir, "index.md")
    ) {
      fs.renameSync(file, path.dirname(file) + ".md");
      flattened++;
    }
  }
  console.log(
    `Markdown: rewrote links in ${rewritten} files; flattened ${flattened} index.md files`,
  );
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
