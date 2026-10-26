// Limit the eager Pierre bundle to languages used by Docker docs.
export * from "shiki/core";
export { createHighlighterCore as createHighlighter } from "shiki/core";
export { createJavaScriptRegexEngine } from "shiki/engine/javascript";
export { createOnigurumaEngine } from "shiki/engine/oniguruma";
export const bundledThemes = {};
export const bundledLanguages = {
  "docs-console": () => import("./console.js"),
  shellscript: () => import("shiki/langs/shellscript.mjs"),
  docker: () => import("shiki/langs/docker.mjs"),
  yaml: () => import("shiki/langs/yaml.mjs"),
  json: () => import("shiki/langs/json.mjs"),
  javascript: () => import("shiki/langs/javascript.mjs"),
  typescript: () => import("shiki/langs/typescript.mjs"),
  tsx: () => import("shiki/langs/tsx.mjs"),
  python: () => import("shiki/langs/python.mjs"),
  go: () => import("shiki/langs/go.mjs"),
  java: () => import("shiki/langs/java.mjs"),
  hcl: () => import("shiki/langs/hcl.mjs"),
  powershell: () => import("shiki/langs/powershell.mjs"),
  xml: () => import("shiki/langs/xml.mjs"),
  diff: () => import("shiki/langs/diff.mjs"),
  ini: () => import("shiki/langs/ini.mjs"),
  sql: () => import("shiki/langs/sql.mjs"),
  properties: () => import("shiki/langs/properties.mjs"),
  toml: () => import("shiki/langs/toml.mjs"),
  html: () => import("shiki/langs/html.mjs"),
  csharp: () => import("shiki/langs/csharp.mjs"),
  nginx: () => import("shiki/langs/nginx.mjs"),
  http: () => import("shiki/langs/http.mjs"),
  groovy: () => import("shiki/langs/groovy.mjs"),
  markdown: () => import("shiki/langs/markdown.mjs"),
  graphql: () => import("shiki/langs/graphql.mjs"),
  c: () => import("shiki/langs/c.mjs"),
  cpp: () => import("shiki/langs/cpp.mjs"),
  ruby: () => import("shiki/langs/ruby.mjs"),
  php: () => import("shiki/langs/php.mjs"),
  makefile: () => import("shiki/langs/makefile.mjs"),
  dotenv: () => import("shiki/langs/dotenv.mjs"),
  "git-commit": () => import("shiki/langs/git-commit.mjs"),
};
