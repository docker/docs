
function navClicked(sourceLink) {
    let el = document.getElementById("#item"+sourceLink)
    if (el) {
        el.classList.toggle("in")
    }
}

let outputLetNav = [], totalTopics = 0;

function pageIsInSection(tree) {
    function processBranch(branch) {
        for (let k = 0; k < branch.length; k++) {
            if (branch[k].section) {
                processBranch(branch[k].section);
            } else {
                if (branch[k].path === pageURL && !branch[k].nosync) {
                    found = true;
                    break;
                }
            }
        }
    }

    let found = false;
    processBranch(tree);
    return found;
}

function walkTree(tree) {
    for (const page of tree) {
        totalTopics++;
        if (page.section) {
            let sectionHasPath = pageIsInSection(page.section);
            outputLetNav.push('<li><a onclick="navClicked(' + totalTopics + ')" data-target="#item' + totalTopics + '" data-toggle="collapse" data-parent="#stacked-menu"')
            if (sectionHasPath) {
                outputLetNav.push('aria-expanded="true"')
            } else {
                outputLetNav.push('class="collapsed" aria-expanded="false"')
            }
            outputLetNav.push(">" + page.sectiontitle + '<span class="caret arrow"></span></a>');
            outputLetNav.push('<ul class="nav collapse');
            if (sectionHasPath) outputLetNav.push(" in");
            outputLetNav.push('" id="#item' + totalTopics + '" aria-expanded="');
            if (sectionHasPath) {
                outputLetNav.push("true");
            } else {
                outputLetNav.push("false");
            }
            outputLetNav.push('">');
            walkTree(page.section);
            outputLetNav.push("</ul></li>");
        } else {
            // just a regular old topic; this is a leaf, not a branch; render a link!
            outputLetNav.push('<li><a href="' + page.path + '"')
            if (page.path === pageURL && !page.nosync) {
                outputLetNav.push('class="active currentPage"')
            }
            outputLetNav.push(">" + page.title + "</a></li>")
        }
    }
}

function renderNav() {
    getJSON( "/js/toc.json", function( data ) {
        for (const item of data.horizontalnav) {
            if (item.path === pageURL || pageIsInSection(data[item.node])) {
                // This is the current section. Set the corresponding header-nav link
                // to active, and build the left-hand (vertical) navigation
                _('#'+item.node).closest("li").classList.add("active")
                walkTree(data[item.node]);
                _("#jsTOCLeftNav").innerHTML = outputLetNav.join("");
            }
        }
        // Scroll the current menu item into view. We actually pick the item *above*
        // the current item to give some headroom above
        scrollMenuItem("#jsTOCLeftNav a.currentPage")
    });
}

// Scroll the given menu item into view. We actually pick the item *above*
// the current item to give some headroom above
function scrollMenuItem(selector) {
    let item = document.querySelector(selector)
    if (item) {
        item = item.closest("li")
    }
    if (item) {
        item = item.previousElementSibling
    }
    if (item) {
        item.scrollIntoView(true)
        if (window.location.hash.length < 2) {
            // Scrolling the side-navigation may scroll the whole page as well
            // this is a dirty hack to scroll the main content back to the top
            // if we're not on a specific anchor
            document.querySelector("main.col-content").scrollIntoView(true)
        }
    }
}

function highlightRightNav(heading) {
    $("#my_toc a.active").removeClass("active");

    if (heading !== "title") {
        $("#my_toc a[href='#" + heading + "']").addClass("active");
    }
}

var currentHeading = "";
$(window).scroll(function () {
    var headingPositions = [];
    $("h1, h2, h3, h4, h5, h6").each(function () {
        if (this.id === "") this.id = "title";
        headingPositions[this.id] = this.getBoundingClientRect().top;
    });
    headingPositions.sort();
    // the headings have all been grabbed and sorted in order of their scroll
    // position (from the top of the page). First one is toppermost.
    for (var key in headingPositions) {
        if (!headingPositions.hasOwnProperty(key)) {
            continue;
        }
        if (headingPositions[key] > 0 && headingPositions[key] < 200) {
            if (currentHeading !== key) {
                // a new heading has scrolled to within 200px of the top of the page.
                // highlight the right-nav entry and de-highlight the others.
                highlightRightNav(key);
                currentHeading = key;
            }
            break;
        }
    }
});

var navHeight = $(".navbar").outerHeight(true) + 80;

$(document.body).scrollspy({
    target: "#leftCol",
    offset: navHeight
});

function loadHash(hashObj) {
    // Using jQuery's animate() method to add smooth page scroll
    // The optional number (800) specifies the number of milliseconds it takes to scroll to the specified area
    $("html, body").animate({scrollTop: $(hashObj).offset().top - 80}, 800);
}

$(document).ready(function () {
    // Add smooth scrolling to all links
    $(".toc-nav a").on("click", function (event) {
        // Make sure this.hash has a value before overriding default behavior
        if (this.hash !== "") {
            // Prevent default anchor click behavior
            event.preventDefault();

            // Store hash
            var hash = this.hash;
            loadHash(hash);

            // Add hash (#) to URL when done scrolling (default click behavior)
            window.location.hash = hash;
        }
    });
    if (window.location.hash) loadHash(window.location.hash);
});


$(document).ready(function () {
    // Add smooth scrolling to all links
    $(".nav-sidebar ul li a").on("click", function (event) {

        // Make sure this.hash has a value before overriding default behavior
        if (this.hash !== "") {
            // Prevent default anchor click behavior
            event.preventDefault();

            // Store hash
            var hash = this.hash;

            // Using jQuery's animate() method to add smooth page scroll
            // The optional number (800) specifies the number of milliseconds it takes to scroll to the specified area
            $("html, body").animate({
                scrollTop: $(hash).offset().top - 80
            }, 800, function () {
                // Add hash (#) to URL when done scrolling (default click behavior)
                window.location.hash = hash;
            });
        }
    });
});

function initAcceptEULAModal() {
    $("main").on("click", "a.accept-eula", function (e) {
        e.preventDefault();
        _("#accept-eula .btn-primary").href = e.target.href;
        $('#accept-eula').modal('show')
    });
}

function initNavToggle() {
    $("#menu-toggle").click(function (e) {
        e.preventDefault();
        $(".wrapper").toggleClass("right-open");
        $(".col-toc").toggleClass("col-toc-hidden");
    });

    $(".navbar-toggle").click(function () {
        $("#sidebar-nav").each(function () {
            $(this).toggleClass("hidden-sm");
            $(this).toggleClass("hidden-xs");
        });
    });
}

ready(() => {
    renderNav()
    initNavToggle()
    initAcceptEULAModal()
    $('[data-toggle="tooltip"]').tooltip()

    // sync tabs with the same data-group
    $(".nav-tabs > li > a").click(function () {
        const group = $(this).attr("data-group");
        $(`.nav-tabs > li > a[data-group='${ group }']`).tab("show");
    });

    $('.language-dockerfile span.k').tooltip({
        title: function() {
            let c = this.textContent;
            this.style.cursor = 'help';
            $(this).on('click', () => { window.location.href = "/engine/reference/builder/#"+c.toLowerCase()});
            return 'Learn more about the "'+ c + '" Dockerfile command.'
        },
        placement: "auto"
    })
});
