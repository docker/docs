var metadata;
var autoCompleteShowing = false;
var displayingAutcompleteResults = new Array();
document.onkeydown = checkKey;
function checkKey(e) {
  if (autoCompleteShowing) {
    e = e || window.event;
    if (e.keyCode == '38') {
        // up arrow
    }
    else if (e.keyCode == '40') {
        // down arrow
    }
    else if (e.keyCode == '37') {
       // left arrow
    }
    else if (e.keyCode == '39') {
       // right arrow
    }
  }
}
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
  $("#st-search-input").on('keyup change', function() {
    var results = new Array();
    var searchVal = $("#st-search-input").val();
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
      resultsOutput.push("<ul class='autocompleteList'>")
      for (i=0; resultsShown < 3 && i < results.length; i++)
      {
        if (results[i] != lastIndex) {
          // show results
          //console.log("match:",metadata.pages[results[i]]);
          displayingAutcompleteResults.push(i);
          resultsOutput.push("<li id='autoTitle" + i + "' class='autocompleteTitle'>")
          resultsOutput.push("<a href=" + metadata.pages[results[i]].url + ">" + highlightMe(metadata.pages[results[i]].title,searchVal) + "</a>");
          resultsOutput.push("</li>");
          resultsOutput.push("<li id='autoUrl" + i + "' class='autocompleteUrl'>")
          resultsOutput.push(highlightMe(metadata.pages[results[i]].url,searchVal));
          resultsOutput.push("</li>");
          /*
          resultsOutput.push("<li id='autoBreadcrumb" + i + "' class='autocompleteBreadcrumb'>")
          resultsOutput.push("Breadcrumb: " + breadcrumbString(metadata.pages[results[i]].url));
          resultsOutput.push("</li>");
          */
          if (metadata.pages[results[i]].keywords)
          {
          resultsOutput.push("<li id='autoKeywords" + i + "' class='autocompleteKeywords'>")
          resultsOutput.push("Keywords: <i>" + highlightMe(metadata.pages[results[i]].keywords,searchVal) + "</i>");
          resultsOutput.push("</li>");
          }
          if (metadata.pages[results[i]].description)
          {
          resultsOutput.push("<li id='autoDescription" + i + "' class='autocompleteDescription'>")
          resultsOutput.push("Description: " + highlightMe(metadata.pages[results[i]].description,searchVal));
          resultsOutput.push("</li>");
          }
          resultsShown++;
        }
        lastIndex = results[i];
      }
      resultsOutput.push("<li class='autocompleteTitle'><a href='/search/?q=" + searchVal + "'><b>See all results...</b></a></li>")
      resultsOutput.push("</ul>");
      $("#autocompleteResults").css("display","block");
      $("#autocompleteResults").html(resultsOutput.join(""));
      autoCompleteShowing = true;
    } else {
      $("#autocompleteResults").css("display","none");
      $("#autocompleteResults").html("");
      autoCompleteShowing = false;
    }
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
