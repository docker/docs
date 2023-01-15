import algoliasearch from "algoliasearch/lite"
import { autocomplete, getAlgoliaResults } from "@algolia/autocomplete-js"

const searchClient = algoliasearch(
  "1G07Z8655T",
  "449796e3912f98c903ac2b2ca0a35020"
)

autocomplete({
  container: "#autocomplete",
  placeholder: "Search the docs",
  getSources({ query }) {
    return [
      {
        sourceId: "objectID",
        getItems() {
          return getAlgoliaResults({
            searchClient,
            queries: [
              {
                indexName: "docs_test",
                query,
                params: {
                  hitsPerPage: 5,
                  attributesToSnippet: ["title:10"],
                  snippetEllipsisText: "…"
                }
              }
            ]
          })
        },
        templates: {
          item({ item, components, html }) {
            return html`<div class="block">
              <a class="inline-block" href="${item.permalink}">
                ${components.Highlight({
                  hit: item,
                  attribute: "title"
                })}</a
              >
            </div>`
          }
        },
        getItemUrl({ item }) {
          return item.permalink
        }
      }
    ]
  }
})
