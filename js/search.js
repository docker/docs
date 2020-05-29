var metadata, glossary;
var autoCompleteShowing = false;
var displayingAutcompleteResults = new Array();
var autoCompleteShowingID = 0;
var lastSearch = "";
var autoCompleteResultLimit = 3;
var results = new Array();
var scoreForTitleMatch = 10;
var scoreForURLMatch = 5;
var scoreForKeywordMatch = 3;
var scoreForDescriptionMatch = 1
function addResult(topic, matchesTitle, matchesDescription, matchesURL, matchesKeywords)
{
  var matchScore = (matchesTitle * scoreForTitleMatch) + (matchesDescription * scoreForDescriptionMatch) + (matchesURL * scoreForURLMatch) + (matchesKeywords * scoreForKeywordMatch);
  if (matchScore > 0)
  {
    var resultIndex = results.length;
    results[resultIndex] = new Array();
    results[resultIndex].topic = topic;
    results[resultIndex].score = matchScore;
  }
}
function loadPage(url)
{
  window.location.replace(url);
  window.location.href = url;
}

$(document).on("keypress", function(event) {
    if (event.keyCode == 13) {
      if(autoCompleteShowing) event.preventDefault();
    }
});

function highlightMe(inputTxt,keyword)
{
  inputTxt = String(inputTxt);
  simpletext = new RegExp("(" + keyword + ")","gi");
  return inputTxt.replace(simpletext, "<span>$1</span>")
}
function matches(inputTxt,searchTxt)
{
  var subs = inputTxt.split(searchTxt);
  return subs.length - 1;
}
function bindSearch()
{
  $("#st-search-input").on('keyup change', function(e) {
    e = e || window.event;
    if (autoCompleteShowing)
    {
      if (e.keyCode == '38') {
          // up arrow
          if (autoCompleteShowingID > -1)
          {
            // go up a result
            $("#autoCompleteResult" + autoCompleteShowingID).removeClass("autocompleteSelected");
            autoCompleteShowingID = autoCompleteShowingID - 1;
            $("#autoCompleteResult" + autoCompleteShowingID).addClass("autocompleteSelected");
            $("#autocompleteShowAll").removeClass("autocompleteSelected");
          } else {
            // de-selection auto-complete; reverting to raw search
            $("#autoCompleteResult0").removeClass("autocompleteSelected");
            autoCompleteShowingID = -1;
          }
      } else if (e.keyCode == '40') {
          // down arrow
          if (autoCompleteShowingID < (displayingAutcompleteResults.length - 1))
          {
            // go down to the next result
            $("#autoCompleteResult" + autoCompleteShowingID).removeClass("autocompleteSelected");
            autoCompleteShowingID = autoCompleteShowingID + 1;
            $("#autoCompleteResult" + autoCompleteShowingID).addClass("autocompleteSelected");
          } else {
            // select "See all results..." and go no further
            $("#autoCompleteResult" + autoCompleteShowingID).removeClass("autocompleteSelected");
            $("#autocompleteShowAll").addClass("autocompleteSelected");
            autoCompleteShowingID = autoCompleteResultLimit;
          }
      } else if (e.keyCode == '13') {
        // return key
        e.preventDefault();
        if (autoCompleteShowingID==autoCompleteResultLimit || autoCompleteShowingID == -1 || autoCompleteShowing == false)
        {
          // "see all" is selected or they don't have an autocomplete result selected
          loadPage("/search/?q=" + $("#st-search-input").val());
        } else {
          // an autocomplete result is selected
          loadPage(pages[displayingAutcompleteResults[autoCompleteShowingID]].url);
        }
      }
      //console.log('autoCompleteShowingID:',autoCompleteShowingID,'displayingAutcompleteResults[id]:',displayingAutcompleteResults[autoCompleteShowingID],'pages[id].url:',pages[displayingAutcompleteResults[autoCompleteShowingID]].url);
    }
    var searchVal = $("#st-search-input").val();
    if (lastSearch != searchVal)
    {
      displayingAutcompleteResults = [];
      results = [];
      var uppercaseSearchVal = searchVal.toUpperCase();
      //console.log("input changed: ",$("#st-search-input").val());

      if (searchVal.length > 2) {
        for (i=0;i<pages.length;i++)
        {
          // search url, description, title, and keywords for search input
          var thisPage = pages[i];
          var matchesTitle=0, matchesDescription=0, matchesURL=0, matchesKeywords=0;
          var matchesTitle = matches(String(thisPage.title).toUpperCase(),uppercaseSearchVal);
          //if (titleMatches > 0) console.log(uppercaseSearchVal,'matches',thisPage.title,titleMatches,'times');
          if (thisPage.description != null) {
            matchesDescription = matches(String(thisPage.description).toUpperCase(),uppercaseSearchVal);
          }
          if (thisPage.url != null) {
            matchesURL = matches(String(thisPage.url).toUpperCase(),uppercaseSearchVal);
          }
          if (thisPage.keywords != null) {
            matchesKeywords = matches(String(thisPage.keywords).toUpperCase(),uppercaseSearchVal);
          }
          addResult(i, matchesTitle, matchesDescription, matchesURL, matchesKeywords);
        }
        results.sort(function(a,b) {
          return b.score - a.score;
        });
      }
      if (results.length > 0)
      {
        autoCompleteShowingID = -1;
        var resultsShown = 0;
        var resultsOutput = new Array();
        resultsOutput.push("<div id='autoContainer'>")
        //console.log(results);
        for (i=0; i < autoCompleteResultLimit && i < results.length; i++)
        {
          //console.log(i, "of", autoCompleteResultLimit, "is underway");
          displayingAutcompleteResults.push(results[i].topic); //log results to global array
          resultsOutput.push("<div class='autoCompleteResult' id='autoCompleteResult" + i + "' onclick='loadPage(\"" + pages[results[i].topic].url + "\")'>");
          resultsOutput.push("<ul class='autocompleteList'>");
          resultsOutput.push("<li id='autoTitle" + i + "' class='autocompleteTitle'>")
          resultsOutput.push("<a href=" + pages[results[i].topic].url + ">" + highlightMe(pages[results[i].topic].title,searchVal) + "</a>");
          resultsOutput.push("</li>");
          resultsOutput.push("<li id='autoUrl" + i + "' class='autocompleteUrl'>")
          resultsOutput.push(highlightMe(pages[results[i].topic].url,searchVal));
          resultsOutput.push("</li>");
          /*
          resultsOutput.push("<li id='autoBreadcrumb" + i + "' class='autocompleteBreadcrumb'>")
          resultsOutput.push("Breadcrumb: " + breadcrumbString(pages[results[i]].url));
          resultsOutput.push("</li>");
          */
          if (pages[results[i].topic].keywords)
          {
          resultsOutput.push("<li id='autoKeywords" + i + "' class='autocompleteKeywords'>")
          resultsOutput.push("<b>Keywords</b>: <i>" + highlightMe(pages[results[i].topic].keywords,searchVal) + "</i>");
          resultsOutput.push("</li>");
          }
          if (pages[results[i].topic].description)
          {
          resultsOutput.push("<li id='autoDescription" + i + "' class='autocompleteDescription'>")
          resultsOutput.push("<b>Description</b>: " + highlightMe(pages[results[i].topic].description,searchVal));
          resultsOutput.push("</li>");
          }
          resultsOutput.push("</ul>");
          resultsOutput.push("</div>")
          resultsShown++;
        }
        var resultsShownText = (resultsShown > 1) ? resultsShown + " of " + results.length + " docs" : "doc";
        resultsOutput.push("<div id='autocompleteShowAll'><ul class='autocompleteList'><li class='autocompleteTitle' id='autoSeeAll'><a href='/search/?q=" + searchVal + "'><b>Showing top " + resultsShownText + ". See all results...</b></a></li></ul></div>")
        resultsOutput.push("</div>");
        $("#autocompleteResults").css("display","block");
        $("#autocompleteResults").html(resultsOutput.join(""));
        autoCompleteShowing = true;
      } else {
        $("#autocompleteResults").css("display","none");
        $("#autocompleteResults").html("");
        autoCompleteShowing = false;
      }
      lastSearch = searchVal;
    } // if searchVal != lastSearch
  });
}

