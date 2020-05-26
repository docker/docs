function hookupTOCEvents() {
  $('.expand-menu').on('click', function(elem) {
    var menu = elem.currentTarget.parentElement
    if (menu.classList.contains("menu-closed")) {
      menu.classList.remove("menu-closed")
      menu.classList.add("menu-open")
    } else {
      menu.classList.add("menu-closed")
      menu.classList.remove("menu-open")
    }
    return false;
  });

  $(".left-off-canvas-menu").css("display","block");
}
$(document).ready(hookupTOCEvents);
