const toc = document.querySelector("#TableOfContents .toc");
const headings = document.querySelectorAll("main article h2, main article h3");

if (toc) {
  window.addEventListener("scroll", updateToc);
  updateToc();
}

function updateToc() {
  // find the heading currently in view
  const currentSection = Array.from(headings).reduce((previous, current) => {
    const { top } = current.getBoundingClientRect();
    // if current is in the top 100px of the viewport
    // 100px is an arbitrary value
    // Should be header height + margin
    if (top < 100) {
      return current;
    }
    return previous
  });
  const prev = toc.querySelector('a[aria-current="true"]');
  const next = toc.querySelector(`a[href="#${currentSection.id}"]`);
  if (prev) {
    prev.removeAttribute("aria-current");
  }
  if (next) {
    next.setAttribute("aria-current", "true");
  }
}
