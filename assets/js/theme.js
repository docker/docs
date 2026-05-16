// update root class based on os setting or localstorage
const storedTheme = localStorage.getItem("theme-preference");
const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
document.firstElementChild.className =
  storedTheme === "dark" || storedTheme === "light"
    ? storedTheme
    : prefersDark
      ? "dark"
      : "light";
document.firstElementChild.dataset.themePreference = storedTheme || "auto";
