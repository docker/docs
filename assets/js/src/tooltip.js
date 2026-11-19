import { computePosition, flip, shift, offset, arrow } from "@floating-ui/dom";

/* Regular tooltips (partial) */

const tooltipWrappers = Array.from(
  document.querySelectorAll("[data-tooltip-wrapper]"),
);

for (const [index, tooltipWrapper] of tooltipWrappers.entries()) {
  const button = tooltipWrapper.firstElementChild;
  const tooltip = button.nextElementSibling;
  const arrowElement = tooltip.firstElementChild;
  tooltip.id ||= `docs-tooltip-${index}`;
  button.setAttribute("aria-describedby", tooltip.id);

  function update() {
    computePosition(button, tooltip, {
      placement: "top",
      middleware: [
        offset(6),
        flip(),
        shift({ padding: 5 }),
        arrow({ element: arrowElement }),
      ],
    }).then(({ x, y, placement, middlewareData }) => {
      tooltip.dataset.placement = placement.split("-")[0];
      Object.assign(tooltip.style, {
        left: `${x}px`,
        top: `${y}px`,
      });

      // Accessing the data
      const { x: arrowX, y: arrowY } = middlewareData.arrow;

      const staticSide = {
        top: "bottom",
        right: "left",
        bottom: "top",
        left: "right",
      }[placement.split("-")[0]];

      Object.assign(arrowElement.style, {
        left: arrowX != null ? `${arrowX}px` : "",
        top: arrowY != null ? `${arrowY}px` : "",
        right: "",
        bottom: "",
        [staticSide]: "-4px",
      });
    });
  }

  function showTooltip() {
    if (button.disabled) return;
    tooltip.classList.remove("hidden");
    document.addEventListener("keydown", dismissOnEscape);
    update();
  }

  function hideTooltip() {
    tooltip.classList.add("hidden");
    document.removeEventListener("keydown", dismissOnEscape);
  }

  function dismissOnEscape(event) {
    if (event.key === "Escape") {
      if (tooltipWrapper.contains(event.target)) event.stopPropagation();
      hideTooltip();
    }
  }

  tooltipWrapper.addEventListener("mouseenter", showTooltip);
  tooltipWrapper.addEventListener("mouseleave", () => {
    if (!tooltipWrapper.contains(document.activeElement)) hideTooltip();
  });
  button.addEventListener("focus", showTooltip);
  button.addEventListener("blur", () => {
    if (!tooltipWrapper.matches(":hover")) hideTooltip();
  });
}
