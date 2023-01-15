const toc = document.querySelector("#TableOfContents")

if (toc) {
  const prose = document.querySelector("article.prose")
  const headings = prose.querySelectorAll("h2, h3, h4")
  // grab the yposition and targets for all anchor links
  const anchorLinks = Array.from(headings)
    .map((h) => h.previousElementSibling)
    .map((a) => [a.offsetTop - 64, `#${a.name}`])
    .sort((a, b) => (a[0] > b[0] ? 1 : 0))

  let closestHeading

  // find the closest anchor link based on window scrollpos
  function findClosestHeading() {
    yPos = window.scrollY
    closest = anchorLinks.reduce((prev, curr) => {
      return Math.abs(curr[0] - yPos) < Math.abs(prev[0] - yPos) &&
        curr[0] <= yPos
        ? curr
        : prev
    })
    return closest
  }

  function updateToc() {
    const prev = toc.querySelector('a[aria-current="true"]')
    const next = toc.querySelector(`a[href="${closestHeading[1]}"]`)
    if (prev) {
      prev.removeAttribute("aria-current")
    }
    if (next) {
      next.setAttribute("aria-current", "true")
    }
  }

  function handleScroll() {
    closestHeading = findClosestHeading()
    updateToc()
  }

  window.addEventListener("scroll", handleScroll)
}
