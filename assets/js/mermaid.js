import mermaid from "mermaid";

const SELECTOR = "pre.mermaid";

const isDark = () => document.documentElement.classList.contains("dark");

// Rendering replaces the element content with an SVG, so the diagram source
// is captured once per element and reused for every redraw.
const sources = new WeakMap();

const getSource = (el) => {
  let source = sources.get(el);
  if (source === undefined) {
    // Same treatment mermaid.run() gives element content: decode HTML
    // entities, but keep <br/> tags because labels use them for line breaks.
    const decoder = document.createElement("textarea");
    decoder.innerHTML = el.innerHTML;
    source = decoder.value.replace(/<br\s*\/?>/gi, "<br/>").trim();
    sources.set(el, source);
  }
  return source;
};

// Incremented on every render request. A render that is no longer the latest
// stops before writing to the page, so a slow render cannot overwrite a
// newer theme.
let generation = 0;
let nextId = 0;
let renderedDark;

const render = async () => {
  const current = ++generation;
  const dark = isDark();
  renderedDark = dark;

  const elements = Array.from(document.querySelectorAll(SELECTOR));
  if (elements.length === 0) return;

  mermaid.initialize({
    startOnLoad: false,
    securityLevel: "strict",
    theme: dark ? "dark" : "default",
    // On a syntax error, render() would otherwise leave Mermaid's error
    // graphic in a temporary element at the end of <body>. The failing
    // element shows its source instead (see the catch below).
    suppressErrorRendering: true,
  });

  for (const el of elements) {
    const source = getSource(el);
    try {
      // mermaid.render() returns the SVG as a string instead of writing it
      // into the element, so an existing diagram stays on screen until its
      // replacement is ready.
      const { svg, bindFunctions } = await mermaid.render(
        `mermaid-${nextId++}`,
        source,
      );
      if (current !== generation) return;
      el.innerHTML = svg;
      bindFunctions?.(el);
      delete el.dataset.error;
    } catch (error) {
      if (current !== generation) return;
      // Unhide the element so the source and the console error can be
      // matched up.
      el.dataset.error = "";
      console.error("Mermaid failed to render a diagram", error);
    }
  }
};

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", render);
} else {
  render();
}

// theme.js and the header theme button both set the light/dark class on
// <html>. Redraw with the matching Mermaid theme whenever it changes.
new MutationObserver(() => {
  if (isDark() !== renderedDark) render();
}).observe(document.documentElement, {
  attributes: true,
  attributeFilter: ["class"],
});
