// Run before paint, independently of Alpine. Storage can be unavailable in private contexts.
(() => {
  const root = document.documentElement;
  const media = window.matchMedia("(prefers-color-scheme: dark)");
  let preference;
  try {
    preference = localStorage.getItem("theme-preference");
  } catch {}
  if (!["light", "dark"].includes(preference)) preference = null;
  function apply() {
    const dark = preference === "dark" || (preference == null && media.matches);
    root.classList.toggle("dark", dark);
    root.classList.toggle("light", !dark);
    root.dataset.themePreference = dark ? "dark" : "light";
    document
      .querySelector("[data-theme-toggle]")
      ?.setAttribute("aria-label", `Switch to ${dark ? "light" : "dark"} mode`);
  }
  apply();
  document.addEventListener("DOMContentLoaded", apply, { once: true });
  media.addEventListener("change", apply);
  document.addEventListener("click", (event) => {
    if (!event.target.closest("[data-theme-toggle]")) return;
    preference = root.classList.contains("dark") ? "light" : "dark";
    try {
      localStorage.setItem("theme-preference", preference);
    } catch {}
    apply();
  });
})();
