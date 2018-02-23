// Right nav highlighting
var sidebarObj = (document.getElementsByClassName("sidebar")[0]) ? document.getElementsByClassName("sidebar")[0] : document.getElementsByClassName("sidebar-home")[0]
var sidebarBottom = sidebarObj.getBoundingClientRect().bottom;
var footerTop = document.getElementsByClassName("footer")[0].getBoundingClientRect().top;
var headerOffset = document.getElementsByClassName("container-fluid")[0].getBoundingClientRect().bottom;

// ensure that the left nav visibly displays the current topic
var current = document.getElementsByClassName("active currentPage");
var body = document.getElementsByClassName("col-content content");
if (current[0]) {
    if (sidebarObj) {
      current[0].scrollIntoView(true);
      body[0].scrollIntoView(true);
    }
    // library hack
    if (document.location.pathname.indexOf("/samples/") > -1){
      $(".currentPage").closest("ul").addClass("in");
    }
  }

function addMyClass(classToAdd) {
  var classString = this.className; // returns the string of all the classes for myDiv
   // Adds the class "main__section" to the string (notice the leading space)
  this.className = newClass; // sets className to the new string
}

function navClicked(sourceLink)
{
  var classString = document.getElementById('#item' + sourceLink).className;
  if (classString.indexOf(' in') > -1)
  {
    //collapse
    var newClass = classString.replace(' in','');
    document.getElementById('#item' + sourceLink).className = newClass;
  } else {
    //expand
    var newClass = classString.concat(' in');
    document.getElementById('#item' + sourceLink).className = newClass;
  }

}

var outputHorzTabs = new Array();
var outputLetNav = new Array();
var totalTopics = 0;
var currentSection;
var sectionToHighlight;
function findMyTopic(tree)
{
  function processBranch(branch)
  {
    for (var k=0;k<branch.length;k++)
    {
      if (branch[k].section) {
        processBranch(branch[k].section);
      } else {
        if (branch[k].path == pageURL && !branch[k].nosync)
        {
          // console.log(branch[k].path + ' was == ' + pageURL)
          thisIsIt = true;
          break;
        } else {
          // console.log(branch[k].path + ' was != ' + pageURL)
        }
      }
    }
  }
  var thisIsIt = false;
  processBranch(tree)
  return thisIsIt;
}
function walkTree(tree)
{
  for (var j=0;j<tree.length;j++)
  {
    totalTopics++;
    if (tree[j].section)
    {
      var sectionHasPath = findMyTopic(tree[j].section);
      outputLetNav.push('<li><a onclick="navClicked(' + totalTopics +')" data-target="#item' + totalTopics +'" data-toggle="collapse" data-parent="#stacked-menu"')
      if (sectionHasPath)
      {
        outputLetNav.push('aria-expanded="true"')
      } else {
        outputLetNav.push('class="collapsed" aria-expanded="false"')
      }
      outputLetNav.push('>' + tree[j].sectiontitle + '<span class="caret arrow"></span></a>');
      outputLetNav.push('<ul class="nav collapse');
      if (sectionHasPath) outputLetNav.push(' in');
      outputLetNav.push('" id="#item' + totalTopics + '" aria-expanded="');
      if (sectionHasPath)
      {
        outputLetNav.push('true');
      } else {
        outputLetNav.push('false');
      }
      outputLetNav.push('">');
      var subTree = tree[j].section;
      walkTree(subTree);
      outputLetNav.push('</ul></li>');
    } else if (tree[j].generateTOC) {
      // auto-generate a TOC from a collection
      walkTree(collectionsTOC[tree[j].generateTOC])
    } else {
      // just a regular old topic; this is a leaf, not a branch; render a link!
      outputLetNav.push('<li><a href="' + tree[j].path + '"')
      if (tree[j].path == pageURL && !tree[j].nosync)
      {
        sectionToHighlight = currentSection;
        outputLetNav.push('class="active currentPage"')
      }
      outputLetNav.push('>'+tree[j].title+'</a></li>')
    }
  }
}
function renderNav(docstoc) {
  for (i=0;i<docstoc.horizontalnav.length;i++)
  {
    if (docstoc.horizontalnav[i].node != "glossary")
    {
      currentSection = docstoc.horizontalnav[i].node;
      // build vertical nav
      var itsHere = findMyTopic(docstoc[docstoc.horizontalnav[i].node]);
      if (itsHere || docstoc.horizontalnav[i].path == pageURL)
      {
        walkTree(docstoc[docstoc.horizontalnav[i].node]);
      }
    }
    // build horizontal nav
    outputHorzTabs.push('<li id="' + docstoc.horizontalnav[i].node + '"');
    if (docstoc.horizontalnav[i].path==pageURL || docstoc.horizontalnav[i].node==sectionToHighlight)
    {
      outputHorzTabs.push(' class="active"');
    }
    outputHorzTabs.push('><a href="'+docstoc.horizontalnav[i].path+'">'+docstoc.horizontalnav[i].title+'</a></li>\n');
  }
  if (outputLetNav.length==0)
  {
    // didn't find the current topic in the standard TOC; maybe it's a collection;
    for (var key in collectionsTOC)
    {
      var itsHere = findMyTopic(collectionsTOC[key]);
      if (itsHere) {
        walkTree(collectionsTOC[key]);
        break;
      }
    }
    // either glossary was true or no left nav has been built; default to glossary
    // show pages tagged with term and highlight term in left nav if applicable
    renderTagsPage()
    for (var i=0;i<glossary.length;i++)
    {
      var highlightGloss = '';
      if (tagToLookup) highlightGloss = (glossary[i].term.toLowerCase()==tagToLookup.toLowerCase()) ? ' class="active currentPage"' : '';
      outputLetNav.push('<li><a'+highlightGloss+' href="/glossary/?term=' + glossary[i].term + '">'+glossary[i].term+'</a></li>');
    }
  }
  document.getElementById('jsTOCHorizontal').innerHTML = outputHorzTabs.join('');
  document.getElementById('jsTOCLeftNav').innerHTML = outputLetNav.join('');
}

