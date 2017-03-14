$(document).ready(function ()
{

  /* Obly run this if we are online*/
  if (window.navigator.onLine) {
    var dockerVersion = '1.6';
    /* This JSON file contains a current list of all docs versions of Docker */
    $.getJSON("https://docs.docker.com/js/archives.json", function(result){
      var outerDivStart = '<div class="row" style="padding-top: 10px; padding-bottom: 10px; min-height: 34px; background-color: #FFD2A1; color: #254356"><div class="textcenter"><span id="archive-list">This is <b><a href="https://docs.docker.com/docsarchive/" style="color: #254356; text-decoration: underline !important">archived documentation</a></b> for Docker&nbsp;' + dockerVersion + '. Go to the <a style="color: #254356; text-decoration: underline !important" href="https://docs.docker.com/">latest docs</a> or a different version:&nbsp;&nbsp;</span>' +
                                 '<span style="z-index: 1001" class="dropdown">';
      var listStart = '<ul class="dropdown-menu" role="menu" aria-labelledby="archive-menu">';
      var listEnd = '</ul>';
      var outerDivEnd = '</span></div></div>';
      var buttonCode = null;
      var listItems = new Array();
      $.each(result, function(i, field){
        var prettyName = 'Docker ' + field.name.replace("v", "");
        // If this archive has current = true, and we don't already have a button
        if ( field.current && buttonCode == null ) {
          // Get the button code
          buttonCode = '<button id="archive-menu" data-toggle="dropdown" class="btn dropdown-toggle" style="border: 1px solid #82949E; background-color: #FBFBFC; color: #254356;">' + prettyName + '&nbsp;(current) &nbsp;<span class="caret"></span></button>';
          // The link is different for the current release
          listItems.push('<li role="presentation"><a role="menuitem" tabindex="-1" href="https://docs.docker.com/">' + prettyName + '</a></li>');
        } else {
          listItems.push('<li role="presentation"><a role="menuitem" tabindex="-1" href="https://docs.docker.com/' + field.name + '/">' + prettyName + '</a></li>');
        }
      });
      $( 'body' ).prepend(outerDivStart + buttonCode + listStart + listItems.join("") + listEnd + outerDivEnd);
    });
  }

  prettyPrint();

  // Resizing
  resizeMenuDropdown();
  // checkToScrollTOC();
  $(window).resize(function() {
    if(this.resizeTO)
    {
      clearTimeout(this.resizeTO);
    }
    this.resizeTO = setTimeout(function ()
    {
      resizeMenuDropdown();
      // checkToScrollTOC();
    }, 500);
  });

  /* Follow TOC links (ScrollSpy) */
  $('body').scrollspy({
    target: '#toc_table',
  });

  /* Prevent disabled link clicks */
  $("li.disabled a").click(function ()
  {
    event.preventDefault();
  });

  // Submenu ensured drop-down functionality for desktops & mobiles
  $('.dd_menu').on({
    click: function ()
    {
      $(this).toggleClass('dd_on_hover');
    },
    mouseenter: function ()
    {
      $(this).addClass('dd_on_hover');
    },
    mouseleave: function ()
    {
      $(this).removeClass('dd_on_hover');
    },
  });

  function getURLP(name)
  {
    return decodeURIComponent((new RegExp('[?|&]' + name + '=' + '([^&;]+?)(&|#|;|$)').exec(location.search)||[,""])[1].replace(/\+/g, '%20')) || null;
  }
  if (getURLP("q")) {
    // Tipue Search activation
    $('#tipue_search_input').tipuesearch({
      'mode': 'json',
      'contentLocation': '/search_content.json.gz'
    });
  }

});

function resizeMenuDropdown ()
{
  $('.dd_menu > .dd_submenu').css("max-height", ($('body').height() - 160) + 'px');
}

// https://github.com/bigspotteddog/ScrollToFixed
function checkToScrollTOC ()
{
  if ( $(window).width() >= 768 )
  {
    // If TOC is hidden, expand.
    $('#toc_table > #toc_navigation').css("display", "block");
    // Then attach or detach fixed-scroll
    if ( ($('#toc_table').height() + 100) >= $(window).height() )
    {
      $('#toc_table').trigger('detach.ScrollToFixed');
      $('#toc_navigation > li.active').removeClass('active');
    }
    else
    {
      $('#toc_table').scrollToFixed({
        marginTop: $('#nav_menu').height(),
        limit: function () { return $('#footer').offset().top - 450; },
        zIndex: 1,
        minWidth: 768,
        removeOffsets: true,
      });
    }
  }
}

function getCookie(cname) {
  var name = cname + "=";
  var ca = document.cookie.split(';');
  for(var i=0; i<ca.length; i++) {
      var c = ca[i].trim();
      if (c.indexOf(name) == 0) return c.substring(name.length,c.length);
  }
  return "";
}
