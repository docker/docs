function toggleMenuItem(event) {
  const section = event.currentTarget.parentElement;
  const icons = event.currentTarget.querySelectorAll(".icon-svg");
  const subsection = section.querySelector("ul");
  subsection.classList.toggle("hidden");
  icons.forEach(i => i.classList.toggle('hidden'))
}

const sectiontree = document.querySelector("#sectiontree");
if (sectiontree) {
  for (const button of sectiontree.querySelectorAll("button")) {
    button.addEventListener("click", toggleMenuItem);
  }
}
