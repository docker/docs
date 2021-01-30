
const maxResults = 3, titleWeight = 10, urlWeight = 5, keywordWeight = 3, descriptionWeight = 1
let searchVal = "", pages = []

function handleKeyNav(/* KeyboardEvent */ e) {
    let row = _(".autocompleteResult.selected")
    switch (e.key) {
        case "ArrowUp":
            if (row && row.previousElementSibling) {
                row.classList.remove("selected")
                row = row.previousElementSibling
                row.classList.add("selected")
            }
            break;
        case "ArrowDown":
            if (!row) {
                // pick the first one
                row = _(".autocompleteResult")
            } else if (row.nextElementSibling) {
                row.classList.remove("selected")
                row = row.nextElementSibling
            }
            if (row) {
                row.classList.add("selected")
            }
            break;
        case "Enter":
            e.preventDefault();
            if (!row || row.id === "autocompleteShowAll") {
                // "see all" is selected or no autocomplete result selected
                window.location.href = "/search/?q=" + e.target.value;
            } else {
                // an autocomplete result is selected
                row.click()
            }
            break;
    }
}

function matches(input, search) {
    return String(input).toUpperCase().split(search.toUpperCase()).length - 1;
}

function handleSearch(/* KeyboardEvent */ e) {
    if (e.target.value === searchVal) {
        // no new search
        return
    }
    searchVal = e.target.value
    let results = [];
    if (searchVal.length > 2) {
        for (let i = 0; i < pages.length; i++) {
            // search url, description, title, and keywords for search input
            const p = pages[i];
            if (!p.title) {
                continue
            }
            let score = (matches(p.title, searchVal) * titleWeight)
            if (p.description != null) {
                score += (matches(p.description, searchVal) * descriptionWeight)
            }
            if (p.url != null) {
                score += (matches(p.url, searchVal) * urlWeight)
            }
            if (p.keywords != null) {
                score += (matches(p.keywords, searchVal) * keywordWeight)
            }
            if (score > 0) {
                results.push({ "topic": i, "score": score });
            }
        }
    }
    let rows = []
    if (results.length > 0) {
        results.sort((a, b) => b.score - a.score);
        const match = new RegExp(`(${searchVal})`, "gi");
        const highlight = function (/* String */ content) {
            return content.replace(match, "<span>$1</span>")
        }
        for (let i = 0; i < maxResults && i < results.length; i++) {
            const p = pages[results[i].topic];
            rows.push(`<div class='autocompleteResult' onclick='window.location.href = "${p.url}"'><ul>`);
            rows.push(`<li><a class='title' href='${p.url}'>${ highlight(p.title) }</a></li>`);
            if (p.description && p.description !== p.title) {
                // Omit description if it's the same as the title
                rows.push(`<li>${ highlight(p.description) }</li>`);
            }
            if (p.keywords) {
                rows.push(`<li class='keywords'><span class='glyphicon glyphicon-tags'></span><i>${ highlight(p.keywords) }</i></li>`);
            }
            rows.push("</ul></div>")
        }
        let shown = Math.min(results.length, maxResults)
        rows.push(`<div class='autocompleteResult' id='autocompleteShowAll'><ul><li>Showing ${shown} of ${results.length} ${ (shown > 1) ? "results" : "result" }. <a href='/search/?q=${searchVal}'>See all results...</a></li></ul></div>`)
    }

    const out = _("#autocompleteResults")
    if (out) {
        out.innerHTML = rows.join("");
        let shown = Math.min(results.length, maxResults)
        out.style.display = shown === 0 ? "none" : "block"
    }
}

ready(() => {
    getJSON( "/js/metadata.json", function(data) {
        pages = data
        const input = _("#st-search-input")
        if (/* HTMLInputElement */ input) {
            input.form.addEventListener('submit', (e) => e.preventDefault());
            input.addEventListener('keyup', handleKeyNav, {capture: true })
            input.addEventListener('keyup', debounce(handleSearch, 100), {capture: true })
        }
    });
})
