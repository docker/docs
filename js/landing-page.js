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

var prefersDark = window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches;
var selectedNightTheme = readCookie("night");

if (selectedNightTheme === "true" || (selectedNightTheme === null && prefersDark)) {
  applyNight();
  $("#switch-style").prop("checked", true);
} else {
  applyDay();
  $("#switch-style").prop("checked", false);
}

function applyNight() {
  $("body").addClass("night");
}

function applyDay() {
  $("body").removeClass("night");
}

$("#switch-style").change(function () {
  if ($(this).is(":checked")) {
      applyNight();
      createCookie("night", true, 999)
  } else {
      applyDay();
      createCookie("night", false, 999);
  }
});