function highlightRightNav(heading)
{
  if (document.location.pathname.indexOf("/glossary/")<0){
    //console.log("highlightRightNav called on",document.location.pathname)
    if (heading == "title")
    {
      history.replaceState({},"Top of page on " + document.location.pathname,document.location.protocol +"//"+ document.location.hostname + (location.port ? ':'+location.port: '') + document.location.pathname);
      $("#my_toc a").each(function(){
        $(this).removeClass("active");
      });
      $("#sidebar-wrapper").animate({
        scrollTop: 0
      },800);
    } else {
      var targetAnchorHREF = document.location.protocol +"//"+ document.location.hostname + (location.port ? ':'+location.port: '') + document.location.pathname + "#" + heading;
      // make sure we aren't filtering out that heading level
      var noFilterFound = false;
      $("#my_toc a").each(function(){
        if (this.href==targetAnchorHREF) {
          noFilterFound = true;
        }
      });
      // now, highlight that heading
      if (noFilterFound)
      {
        $("#my_toc a").each(function(){
          //console.log("right-nav",this.href);
          if (this.href==targetAnchorHREF)
          {
            history.replaceState({},this.innerText,targetAnchorHREF);
            $(this).addClass("active");
            var sidebarOffset = (sidebarBottom > 200) ? 200 : headerOffset - 20;
            $("#sidebar-wrapper").animate({
              scrollTop: $("#sidebar-wrapper").scrollTop() + $(this).position().top - sidebarOffset
            },100);
            //document.getElementById("sidebar-wrapper").scrollTop = this.getBoundingClientRect().top - 200;
          } else {
            $(this).removeClass("active");
          }
        });
      }
    }
  }
}
var currentHeading = "";
$(window).scroll(function(){
  var headingPositions = new Array();
  $("h1, h2, h3, h4, h5, h6").each(function(){
    if (this.id == "") this.id="title";
    headingPositions[this.id]=this.getBoundingClientRect().top;
  });
  headingPositions.sort();
  // the headings have all been grabbed and sorted in order of their scroll
  // position (from the top of the page). First one is toppermost.
  for(var key in headingPositions)
  {
    if (headingPositions[key] > 0 && headingPositions[key] < 200)
    {
      if (currentHeading != key)
      {
        // a new heading has scrolled to within 200px of the top of the page.
        // highlight the right-nav entry and de-highlight the others.
        highlightRightNav(key);
        currentHeading = key;
      }
      break;
    }
  }
});


// Cookie functions
function createCookie(name,value,days) {
    var expires = "";
    if (days) {
        var date = new Date();
        date.setTime(date.getTime() + (days*24*60*60*1000));
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + value + expires + "; path=/";
}

function readCookie(name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(';');
    for(var i=0;i < ca.length;i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1,c.length);
        if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
    }
    return null;
}

function eraseCookie(name) {
    createCookie(name,"",-1);
}

if (readCookie("night") == "true") {
  applyNight();
  $('#switch-style').prop('checked', true);
} else {
  applyDay();
  $('#switch-style').prop('checked', false);
}



/*
 *
 * toggle menu *********************************************************************
 *
 */

$("#menu-toggle").click(function(e) {
        e.preventDefault();
        $(".wrapper").toggleClass("right-open");
        $(".col-toc").toggleClass("col-toc-hidden");
    });
