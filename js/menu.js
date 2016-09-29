var tocData;
var treeOutput = new Array();
function renderTree(tree)
{
  for (var i=0; i < tree.length; i++)
  {
    if (tree[i].heading)
    {
      // render a heading
    } else {
      if (tree[i].sectiontitle) {
        treeOutput.push('<li class="leaf menu-closed"><a href="#" class="expand-menu "><span class="menu-icon" aria-hidden="true"></span>' + tree[i].sectiontitle + '</a>');
        treeOutput.push('<ul class="nav-sub">');
        renderTree(tree[i].section);
        treeOutput.push('</ul>');
      } else {
        treeOutput.push('<li class="leaf"><a href="' + tree[i].path +'" class="');
        if (tree[i].path == window.location.pathname) treeOutput.push('active currentPage');
        treeOutput.push('">' + tree[i].title + '</a></li>');
      }
    }
  }
}
function hookupTOCEvents()
{
  // do after tree render
  $('.expand-menu').on('click touchstart', function(elem) {
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
}
jQuery(document).ready(function(){
    $.getJSON( "/toc.txt", function( data ) {
      tocData = data;
      renderTree(data.toc);
      $(".nav-sub").html(treeOutput.join(''));
      hookupTOCEvents();
    });
    $("#TableOfContents ul").empty();

    var prevH2Item = null;
    var prevH2List = null;

    var index = 0;
    var currentHeader = 0, lastHeader = 0;
    var output = "";
    $("h2, h3, h4").each(function() {
        var li= "<li><a href='" + window.location + "#" + $(this).attr('id') + "'>" + $(this).text().replace("¶","") + "</a></li>";
        lastHeader = currentHeader;
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
        if (currentHeader > lastHeader)
        {
            // nest further
            output += "<ul>" + li;
        } else if (lastHeader < currentHeader)
        {
            // close nesting
            output += "</ul>" + li
        } else {
            // continue, no change in nesting
            output += li;
        }
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
