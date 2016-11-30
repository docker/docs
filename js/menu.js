var metadata;
var autoCompleteShowing = false;
var displayingAutcompleteResults = new Array();
var autoCompleteShowingID = 0;
var lastSearch = "";
var autoCompleteResultLimit = 3;
function loadPage(url)
{
  window.location.replace(url);
  window.location.href = url;
}
$(document).on("keypress", function(event) {
    if (event.keyCode == 13) {
        event.preventDefault();
    }
});
function highlightMe(inputTxt,keyword)
{
  inputTxt = String(inputTxt);
  simpletext = new RegExp("(" + keyword + ")","gi");
  return inputTxt.replace(simpletext, "<span style='background-color:yellow'>$1</span>")
}
function hookupTOCEvents()
{
  // do after tree render
  $('.expand-menu').on('mouseup touchend', function(elem) {
//      menu = elem.currentTarget.nextElementSibling
    menu = elem.currentTarget.parentElement
    if (menu.classList.contains("menu-closed")) {
      menu.classList.remove("menu-closed")
      menu.classList.add("menu-open")
    } else {
      menu.classList.add("menu-closed")
      menu.classList.remove("menu-open")
    }
    return false;
  });
  $(".currentPage").each(function(){
    $(this).parentsUntil($('.docsidebarnav_section')).addClass("active").removeClass("menu-closed").addClass("menu-open");
  });
  $(".left-off-canvas-menu").css("display","block");
  // console.log(metadata);
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
            $("#autoSeeAll").removeClass("autocompleteSelected");
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
            $("#autoSeeAll").addClass("autocompleteSelected");
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
          loadPage(metadata.pages[displayingAutcompleteResults[autoCompleteShowingID]].url);
        }
      }
      //console.log('autoCompleteShowingID:',autoCompleteShowingID,'displayingAutcompleteResults[id]:',displayingAutcompleteResults[autoCompleteShowingID],'metadata.pages[id].url:',metadata.pages[displayingAutcompleteResults[autoCompleteShowingID]].url);
    }
    var searchVal = $("#st-search-input").val();
    if (lastSearch != searchVal)
    {
      displayingAutcompleteResults = [];
      var results = new Array();
      var uppercaseSearchVal = searchVal.toUpperCase();
      //console.log("input changed: ",$("#st-search-input").val());
      if (searchVal.length > 2) {
        for (i=0;i<metadata.pages.length;i++) {
          // search url, description, title, and keywords for search input
          var thisPage = metadata.pages[i];
          if (String(thisPage.title).toUpperCase().indexOf(uppercaseSearchVal) > -1) results.push(i);
          if (thisPage.description != null) if (String(thisPage.description).toUpperCase().indexOf(uppercaseSearchVal) > -1) results.push(i);
          if (thisPage.url != null) if (String(thisPage.url).toUpperCase().indexOf(uppercaseSearchVal) > -1) results.push(i);
          if (thisPage.keywords != null) if (String(thisPage.keywords).toUpperCase().indexOf(uppercaseSearchVal) > -1) results.push(i);
        }
      }
      if (results.length > 0)
      {
        var lastIndex = -1;
        var resultsOutput = new Array();
        var resultsShown  = 0;
        var highlightHTML = "";
        resultsOutput.push("<div id='autoContainer'>")
        for (i=0; resultsShown < autoCompleteResultLimit && i < results.length; i++)
        {
          if (results[i] != lastIndex) {
            // show results
            //console.log("match:",metadata.pages[results[i]]);
            displayingAutcompleteResults.push(results[i]); //log results to global array
            if (lastIndex==-1) {
              //first result!
              autoCompleteShowingID = -1;
            } else {
              highlightHTML = "";
            }
            resultsOutput.push("<div id='autoCompleteResult" + resultsShown + "'" + highlightHTML + " onclick='loadPage(\"" + metadata.pages[results[i]].url + "\")'>");
            resultsOutput.push("<ul class='autocompleteList'>");
            resultsOutput.push("<li id='autoTitle" + resultsShown + "' class='autocompleteTitle'>")
            resultsOutput.push("<a href=" + metadata.pages[results[i]].url + ">" + highlightMe(metadata.pages[results[i]].title,searchVal) + "</a>");
            resultsOutput.push("</li>");
            resultsOutput.push("<li id='autoUrl" + resultsShown + "' class='autocompleteUrl'>")
            resultsOutput.push(highlightMe(metadata.pages[results[i]].url,searchVal));
            resultsOutput.push("</li>");
            /*
            resultsOutput.push("<li id='autoBreadcrumb" + i + "' class='autocompleteBreadcrumb'>")
            resultsOutput.push("Breadcrumb: " + breadcrumbString(metadata.pages[results[i]].url));
            resultsOutput.push("</li>");
            */
            if (metadata.pages[results[i]].keywords)
            {
            resultsOutput.push("<li id='autoKeywords" + resultsShown + "' class='autocompleteKeywords'>")
            resultsOutput.push("Keywords: <i>" + highlightMe(metadata.pages[results[i]].keywords,searchVal) + "</i>");
            resultsOutput.push("</li>");
            }
            if (metadata.pages[results[i]].description)
            {
            resultsOutput.push("<li id='autoDescription" + resultsShown + "' class='autocompleteDescription'>")
            resultsOutput.push("Description: " + highlightMe(metadata.pages[results[i]].description,searchVal));
            resultsOutput.push("</li>");
            }
            resultsOutput.push("</ul>");
            resultsOutput.push("</div>")
            resultsShown++;
          }
          lastIndex = results[i];
        }
        resultsOutput.push("<ul class='autocompleteList'><li class='autocompleteTitle' id='autoSeeAll'><a href='/search/?q=" + searchVal + "'><b>See all results...</b></a></li></ul>")
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

jQuery(document).ready(function(){
    $.getJSON( "/metadata.txt", function( data ) {
      metadata = data;
      hookupTOCEvents();
    });
    $("#TableOfContents ul").empty();

    var prevH2Item = null;
    var prevH2List = null;

    var index = 0;
    var currentHeader = 0, lastHeader = 0;
    var output = "<ul>";
    $("h1, h2, h3, h4").each(function() {
        var li= "<li><a href='" + window.location + "#" + $(this).attr('id') + "'>" + $(this).text().replace("¶","") + "</a></li>";
        if( $(this).is("h2") ){
          // h2
          currentHeader = 2;
        } else if( $(this).is("h3") ){
          // h3
          currentHeader = 3;
        } else if( $(this).is("h4") ) {
          // h4
          currentHeader = 4;
        }
        //console.log("currentHeader ",currentHeader, "lastHeader ",lastHeader, "text ", $(this).text());
        if (currentHeader > lastHeader) {
            // nest further
            output += "<ul>"
        }
        if (currentHeader < lastHeader && lastHeader > 0) {
            // close nesting
            //console.log("Closing nesting because ", lastHeader, "is <", currentHeader);
            for (i=0; i < (lastHeader - currentHeader); i++)
            {
              output += "</ul>"
            }
        }
        output += li;
        lastHeader = currentHeader;
        /*
        if( $(this).is("h2") ){
            prevH2List = $("<ul></ul>");
            prevH2Item = $(li);
            prevH2Item.append(prevH2List);
            prevH2Item.appendTo("#TableOfContents ul");
        } else {
            prevH2List.append(li);
        }
        index++;*/
    });
    output += "</ul>";
    $("#TableOfContents").html(output);
});
