import { createCSSVariablesTheme, registerCustomTheme } from "@pierre/diffs";

// One set of syntax roles for every language; Trident supplies both palettes.
const colors = {
  foreground: "var(--tri-foreground)",
  background: "var(--tri-muted)",
  "token-keyword": "var(--tri-primary)",
  "token-function": "var(--tri-primary)",
  "token-link": "var(--tri-primary)",
  "token-comment": "var(--tri-muted-foreground)",
  "token-string": "var(--tri-foreground)",
  "token-string-expression": "var(--tri-foreground)",
  "token-constant": "var(--tri-foreground)",
  "token-parameter": "var(--tri-foreground)",
  "token-punctuation": "var(--tri-foreground)",
  "token-inserted": "var(--tri-alert-success-title)",
  "token-deleted": "var(--tri-alert-error-title)",
  "token-changed": "var(--tri-primary)",
};

export const themes = ["docs-code-light", "docs-code-dark"];
for (const [index, name] of themes.entries()) {
  const theme = createCSSVariablesTheme({
    name,
    variablePrefix: "--docs-syntax-",
    variableDefaults: colors,
    fontStyle: false,
  });
  theme.type = index === 0 ? "light" : "dark";
  registerCustomTheme(name, () => Promise.resolve(theme));
}
