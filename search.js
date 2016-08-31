


var search = instantsearch({
  appId: 'BH4D9OD16A',
  apiKey: '828e35034e76dbb6fdc5b33c03d5702f',
  indexName: 'docker',
  urlSync: true,
  searchFunction: function(helper) {
    if (helper.state.query === '') {
      // preload the search with the url path if its the 404 page?
      if (document.title == '404 Page not found') {
        helper.state.query = document.location.pathname;
      } else {
        return;
      }
    }
    helper.search();
  }
});
search.addWidget(
  instantsearch.widgets.searchBox({
    container: '#search-box',
    placeholder: 'Search the Documentation...',
    queryHook: function(query, searchFn) {
      if (query === '') {
        // preload the search with the url path if its the 404 page?
        if (document.title == '404 Page not found') {
          query = document.location.pathname;
        } else {
          return;
        }
      }
    searchFn(query);
  },
    hitsPerPage: 10,
    autofocus: true,
    poweredBy: true
  })
);

var renderItem =
'    <div class="algolia-docsearch-suggestion--category-header">' +
'        <span class="algolia-docsearch-suggestion--category-header-lvl0">{{_highlightResult.hierarchy.lvl0.value}}</span>' +
'    </div>' +
'    <div class="algolia-docsearch-suggestion--wrapper">' +
'      <div class="algolia-docsearch-suggestion--subcategory-column">' +
'        <span class="algolia-docsearch-suggestion--subcategory-column-text">{{_highlightResult.hierarchy.lvl1.value}}</span>' +
'      </div>' +
'      <div class="algolia-docsearch-suggestion--content">' +
'        <div class="algolia-docsearch-suggestion--subcategory-inline">{{_highlightResult.hierarchy.lvl2.value}}</div>' +
'        <div class="algolia-docsearch-suggestion--title">{{_highlightResult.hierarchy.lvl3.value}}</div>' +
'        <div class="algolia-docsearch-suggestion--text">{{content}} </div>' +
'      </div>' +
'    </div>' +
'  </div>';

var lastLvl0Hit = '';
var lastLvl1Hit = '';
var onRenderHandler = function() {
  lastLvl0Hit = '';
  lastLvl1Hit = '';
};

search.on('render', onRenderHandler);

search.addWidget(
  instantsearch.widgets.hits({
    container: '#hits-container',
    hitsPerPage: 10,
    templates: {
      item: function(data) {
  result = '';
  if (data._highlightResult.hierarchy.hasOwnProperty('lvl0')) {
    // Only show the lvl0 category if its different from the hit above.
    if (lastLvl0Hit != data._highlightResult.hierarchy.lvl0.value) {
      result += '    <div class="algolia-docsearch-suggestion--category-header">';
      result += '        <span class="algolia-docsearch-suggestion--category-header-lvl0">'+ data._highlightResult.hierarchy.lvl0.value + '</span>';
      result += '    </div>';
      lastLvl0Hit = data._highlightResult.hierarchy.lvl0.value
    }
  }
  result += '    <div class="algolia-docsearch-suggestion--wrapper">';
  if (data._highlightResult.hierarchy.hasOwnProperty('lvl1')) {
    // Only show the lvl1 category if its different from the hit above.
    if (lastLvl1Hit != data._highlightResult.hierarchy.lvl1.value) {
      result += '      <div class="algolia-docsearch-suggestion--subcategory-column">';
  result += '<a href="'+ data.url +'">';
      result += '        <span class="algolia-docsearch-suggestion--subcategory-column-text">'+ data._highlightResult.hierarchy.lvl1.value +'</span>';
  result += '</a>';
      result += '      </div>';
      lastLvl1Hit = data._highlightResult.hierarchy.lvl1.value
    }
  }
  result += '      <div class="algolia-docsearch-suggestion--content">';
  if (data._highlightResult.hierarchy.hasOwnProperty('lvl2')) {
  result += '<a href="'+ data.url +'">';
    result += '        <div class="algolia-docsearch-suggestion--subcategory-inline">'+ data._highlightResult.hierarchy.lvl2.value +'</div>';
  result += '</a>';
  }
  if (data._highlightResult.hierarchy.hasOwnProperty('lvl3')) {
  result += '<a href="'+ data.url +'">';
    result += '        <div class="algolia-docsearch-suggestion--title">'+ data._highlightResult.hierarchy.lvl3.value +'</div>';
  result += '</a>';
  }
  if (data.content != null) {
  result += '<a href="'+ data.url +'">';
    result += '        <div class="algolia-docsearch-suggestion--text">'+ data.content +'</div>';
  result += '</a>';
  }
  result += '    </div>';
  result += '  </div>';
  return result;
  }
    }
  })
);

search.addWidget(
  instantsearch.widgets.pagination({
    container: '#pagination-container',
    hitsPerPage: 10
  })
);
search.start();
