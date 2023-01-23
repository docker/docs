// Scroll the given menu item into view. We actually pick the item *above*
// the current item to give some headroom above
function scrollMenuItem() {
  const sectiontree = document.querySelector("#sectiontree")
  const sidebar = document.querySelector("#sidebar")
  let item = sectiontree.querySelector('[aria-current="page"]')
  if (!item) return
  item = item.parentElement.closest("li")
  if (item) {
    const itemY = item.getBoundingClientRect().y
    // scroll to the item y-coord (with a 150px padding for some head room)
    if (itemY > window.innerHeight - 150) {
      sidebar.scrollTop = itemY - 150
    }
  }
}

scrollMenuItem()
