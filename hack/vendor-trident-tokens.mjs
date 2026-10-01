#!/usr/bin/env node
// Import generated CSS from an extracted package or a built Trident checkout.
import { createHash } from "node:crypto";
import { mkdir, readFile, readdir, writeFile } from "node:fs/promises";
import { resolve } from "node:path";
import { fileURLToPath } from "node:url";

const [packagePath, revision] = process.argv.slice(2);
if (!packagePath || !/^[a-f0-9]{40}$/.test(revision || "")) {
  throw new Error(
    "Usage: node hack/vendor-trident-tokens.mjs <package-directory> <source-commit>",
  );
}
const source = resolve(packagePath);
const pkg = JSON.parse(await readFile(resolve(source, "package.json"), "utf8"));
if (pkg.name !== "@docker/trident-tokens") {
  throw new Error("Expected @docker/trident-tokens");
}
const styles = resolve(source, "dist/styles");
const names = (await readdir(styles))
  .filter((name) => /^tri-[a-z-]+\.css$/.test(name))
  .sort();
const files = new Map(
  await Promise.all(
    names.map(async (name) => [name, await readFile(resolve(styles, name))]),
  ),
);
if (!files.has("tri-tokens.css"))
  throw new Error("Missing generated CSS entry point");
for (const content of files.values()) {
  for (const match of content
    .toString()
    .matchAll(/@import\s+['"]\.\/([^'"]+)['"]/g)) {
    if (!files.has(match[1]))
      throw new Error(`Missing imported file: ${match[1]}`);
  }
}
const target = fileURLToPath(
  new URL("../assets/css/vendor/trident/", import.meta.url),
);
await mkdir(target, { recursive: true });
const hashes = {};
for (const [name, content] of files) {
  await writeFile(resolve(target, name), content);
  hashes[name] = createHash("sha256").update(content).digest("hex");
}
await writeFile(
  resolve(target, "manifest.json"),
  JSON.stringify(
    {
      package: pkg.name,
      version: pkg.version,
      repository: "https://github.com/docker/trident",
      revision,
      license: pkg.license || null,
      files: hashes,
    },
    null,
    2,
  ) + "\n",
);
console.log(`Vendored ${files.size} CSS files from ${pkg.name}@${pkg.version}`);
