import fs from "node:fs";
import path from "node:path";
import { marked } from "marked";

export function resolveRedirect(url, redirects) {
  const seen = new Set();
  while (redirects[url.pathname] && !seen.has(url.pathname)) {
    seen.add(url.pathname);
    const next = new URL(redirects[url.pathname], url);
    if (!next.hash) next.hash = url.hash;
    url = next;
  }
  return url;
}
export const decode = (s) =>
  s
    .replace(/&#(\d+);/g, (_, n) => String.fromCodePoint(Number(n)))
    .replaceAll("&amp;", "&")
    .replaceAll("&#34;", '"')
    .replaceAll("&#39;", "'")
    .replaceAll("&lt;", "<")
    .replaceAll("&gt;", ">");
// Hugo's production minifier removes optional attribute quotes.
export function attributes(html, name) {
  const pattern = new RegExp(
    `\\s${name}=(?:"([^"]*)"|'([^']*)'|([^\\s>]+))`,
    "g",
  );
  return [...html.matchAll(pattern)].map((m) => decode(m[1] ?? m[2] ?? m[3]));
}

export function localLinkProblems(
  base,
  file,
  hrefs,
  redirects = {},
  {
    checkFragments = true,
    label = "API",
    hosts = ["localhost", "docs.docker.com"],
  } = {},
) {
  const problems = [];
  for (const href of hrefs) {
    let url;
    try {
      url = new URL(href, "http://localhost:1314/" + path.relative(base, file));
    } catch {
      continue;
    }
    if (
      !hosts.includes(url.hostname) ||
      !["http:", "https:"].includes(url.protocol)
    )
      continue;
    url = resolveRedirect(url, redirects);
    if (!hosts.includes(url.hostname)) continue;
    let target = path.join(
      base,
      decodeURIComponent(url.pathname),
      url.pathname.endsWith("/") ? "index.html" : "",
    );
    if (fs.existsSync(target) && fs.statSync(target).isDirectory())
      target = path.join(target, "index.html");
    if (!fs.existsSync(target)) {
      problems.push(
        `Broken ${label} link: ${path.relative(base, file)} -> ${href}`,
      );
    } else if (checkFragments && url.hash && target.endsWith(".html")) {
      const targetHTML = fs.readFileSync(target, "utf8");
      const anchors = [
        ...attributes(targetHTML, "id"),
        ...attributes(targetHTML, "name"),
      ];
      if (!anchors.includes(decodeURIComponent(url.hash.slice(1))))
        problems.push(
          `Broken ${label} anchor: ${path.relative(base, file)} -> ${href}`,
        );
    }
  }
  return problems;
}

export function markdownLinks(markdown) {
  return attributes(marked.parse(markdown), "href");
}

export function engineOutputProblems(base, sourceDir, apis) {
  const problems = [];
  const versions = fs
    .readdirSync(sourceDir)
    .filter((file) => /^v\d+\.\d+\.md$/.test(file))
    .map((file) => file.slice(0, -3));
  for (const version of versions) {
    const url = `/reference/api/engine/version/${version}/`;
    const file = path.join(base, url, "index.html");
    if (!fs.existsSync(file)) {
      problems.push(`Missing Engine page: ${url}`);
      continue;
    }
    const html = fs.readFileSync(file, "utf8");
    const primarySpec = url.slice(0, -1) + ".yaml";
    const generated = apis.some((api) => api.url === url);
    if (generated) {
      if (
        !attributes(html, "data-api-view").includes("overview") ||
        html.includes("<redoc")
      )
        problems.push(`Generated Engine renderer missing: ${url}`);
    } else if (
      !attributes(html, "spec-url").includes(primarySpec) ||
      !html.includes("redoc/v2.5.1/bundles/redoc.standalone.js") ||
      attributes(html, "class").some((value) =>
        value.split(/\s+/).includes("docs-site"),
      ) ||
      html.includes("data-api-view=")
    )
      problems.push(`Legacy Engine renderer changed: ${url}`);
    if (!fs.existsSync(path.join(base, primarySpec)))
      problems.push(`Missing primary Engine specification: ${primarySpec}`);
    for (const specURL of attributes(html, "spec-url")) {
      if (!fs.existsSync(path.join(base, specURL)))
        problems.push(`Missing Engine specification: ${url} -> ${specURL}`);
    }
  }
  return problems;
}
