jQuery(document).ready(function(){
    $('.expand-menu').on('click touchstart', function(elem) {
//      menu = elem.currentTarget.nextElementSibling
      console.log("menu.js firing!")
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

    $("#TableOfContents ul").empty();

    var prevH2Item = null;
    var prevH2List = null;

    var index = 0;
    $("h2, h3").each(function() {

        //insert an anchor to jump to, from the TOC link.
        var anchor = "<a name='" + index + "'></a>";
        $(this).before(anchor);

        var li= "<li><a href='" + window.location + "#" + $(this).id + "'>" + $(this).text() + "</a></li>";

        if( $(this).is("h2") ){
            prevH2List = $("<ul></ul>");
            prevH2Item = $(li);
            prevH2Item.append(prevH2List);
            prevH2Item.appendTo("#TableOfContents ul");
        } else {
            prevH2List.append(li);
        }
        index++;
    });
});
