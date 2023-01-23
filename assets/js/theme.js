// return true if dark mode, false if light
function getThemePreference() {
  const theme = localStorage.getItem("theme-preference")
  if (theme) return theme
  else
    return window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light"
}

// update root class based on os setting or localstorage
const preference = getThemePreference()
document.firstElementChild.className = preference === "dark" ? "dark" : "light"
localStorage.setItem("theme-preference", preference)
