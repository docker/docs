const toc = document.querySelector("#TableOfContents");

if (toc) {
  const prose = document.querySelector("article.prose");
  const headings = prose.querySelectorAll("h2, h3");
  // grab all heading anchors on this page
  const anchorLinks = Array.from(headings)
    .map((h) => h.previousElementSibling)
    .map((a) => [a.offsetTop, `${a.name}`])
    .sort((a, b) => (a[0] > b[0] ? 1 : 0));

  // find the closest anchor link based on window scrollpos
  function findClosestHeading(anchorLinks) {
    const { innerHeight } = window;
    for (const anchor of anchorLinks) {
      const el = document.querySelector(`a[name="${anchor[1]}"]`);
      const { top } = el.getBoundingClientRect();
      // if the heading is visible and within the top 20% of viewport
      if (top > 0 && top < (innerHeight * 0.2)) {
        return anchor[1];
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
    const closestHeading = findClosestHeading(anchorLinks);
    updateToc(closestHeading);
  }

  window.addEventListener("scroll", updateTocHighlight);
  window.addEventListener("click", updateTocHighlight);
}
