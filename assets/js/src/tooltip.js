import { computePosition, flip, shift, offset, arrow } from "@floating-ui/dom";

/* Regular tooltips (partial) */

const tooltipWrappers = Array.from(
  document.querySelectorAll("[data-tooltip-wrapper]"),
);

for (const tooltipWrapper of tooltipWrappers) {
  const button = tooltipWrapper.firstElementChild;
  const tooltip = button.nextElementSibling;
  const arrowElement = tooltip.firstElementChild;

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
    tooltip.classList.toggle("hidden");
    update();
  }

  function hideTooltip() {
    tooltip.classList.toggle("hidden");
  }

  [
    ["mouseenter", showTooltip],
    ["mouseleave", hideTooltip],
    ["focus", showTooltip],
    ["blur", hideTooltip],
  ].forEach(([event, listener]) => {
    button.addEventListener(event, listener);
  });
}
