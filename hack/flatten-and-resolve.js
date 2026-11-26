#!/usr/bin/env node

/**
 * Apply Hugo-resolved Markdown destinations, then flatten index.md.
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

function resolveLink(token, file, destinations) {
  const { destination, kind, definition, attribute } = token;
  if (/^(?:[a-z][a-z0-9+.-]*:|\/\/|#|\?)/i.test(destination))
    return destination;
  // HTML from shortcodes is already relative to the published HTML page.
  if (attribute) {
    const page = path
      .relative(publicDir, file)
      .replace(/(^|\/)index\.md$/, "$1")
      .replace(/\.md$/, "");
    return (
      new URL(
        destination,
        `https://docs.docker.com/${page.replace(/\/$/, "")}${page ? "/" : ""}`,
      ).pathname + (destination.match(/[?#].*$/s)?.[0] || "")
    );
  }
  const resolved =
    destinations[kind]?.[destination] ??
    (definition ? destinations.images?.[destination] : undefined);
  if (resolved !== undefined) return resolved;
  // Unused reference definitions are not rendered by Hugo.
  if (definition) return destination;
  // Generated layouts already emit published absolute paths.
  if (destination.startsWith("/") && !/\.md(?:[?#]|$)/.test(destination))
    return destination;
  throw new Error(
    `No Hugo resolution for ${path.relative(publicDir, file)} -> ${destination}`,
  );
}

function rewriteLinks(content, tokens, file, destinations) {
  for (const token of tokens) {
    const resolved = resolveLink(token, file, destinations);
    if (resolved === token.destination) continue;
    const replacement = token.attribute
      ? resolved
          .replaceAll("&", "&amp;")
          .replaceAll('"', "&quot;")
          .replaceAll("'", "&#39;")
      : resolved;
    content =
      content.slice(0, token.start) + replacement + content.slice(token.end);
  }
  return content;
}

async function main() {
  const { linkDestinations } = await import("./markdown-links.mjs");
  const metadataDir = path.join(publicDir, "__markdown-links");
  const files = markdownFiles(publicDir);
  const changes = [];
  let rewritten = 0;
  let flattened = 0;
  for (const file of files) {
    const original = fs.readFileSync(file, "utf8");
    const metadataFile = path.join(
      metadataDir,
      path.relative(publicDir, file) + ".json",
    );
    const destinations = fs.existsSync(metadataFile)
      ? JSON.parse(fs.readFileSync(metadataFile, "utf8"))
      : {};
    const content = rewriteLinks(
      original,
      linkDestinations(original),
      file,
      destinations,
    );
    changes.push({ file, original, content });
  }
  for (const { file, original, content } of changes) {
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
  fs.rmSync(metadataDir, { recursive: true, force: true });
  console.log(
    `Markdown: rewrote links in ${rewritten} files; flattened ${flattened} index.md files`,
  );
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
