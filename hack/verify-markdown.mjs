import fs from "node:fs";
import path from "node:path";
import { pathToFileURL } from "node:url";
import { marked } from "marked";
import { parseFragment } from "parse5";
import {
  localLinkProblems,
  resolveRedirect,
} from "./api-docs/output-checks.mjs";

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

export async function verifyMarkdown(base, { request = fetch } = {}) {
  const redirects = JSON.parse(
    fs.readFileSync(path.join(base, "redirects.json"), "utf8"),
  );
  const index = fs.readFileSync(path.join(base, "llms-full.txt"), "utf8");
  const siteOrigin = new URL(
    index.match(/^Markdown:\s*(\S+)/m)?.[1] || "https://docs.docker.com",
  );
  const hosts = [...new Set([siteOrigin.hostname, "docs.docker.com"])];
  const problems = [];
  const remote = new Map();
  const files = markdownFiles(base);
  for (const file of files) {
    const html = marked.parse(fs.readFileSync(file, "utf8"));
    const destinations = [];
    const visit = (node) => {
      for (const attr of node.attrs || []) {
        if (["href", "src"].includes(attr.name)) destinations.push(attr.value);
      }
      for (const child of node.childNodes || []) visit(child);
    };
    visit(parseFragment(html));
    const local = [];
    for (const href of destinations) {
      let url;
      try {
        url = new URL(
          href,
          new URL(path.relative(base, file), siteOrigin.origin + "/"),
        );
      } catch {
        problems.push(
          `Invalid Markdown URL: ${path.relative(base, file)} -> ${href}`,
        );
        continue;
      }
      if (
        !hosts.includes(url.hostname) ||
        !["http:", "https:"].includes(url.protocol)
      )
        continue;
      url = resolveRedirect(url, redirects);
      // Learn is published by docker/learn and is absent from this build.
      if (hosts.includes(url.hostname) && url.pathname.startsWith("/learn/")) {
        url = new URL(url.pathname + url.search, "https://docs.docker.com");
        if (!remote.has(url.href)) remote.set(url.href, new Set());
        remote.get(url.href).add(path.relative(base, file));
      } else {
        local.push(url.href);
      }
    }
    problems.push(
      ...localLinkProblems(base, file, local, redirects, {
        checkFragments: false,
        label: "Markdown",
        hosts,
      }),
    );
  }
  for (const [, href] of index.matchAll(/^Markdown:\s*(\S+)\s*$/gm)) {
    const url = new URL(href);
    const file = path.join(base, decodeURIComponent(url.pathname));
    if (!fs.existsSync(file) || !fs.statSync(file).isFile())
      problems.push(`Missing advertised Markdown: ${href}`);
  }
  for (const [url, sources] of remote) {
    try {
      const response = await request(url, {
        signal: AbortSignal.timeout(15000),
      });
      await response.body?.cancel();
      if (!response.ok)
        problems.push(
          `Broken Markdown link: ${[...sources].join(", ")} -> ${url} (HTTP ${response.status})`,
        );
    } catch (error) {
      problems.push(`Unable to check Markdown link: ${url}: ${error.message}`);
    }
  }
  return { files: files.length, remote: remote.size, problems };
}

if (
  process.argv[1] &&
  import.meta.url === pathToFileURL(process.argv[1]).href
) {
  const result = await verifyMarkdown(
    path.resolve(process.argv[2] || "public"),
  );
  if (result.problems.length) {
    console.error(result.problems.join("\n"));
    process.exitCode = 1;
  } else {
    console.log(
      `Markdown: verified links and images in ${result.files} files, advertised URLs, and ${result.remote} Learn URLs`,
    );
  }
}
