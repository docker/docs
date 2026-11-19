import "../quick-install.js";
// Native dialog supplies focus trapping, Escape dismissal and focus restoration.
const searchDialog = document.getElementById("docs-search");
function openDocsSearch() {
  window.Alpine?.store("gordon").close();
  document.querySelectorAll("[data-mobile-nav]").forEach((menu) => {
    menu.open = false;
  });
  document
    .querySelectorAll(".image-zoom-dialog[open]")
    .forEach((dialog) => dialog.close());
  searchDialog?.showModal();
  searchDialog?.querySelector("input")?.focus();
}
document
  .querySelectorAll("[data-open-search]")
  .forEach((button) => button.addEventListener("click", openDocsSearch));
document
  .querySelector("[data-close-search]")
  ?.addEventListener("click", () => searchDialog.close());
searchDialog?.addEventListener("click", (event) => {
  if (event.target === searchDialog) searchDialog.close();
});
document.addEventListener("keydown", (event) => {
  if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k") {
    event.preventDefault();
    openDocsSearch();
  }
  if (event.key === "Escape") {
    if (searchDialog?.open) {
      event.preventDefault();
      searchDialog.close();
    }
    document.querySelectorAll("[data-mobile-nav]").forEach((menu) => {
      menu.open = false;
    });
  }
});
// Restore synchronously after layout; reveal only the part of the active item outside the viewport.
const sidebarScroll = document.querySelector("[data-sidebar-scroll]");
if (sidebarScroll) {
  const key = `docs-sidebar:${sidebarScroll.querySelector("[data-sidebar-key]")?.dataset.sidebarKey || location.pathname}`;
  try {
    sidebarScroll.scrollTop = Number(sessionStorage.getItem(key)) || 0;
  } catch {}
  const active = sidebarScroll.querySelector('[aria-current="page"]');
  if (active) {
    const box = sidebarScroll.getBoundingClientRect();
    const item = active.getBoundingClientRect();
    if (item.top < box.top) sidebarScroll.scrollTop -= box.top - item.top;
    else if (item.bottom > box.bottom)
      sidebarScroll.scrollTop += item.bottom - box.bottom;
  }
  sidebarScroll.addEventListener(
    "scroll",
    () => {
      try {
        sessionStorage.setItem(key, String(sidebarScroll.scrollTop));
      } catch {}
    },
    { passive: true },
  );
}

const mobileSidebar = window.matchMedia("(width < 700px)");
function updateSidebar({ matches }) {
  document.querySelectorAll(".sidebar-panel").forEach((panel) => {
    panel.open = !matches;
  });
}
updateSidebar(mobileSidebar);
mobileSidebar.addEventListener("change", updateSidebar);
