// Replacement for jQuery $(selector) and $.onready()
const _ = s => document.querySelector(s);
const ready = f => document.readyState !== 'loading' ? f() : document.addEventListener('DOMContentLoaded', f)

function getJSON(url, fn) {
    const xhr = new XMLHttpRequest();
    xhr.open('GET', url, true);
    xhr.responseType = 'json';
    xhr.onload = () => xhr.status === 200 ? fn(xhr.response) : null;
    xhr.send();
}

// throttle / debounce events. taken from https://programmingwithmosh.com/javascript/javascript-throttle-and-debounce-patterns/
function debounce(fn, msec) {
    let id;
    return function(...args) {
        clearTimeout(id); id = setTimeout(() => fn.apply(this, args), msec);
    }
}

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
