const toc = document.querySelector("#TableOfContents");

if (toc) {
  const prose = document.querySelector("article.prose");
  const headings = prose.querySelectorAll("h2, h3");
  // grab all heading anchors on this page
  const headingAnchors = Array.from(headings)
    .map((h) => h.previousElementSibling)
    .map((a) => `${a.name}`);

  // find the closest anchor link based on window scrollpos
  function findClosestHeading(headingAnchors) {
    const { innerHeight } = window;
    for (const anchor of headingAnchors) {
      const el = document.querySelector(`a[name="${anchor}"]`);
      const { top } = el.getBoundingClientRect();
      // if the heading is visible and within the top 20% of viewport
      if (top > 0 && top < (innerHeight * 0.2)) {
        return anchor;
      }
    }
  }

  function updateToc(closestHeading) {
    const prev = toc.querySelector('a[aria-current="true"]');
    const next = toc.querySelector(`a[href="#${closestHeading}"]`);
    console.log(prev, next)
    if (prev && next) {
      prev.removeAttribute("aria-current");
    }
    if (next) {
      next.setAttribute("aria-current", "true");
    }
  }

  function updateTocHighlight(event) {
    if (event.type === "click" && !toc.contains(event.target)) {
      return;
    }
    // grab the anchor id of the closest heading
    const closestHeading = findClosestHeading(headingAnchors);
    updateToc(closestHeading);
  }

  window.addEventListener("scroll", updateTocHighlight);
  window.addEventListener("click", updateTocHighlight);
}
