jQuery(document).ready(function(){
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
});
