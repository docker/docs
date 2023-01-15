function handleSidebarClick(e) {
  if (!e.target.classList.contains("sidebar-expander")) return
  const sectionElement = e.target.closest("li")
  if (sectionElement.dataset.expanded === "true") {
    e.target.textContent = "chevron_right"
    sectionElement.dataset.expanded = "false"
  } else {
    e.target.textContent = "expand_more"
    sectionElement.dataset.expanded = "true"
  }
  sectionElement.nextElementSibling.classList.toggle("hidden")
}

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

const sectiontree = document.querySelector("#sectiontree")
if (sectiontree) {
  scrollMenuItem()
  sectiontree.addEventListener("click", handleSidebarClick)
}
