---
layout: null
---

/* Only run this if we are online*/
if (window.navigator.onLine) {
  var suppressButterBar = false;
  /* This JSON file contains a current list of all docs versions of Docker */
  $.getJSON("/js/archives.json", function(result){
    var outerDivStart = '<div id="archive-butterbar"><div class="container"><div style="text-align: center"><span id="archive-list">This is <b><a href="https://docs.docker.com/docsarchive/" style="color: #254356; text-decoration: underline !important">archived documentation</a></b> for Docker&nbsp;' + dockerVersion + '. Go to the <a style="color: #254356; text-decoration: underline !important" href="https://docs.docker.com/">latest docs</a> or a different version:&nbsp;&nbsp;</span>' +
                               '<span style="z-index: 1001" class="dropdown">';
    var listStart = '<ul class="dropdown-menu" role="menu" aria-labelledby="archive-menu">';
    var listEnd = '</ul>';
    var outerDivEnd = '</span></div></div></div>';
    var buttonCode = null;
    var listItems = new Array();
    $.each(result, function(i, field){
      if ( field.name == dockerVersion && field.current ) {
        // We are the current version so we don't need a butterbar
        suppressButterBar = true;
      } else {
        var prettyName = 'Docker ' + field.name.replace("v", "");
        // If this archive has current = true, and we don't already have a button
        if ( field.current && buttonCode == null ) {
          // Get the button code
          buttonCode = '<button id="archive-menu" data-toggle="dropdown" class="btn dropdown-toggle" style="border: 1px solid #254356; background-color: #fff; color: #254356;">' + prettyName + '&nbsp;(current) &nbsp;<span class="caret"></span></button>';
          // The link is different for the current release
          listItems.push('<li role="presentation"><a role="menuitem" tabindex="-1" href="https://docs.docker.com/">' + prettyName + '</a></li>');
        } else {
          listItems.push('<li role="presentation"><a role="menuitem" tabindex="-1" href="https://docs.docker.com/' + field.name + '/">' + prettyName + '</a></li>');
        }
      }
    });

    // only append the butterbar if we are NOT the current version
    // Also set the isArchive variable to true if it's an archive. It defaults
    // to true, set in _layouts/docs.html. We default to true because it looks
    // better in CSS to show stuff than to hide stuff onLoad.
    if ( suppressButterBar == false ) {
      $( 'body' ).prepend(outerDivStart + buttonCode + listStart + listItems.join("") + listEnd + outerDivEnd);
      isArchive = true;
      // If the butterbar exists, deal with positioning it
      // Depends on some logic in _layout/docs.html
      $(document).scroll(function() {
        if ( $( 'nav' ).hasClass( 'affix' ) ) {
          $('#archive-butterbar').addClass('fixed').removeClass('top');
        } else {
          $('#archive-butterbar').addClass('top').removeClass('fixed');
        }
      });
    } else {
      isArchive = false;
      /* This is only relevant to /enterprise/index.md */
      if (document.getElementById('ee-version-div')) {
        document.getElementById('ee-version-div').textContent += "The latest version of Docker EE Engine is {{ site.docker_ee_version }}.";
      }
    }  });
}
