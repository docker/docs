(function (d) {
    "use strict";
    for (const h of d.querySelectorAll("H1, H2, H3")) {
        if (h.id != null && h.id.length > 0) {
            h.insertAdjacentHTML('beforeend', `<a href="#${h.id}" class="anchorLink">🔗</a>`)
        }
    }
})(document);
