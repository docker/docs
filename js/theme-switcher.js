// Replacement for jQuery $(selector) and $.onready()
const _ = s => document.querySelector(s);
const ready = f => document.readyState !== 'loading' ? f() : document.addEventListener('DOMContentLoaded', f)

const darkMode = window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches
const selectedTheme = window.localStorage ? localStorage.getItem("theme") : null;

if (selectedTheme !== null) {
    if (selectedTheme === "night") _("body").classList.add("night");
} else if (darkMode) {
    _("body").classList.add("night");
}

function themeToggler() {
    const sw = _("#switch-style"), b = _("body");
    if (sw && b) {
        sw.checked = b.classList.contains("night")
        sw.addEventListener("change", function (){
            b.classList.toggle("night", this.checked)
            if (window.localStorage) {
                this.checked ? localStorage.setItem("theme", "night") : localStorage.setItem("theme", "day")
            }
        })
    }
}

ready(themeToggler)
