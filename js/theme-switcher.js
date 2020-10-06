// Cookie functions
function createCookie(name, value, days) {
    var expires = "";
    if (days) {
        var date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + value + expires + "; path=/";
}

function readCookie(name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(";");
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) === " ") c = c.substring(1, c.length);
        if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
}

var rootDom = $("html");
var themeCookieName = "night";
var selectedNightTheme = readCookie(themeCookieName);
var switchStyle = $("#switch-style");

function applyTheme(name) {
    rootDom.removeClass(function(_, className) {
        return (className.match(/theme-\S+/g) || []).join(' ');
    });
    rootDom.addClass("theme-" + name);
}

if (selectedNightTheme === "true") {
    applyTheme("dark");
    switchStyle.prop("checked", true);
} else if (selectedNightTheme === "false") {
    applyTheme("light");
    switchStyle.prop("checked", false);
} else {
    var prefersDark = window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches;
    switchStyle.prop("checked", prefersDark);
}

switchStyle.change(function () {
    if ($(this).is(":checked")) {
        applyTheme("dark");
        createCookie(themeCookieName, true, 999)
    } else {
        applyTheme("light");
        createCookie(themeCookieName, false, 999);
    }
});
