import Fuse from "fuse.js";

let indexed = false;
let handler = null;

const modalSearchInput = document.querySelector("#modal-search-input");
const modalSearchResults = document.querySelector("#modal-search-results");

async function initializeIndex() {
  const index = await fetch("/metadata.json").then((response) =>
    response.json(),
  );

  const options = {
    keys: [
      { name: "title", weight: 2 },
      { name: "description", weight: 1 },
      { name: "keywords", weight: 1 },
      { name: "tags", weight: 1 },
    ],
    minMatchCharLength: 1,
    threshold: 0.2,
    ignoreLocation: true,
    useExtendedSearch: true,
    ignoreFieldNorm: true,
  };

  handler = new Fuse(index, options);
  indexed = true;
}

async function executeSearch(query) {
  !indexed && (await initializeIndex());
  const results = handler.search(query);
  return results;
}

async function modalSearch(e) {
  const query = e.target.value;
  results = await executeSearch(query);

  let resultsHTML = `<div>${results.length} results</div>`;
  resultsHTML += results
    .map(({ item }) => {
      return `<div class="bg-gray-light-100 dark:bg-gray-dark-200 rounded p-4">
      <div class="flex flex-col">
        <a class="link" href="${item.url}">${item.title}</a>
        <p>${item.description}</p>
      </div>
      </div>`;
    })
    .join("");

  modalSearchResults.innerHTML = resultsHTML;
}

modalSearchInput.addEventListener("input", modalSearch);
