---
description: Docker documentation search results
keywords: Search, Docker, documentation, manual, guide, reference, api
noratings: true
notoc: true
notags: true
title: "Docs search <span id='searchTerm'></span>"
tree: false
---

<style type='text/css'>

</style>

<div id="glossaryMatch"></div>

<script defer>
// Replace the subscriptionKey string value with your valid subscription key.
var subscriptionKey = "a71972579d8640d38b3bc859d7c4f1c3";
var customconfig = "3956951448";

var first = 'First'; // Override for Chinese.
var last = 'Last';
var prev = 'Prev';
var next = 'Next';
var mkt = "en-us";
var noResultHtml = 'No results found...';

/* CN Version:

first: '首页',
last: '尾页',
prev: '上页',
next: '下页',
var noResultHtml = '没有结果!';
*/

function doBingPagingSearch(page) {

    var searchText = decodeURI(getQueryString().q);
    if (searchText != "undefined" && searchText != '') {
        if (page == "undefined") {
            page = 1;
            startPos = 0;
        } else {
            startPos = (page - 1) * 10;
        }

        var bingEndPoint = "https://api.cognitive.microsoft.com/bingcustomsearch/v5.0/search";

        // Request parameters.
        var reqParams = {
            "q": searchText,
            "customconfig": customconfig,
            "responseFilter": "Webpages",
            "mkt": mkt,
            "safesearch": "Moderate",
            "count": "10",
            "offset": startPos,
        };

        $.ajax({
            url: bingEndPoint + "?" + $.param(reqParams),
            beforeSend: function (xhrObj) {
            xhrObj.setRequestHeader("Content-Type", "application/json");
            xhrObj.setRequestHeader("Ocp-Apim-Subscription-Key", subscriptionKey);
            },
            type: "GET",
        })
        .done(function (data) {
            var pageHits = data.webPages.value;
            var totalPageHits = data.webPages.totalEstimatedMatches;
            console.log("pageHits: ",pageHits,"totalPageHits: ",totalPageHits);
            if (totalPageHits > 0) {
                var totalPageNum = Math.ceil(totalPageHits / 10);
                var $pagination = $('#pagination-result');
                var paginationOpts = {
                    totalPages: totalPageNum,
                    visiblePages: 5,
                    first: first,
                    last: last,
                    prev: prev,
                    next: next,
                    initiateStartPageClick: false,
                    onPageClick: function (event, page) {
                        doBingPagingSearch(page);
                    }
                };

                $pagination.twbsPagination(paginationOpts);
                var searchResultHtml = "";
                for (var i = 0; i < pageHits.length; i++) {
                    var item = pageHits[i];

                    var title = item.name;
                    var url = item.url;
                    var desc = item.snippet;
                    var descHtml = "<div class='result-desc'>" + desc + "</div>";

                    // hightlight keywords start
                    searchText = searchText.replace(/(\s+)/, "(<[^>]+>)*$1(<[^>]+>)*");
                    var pattern = new RegExp("(" + searchText + ")", "gi");

                    title = title.replace(pattern, "<b>$1</b>");
                    title = title.replace(/(<b>[^<>]*)((<[^>]+>)+)([^<>]*<\/b>)/, "$1</b>$2<b>$4");
                    // hightlight keywords end

                    var titleHtml = "<a class='result-title' href='" + url + "'>" + title + "</a>";

                    var urlHtml = "<div class='result-url'>" + url + "</div>";

                    searchResultHtml += "<div class='result-wrap'>" + titleHtml + urlHtml + descHtml + "</div>";

                }

                $("#search-result").html(searchResultHtml);
            } else {

                $("#search-result").html(noResultHtml);
            }

        })
        .fail(function (jqXHR, textStatus, errorThrown) {
            var errorString = (errorThrown === "") ? "Error. " : errorThrown + " (" + jqXHR.status + "): " + jqXHR.responseText;
            /*
            errorString += (jqXHR.responseText === "") ? "" : (jQuery.parseJSON(jqXHR.responseText).message) ?
            jQuery.parseJSON(jqXHR.responseText).message : jQuery.parseJSON(jqXHR.responseText).error.message;
            */
            console.log(errorString);
        });
    }
}

function getQueryString() {
    var vars = [], hash;
    var hashes = window.location.href.slice(window.location.href.indexOf('?') + 1).split('&');
    for (var i = 0; i < hashes.length; i++) {
        hash = hashes[i].split('=');
        vars.push(hash[0]);
        vars[hash[0]] = hash[1];
    }
    return vars;
}

setTimeout(function(){

  $(document).ready(function () {

    $(document).ajaxStart(function(){
        $("#ajax_loading").show();
    }).ajaxComplete(function(){
        $("#ajax_loading").hide();
    });
    doBingPagingSearch(1);

    if (decodeURI(queryString().q) != "undefined" && decodeURI(queryString().q) && decodeURI(queryString().q).length > 0) {
      $("#st-search-input").val(decodeURI(queryString().q));
      $("#st-search-input").focus();
      $("#searchTerm").html("results for: " + decodeURI(queryString().q))
    }
  });
}, 1);
</script>
<script>var $ = jQuery.noConflict();</script>
<div id="my-cse1">
  <div id="search-result"></div>
  <nav aria-label="Page navigation">
    <ul class="pagination" id="pagination-result">
    </ul>
  </nav>
  <gcse:searchresults-only></gcse:searchresults-only>
</div>