$("#menu-toggle-left").click(function(e) {
        e.preventDefault();
        $(".col-nav").toggleClass("col-toc-hidden");
    });
$(".navbar-toggle").click(function(){
  $("#sidebar-nav").each(function(){
    $(this).toggleClass("hidden-sm");
    $(this).toggleClass("hidden-xs");
  });
});

var navHeight = $('.navbar').outerHeight(true) + 80;

$(document.body).scrollspy({
	target: '#leftCol',
	offset: navHeight
});

function loadHash(hashObj)
{
  // Using jQuery's animate() method to add smooth page scroll
  // The optional number (800) specifies the number of milliseconds it takes to scroll to the specified area
  $('html, body').animate({
    scrollTop: $(hashObj).offset().top-80
  }, 800);
}

$(document).ready(function(){
  // Add smooth scrolling to all links
  // $( ".toc-nav a" ).addClass( "active" );
  $(".toc-nav a").on('click', function(event) {
    // $(this).addClass('active');
    // Make sure this.hash has a value before overriding default behavior
    if (this.hash !== "") {
      // Prevent default anchor click behavior
      event.preventDefault();

      // Store hash
      var hash = this.hash;
      loadHash(hash);

      // Add hash (#) to URL when done scrolling (default click behavior)
      window.location.hash = hash;

    } // End if
  });
  if (window.location.hash) loadHash(window.location.hash);
});


$(document).ready(function(){
  $(".sidebar").Stickyfill();

  // Add smooth scrolling to all links
  $(".nav-sidebar ul li a").on('click', function(event) {

    // Make sure this.hash has a value before overriding default behavior
    if (this.hash !== "") {
      // Prevent default anchor click behavior
      event.preventDefault();

      // Store hash
      var hash = this.hash;

      // Using jQuery's animate() method to add smooth page scroll
      // The optional number (800) specifies the number of milliseconds it takes to scroll to the specified area
      $('html, body').animate({
        scrollTop: $(hash).offset().top-80
      }, 800, function(){

        // Add hash (#) to URL when done scrolling (default click behavior)
        window.location.hash = hash;
      });
    } // End if
  });
});


/*
 *
 * make dropdown show on hover *********************************************************************
 *
 */

$('ul.nav li.dropdown').hover(function() {
  $(this).find('.dropdown-menu').stop(true, true).delay(200).fadeIn(500);
}, function() {
  $(this).find('.dropdown-menu').stop(true, true).delay(200).fadeOut(500);
});

/*
 *
 * swapStyleSheet*********************************************************************
 *
 */

// function swapStyleSheet(sheet) {
//     document.getElementById('pagestyle').setAttribute('href', sheet);
// }

function applyNight()
{
  $( "body" ).addClass( "night" );
}

function applyDay() {
  $( "body" ).removeClass( "night" );
}

$('#switch-style').change(function() {

    if ($(this).is(':checked')) {
        applyNight();
        createCookie("night",true,999)
    } else {
        applyDay();
    //     swapStyleSheet('/css/style.css');
        eraseCookie("night")
    }
});


/*
 *
 * TEMP HACK For side menu*********************************************************************
 *
 */

$('.nav-sidebar ul li a').click(function() {
    $(this).addClass('collapse').siblings().toggleClass('in');
});

if($('.nav-sidebar ul a.active').length != 0)
{
  $('.nav-sidebar ul').click(function() {
      $(this).addClass('collapse in').siblings;
  });
}


/*
 *
 * Components *********************************************************************
 *
 */

$(function () {
  $('[data-toggle="tooltip"]').tooltip()
})

// Enable glossary link popovers
$('.glossLink').popover();

// sync tabs with the same data-group
window.onload = function() {
  $('.nav-tabs > li > a').click(function(e) {
    var group = $(this).attr('data-group');
    $('.nav-tabs > li > a[data-group="'+ group +'"]').tab('show');
  })

  // isArchive is set by logic in archive.js
  if ( isArchive == false ) {
    //console.log("Showing content that should only be in the current version.");
    // Hide elements that are not appropriate for archives
    // PollDaddy
    $('#ratings-div').css("visibility","visible");
    //console.log("Ratings widget shown.");
    // Archive drop-down
    $('.ctrl-right .btn-group').css("visibility","visible");
    //console.log("Archive widget shown.");
    // Swarch
    $('.search-form').css("visibility","visible");
    //console.log("Search widget shown.");
    // Page edit link
    $('.feedback-links li').first().css("visibility","visible");
    //console.log("Page edit link shown.");
  } else {
    //console.log("Keeping non-applicable elements hidden.");
  }
};
