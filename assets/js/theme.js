// return 'light' or 'dark' depending on localStorage (pref) or system setting
function getThemePreference() {
  const theme = localStorage.getItem("theme-preference");
  if (theme) return theme;
  else
    return window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light";
}

// update root class based on os setting or localstorage
const preference = getThemePreference();
document.firstElementChild.className = preference === "dark" ? "dark" : "light";
localStorage.setItem("theme-preference", preference);

// set innertext for the theme switch button
// window.addEventListener("DOMContentLoaded", () => {
//   const themeSwitchButton = document.querySelector("#theme-switch");
//   themeSwitchButton.textContent = `${preference}_mode`;
// });
