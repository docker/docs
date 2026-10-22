import Alpine from "alpinejs";
import collapse from "@alpinejs/collapse";
import persist from "@alpinejs/persist";
import focus from "@alpinejs/focus";
import { marked } from "marked";
window.Alpine = Alpine;

Alpine.plugin(collapse);
Alpine.plugin(persist);
Alpine.plugin(focus);

// Configure marked to escape HTML in text tokens only (not code blocks)
marked.use({
  walkTokens(token) {
    // Escape HTML in text and HTML tokens, preserve code blocks
    if (token.type === "text" || token.type === "html") {
      const text = token.text || token.raw;
      const escaped = text
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;");
      if (token.text) token.text = escaped;
      if (token.raw) token.raw = escaped;
    }
  },
});

// Add $markdown magic for rendering markdown with syntax highlighting
Alpine.magic("markdown", () => {
  return (content) => {
    if (!content) return "";
    const html = marked(content);

    const div = document.createElement("div");
    div.innerHTML = html;
    div.querySelectorAll("pre > code").forEach((code) => {
      const pre = code.parentElement;
      pre.classList.add("not-prose");
      pre.setAttribute("data-pierre-source", "");
      const wrapper = document.createElement("div");
      wrapper.className = "pierre-code not-prose";
      wrapper.dataset.pierreLanguage =
        [...code.classList]
          .find((name) => name.startsWith("language-"))
          ?.slice(9) || "text";
      pre.replaceWith(wrapper);
      wrapper.append(pre);
    });

    // Handle inline code elements (not in pre blocks)
    div.querySelectorAll("code:not(pre code)").forEach((code) => {
      code.classList.add("not-prose");
    });

    return div.innerHTML;
  };
});

// Stores
Alpine.store("showSidebar", false);
Alpine.store("gordon", {
  isOpen: false,
  query: "",
  autoSubmit: false,
  toggle() {
    this.isOpen = !this.isOpen;
  },
  open(query, autoSubmit = false) {
    this.isOpen = true;
    if (query) {
      this.query = query;
      this.autoSubmit = autoSubmit;
    }
  },
  close() {
    this.isOpen = false;
  },
});

Alpine.start();