function queryString()
{
    var vars = [], hash;
    var hashes = window.location.href.slice(window.location.href.indexOf('?') + 1).split('&');
    for(var i = 0; i < hashes.length; i++)
    {
        hash = hashes[i].split('=');
        vars.push(hash[0]);
        vars[hash[0]] = hash[1];
    }
    return vars;
}

function renderTopicsByTagTable(tagToLookup,divID)
{
  var matchingPages = new Array();
  for (i=0;i<pages.length;i++)
  {
    thisPage = pages[i];
    if (thisPage.keywords)
    {
      var keywordArray = thisPage.keywords.toString().split(",");
      for (n=0;n<keywordArray.length;n++)
      {
        if (keywordArray[n].trim().toLowerCase()==tagToLookup.toLowerCase())
        {
          matchingPages.push(i); // log the id of the page w/matching keyword
        }
      }
    }
  }
  var pagesOutput = new Array();
  if (matchingPages.length > 0)
  {
    pagesOutput.push("<h2>Pages tagged with: " + tagToLookup + "</h2>");
    pagesOutput.push("<table><thead><tr><td>Page</td><td>Description</td></tr></thead><tbody>");
    for(i=0;i<matchingPages.length;i++) {
      thisPage = pages[matchingPages[i]];
      pagesOutput.push("<tr><td><a href='" + thisPage.url + "'>" + thisPage.title + "</a></td><td>" + thisPage.description + "</td></tr>");
    }
    pagesOutput.push("</tbody></table>");
  }
  $("#" + divID).html(pagesOutput.join(""));
}

var tagToLookup;
function renderTagsPage()
{
  if(window.location.pathname.indexOf("/glossary/")>-1 || window.location.pathname.indexOf("/search/")>-1)
  {
    if (window.location.pathname.indexOf("/glossary/")>-1)
    {
      // Get ?term=<value>
      tagToLookup = decodeURI(queryString().term);
      $("#keyword").html(tagToLookup);
    }
    else
    {
      // Get ?q=<value>
      tagToLookup = decodeURI(queryString().q);
    }
    // Get the term and definition
    for (i=0;i<glossary.length;i++)
    {
      if (glossary[i].term.toLowerCase()==tagToLookup.toLowerCase())
      {
        var glossaryOutput = glossary[i].def;
      }
    }
    if (glossaryOutput) {
      $("#glossaryMatch").html("<h2>Definition of: " + tagToLookup + "</h2>" + glossaryOutput);
    }
    renderTopicsByTagTable(tagToLookup,"topicMatch",true);
  }
}
$(document).ready(bindSearch);
