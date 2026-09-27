import { File, preloadHighlighter } from "@pierre/diffs";
import { bundledLanguages } from "./pierre/shiki.js";

import { themes } from "./pierre/themes.js";
const aliases = {
  console: "docs-console",
  bash: "shellscript",
  sh: "shellscript",
  shell: "shellscript",
  dockerfile: "docker",
  yml: "yaml",
  js: "javascript",
  ts: "typescript",
  py: "python",
  golang: "go",
  ps: "powershell",
  plaintext: "text",
  txt: "text",
  systemd: "ini",
  env: "dotenv",
  cs: "csharp",
};
const languageFor = (name = "text") => {
  const language = aliases[name.toLowerCase()] || name.toLowerCase();
  return Object.hasOwn(bundledLanguages, language) ? language : "text";
};
const themeType = () =>
  document.documentElement.classList.contains("dark") ? "dark" : "light";
const instances = new Map();
const pending = new WeakSet();
const loads = new Map();
const instructions = new Set([
  "ADD",
  "ARG",
  "CMD",
  "COPY",
  "ENTRYPOINT",
  "ENV",
  "EXPOSE",
  "FROM",
  "HEALTHCHECK",
  "LABEL",
  "ONBUILD",
  "RUN",
  "SHELL",
  "STOPSIGNAL",
  "USER",
  "VOLUME",
  "WORKDIR",
]);

function prepare(language) {
  if (!loads.has(language)) {
    loads.set(language, preloadHighlighter({ themes, langs: [language] }));
  }
  return loads.get(language);
}

function highlightLines(value) {
  const ranges = String(Array.isArray(value) ? value.join(" ") : value || "")
    .split(/[\s,]+/)
    .map((part) => part.match(/^(\d+)(?:-(\d+))?$/))
    .filter(Boolean)
    .map((match) => [Number(match[1]), Number(match[2] || match[1])]);
  return (number) =>
    ranges.some(([start, end]) => number >= start && number <= end);
}

function enhance(wrapper) {
  if (pending.has(wrapper) || instances.has(wrapper)) return;
  const source = wrapper.querySelector("[data-pierre-source]");
  if (!source) return;
  pending.add(wrapper);
  const language = languageFor(wrapper.dataset.pierreLanguage);
  let options = {};
  try {
    options = JSON.parse(wrapper.dataset.pierreOptions || "{}");
  } catch {
    /* Use defaults. */
  }
  const highlighted = highlightLines(options.hl_lines);
  const start = Math.max(1, Number.parseInt(options.linenostart, 10) || 1);
  const host = document.createElement("diffs-container");
  host.className = "pierre-render";
  host.style.visibility = "hidden";
  wrapper.append(host);
  const fail = (error) => {
    console.warn("Code highlighting failed; keeping plain code.", error);
    instances.get(wrapper)?.cleanUp();
    instances.delete(wrapper);
    host.remove();
    source.hidden = false;
    wrapper.classList.add("pierre-fallback");
  };
  prepare(language)
    .then(() => {
      if (!wrapper.isConnected) return;
      const instance = new File({
        theme: { light: themes[0], dark: themes[1] },
        themeType: themeType(),
        disableFileHeader: true,
        disableLineNumbers: !["true", "table", "inline"].includes(
          String(options.linenos),
        ),
        disableErrorHandling: true,
        overflow: "scroll",
        unsafeCSS: `
        :host { --diffs-bg: var(--tri-muted); --diffs-font-family: var(--tri-font-family-mono);
          --diffs-font-size: var(--tri-scale-mono-code-font-size);
          --diffs-line-height: var(--tri-scale-mono-code-line-height);
          --diffs-gap-block: 0px; --diffs-gap-inline: 8px; }
        [data-docs-highlight] { background: var(--tri-accent-subtle); }
        a { color: inherit; text-decoration: underline dashed; text-underline-offset: 3px; }
      `,
        onPostRender(container) {
          const root = container.shadowRoot;
          if (!root?.querySelector("[data-line]")) return;
          for (const line of root.querySelectorAll("[data-line]")) {
            line.toggleAttribute(
              "data-docs-highlight",
              highlighted(Number(line.dataset.line)),
            );
          }
          for (const number of root.querySelectorAll("[data-column-number]")) {
            const label = number.querySelector("[data-line-number-content]");
            if (label)
              label.textContent =
                Number(number.dataset.columnNumber) + start - 1;
          }
          if (language === "docker") {
            for (const span of root.querySelectorAll("[data-line] span")) {
              if (!instructions.has(span.textContent) || span.closest("a"))
                continue;
              const link = document.createElement("a");
              link.href = `/reference/dockerfile/#${span.textContent.toLowerCase()}`;
              link.title = `Learn more about the ${span.textContent} instruction`;
              span.replaceWith(link);
              link.append(span);
            }
          }
          source.hidden = true;
          host.style.visibility = "";
          wrapper.classList.add("pierre-ready");
        },
      });
      instances.set(wrapper, instance);
      instance.render({
        file: { name: "code", lang: language, contents: source.textContent },
        fileContainer: host,
      });
    })
    .catch(fail);
}

function scan(root) {
  if (root.matches?.(".pierre-code")) enhance(root);
  root.querySelectorAll(".pierre-code").forEach(enhance);
  for (const [wrapper, instance] of instances) {
    if (!wrapper.isConnected) {
      instance.cleanUp();
      instances.delete(wrapper);
    }
  }
}

// Warm common Gordon languages immediately, alongside the languages on this page.
for (const language of [
  "shellscript",
  "docker",
  "yaml",
  "json",
  "javascript",
  "python",
  "go",
]) {
  prepare(language).catch(() => {}); // Each block handles its own fallback.
}
scan(document);
const chat = document.querySelector("#gordon-chat");
if (chat) {
  new MutationObserver(() => scan(chat)).observe(chat, {
    childList: true,
    subtree: true,
  });
}
new MutationObserver(() => {
  for (const instance of instances.values()) instance.setThemeType(themeType());
}).observe(document.documentElement, {
  attributes: true,
  attributeFilter: ["class"],
});
