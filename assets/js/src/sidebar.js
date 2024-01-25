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
if (sectiontree) {
  for (const button of sectiontree.querySelectorAll("button")) {
    button.addEventListener("click", toggleMenuItem);
  }
}
