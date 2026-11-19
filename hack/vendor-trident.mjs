#!/usr/bin/env node
// Refresh from a reviewed Trident checkout; ordinary builds never access private packages.
import { createHash } from "node:crypto";
import { execFileSync } from "node:child_process";
import { readFile, writeFile, mkdir, copyFile, access } from "node:fs/promises";
import { resolve, join } from "node:path";
import { fileURLToPath } from "node:url";
const [sourcePath, lucidePath] = process.argv.slice(2);
if (!sourcePath || !lucidePath)
  throw new Error(
    "Usage: node hack/vendor-trident.mjs <built-trident-checkout> <lucide-static-package>",
  );
const repo = fileURLToPath(new URL("../", import.meta.url));
const source = resolve(sourcePath);
const revision = execFileSync("git", ["-C", source, "rev-parse", "HEAD"], {
  encoding: "utf8",
}).trim();
const sha = (bytes) => createHash("sha256").update(bytes).digest("hex");
for (const kind of ["tokens", "icons"]) {
  const target = join(
    repo,
    kind === "tokens" ? "assets/css/vendor/trident" : "assets/icons/trident",
  );
  const previous = JSON.parse(await readFile(join(target, "manifest.json")));
  const pkgPath = join(source, "packages", `trident-${kind}`);
  const pkg = JSON.parse(await readFile(join(pkgPath, "package.json")));
  if (pkg.name !== `@docker/trident-${kind}`)
    throw new Error("Unexpected package");
  const files = {};
  const inputs = new Map();
  for (const name of Object.keys(previous.files)) {
    let input = join(pkgPath, kind === "tokens" ? "dist" : "src/svg", name);
    if (kind === "icons") {
      try {
        await access(input);
      } catch {
        input = join(resolve(lucidePath), "icons", name);
      }
    }
    const bytes = await readFile(input);
    inputs.set(name, bytes);
    files[name] = sha(bytes);
  }
  const manifest = {
    package: pkg.name,
    version: pkg.version,
    revision,
    repository: "https://github.com/docker/trident",
    files,
  };
  if (kind === "icons") {
    const lucide = JSON.parse(
      await readFile(join(resolve(lucidePath), "package.json")),
    );
    if (lucide.version !== pkg.devDependencies["lucide-static"])
      throw new Error("Lucide version does not match Trident");
    manifest.lucide = lucide.version;
  }
  await mkdir(target, { recursive: true });
  for (const [name, bytes] of inputs)
    await writeFile(join(target, name), bytes);
  if (kind === "icons")
    await copyFile(
      join(pkgPath, "LICENSE.lucide"),
      join(target, "LICENSE.lucide"),
    );
  await writeFile(
    join(target, "manifest.json"),
    JSON.stringify(manifest, null, 2) + "\n",
  );
  console.log(`Vendored ${pkg.name}@${pkg.version} at ${revision}`);
}
