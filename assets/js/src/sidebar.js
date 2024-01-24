// Scroll the given menu item into view. We actually pick the item *above*
// the current item to give some headroom above
function scrollMenuItem() {
  let item = sectiontree.querySelector('[aria-current="page"]');
  if (!item) return;
  item = item.parentElement.closest("li");
  if (item) {
    const itemY = item.getBoundingClientRect().y;
    // scroll to the item y-coord (with a 150px padding for some head room)
    if (itemY > window.innerHeight - 150) {
      sidebar.scrollTop = itemY - 150;
    }
  }
}

function toggleMenuItem(event) {
  const section = event.currentTarget.parentElement;
  const icon = section.querySelector(".icon");
  const subsection = section.querySelector("ul");
  subsection.classList.toggle("hidden");
  if (subsection.classList.contains("hidden")) {
    icon.setAttribute("data-icon", "expand_more");
  } else {
    icon.setAttribute("data-icon", "expand_less");
  }
}

const sectiontree = document.querySelector("#sectiontree");
const sidebar = document.querySelector("#sidebar");
if (sectiontree && sidebar) {
  scrollMenuItem();
  for (const button of sectiontree.querySelectorAll("button")) {
    button.addEventListener("click", toggleMenuItem);
  }
}
