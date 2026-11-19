// One clipboard interaction for fenced examples and platform-specific install commands.
const resets = new WeakMap();
document.addEventListener("click", async (event) => {
  const button = event.target.closest("[data-copy-text]");
  if (!button || button.disabled || button.getAttribute("aria-busy") === "true")
    return;
  const control = button.closest(".code-copy");
  const status = control.querySelector('[role="status"]');
  const icon = control.querySelector("[data-copy-icon]");
  const check = control.querySelector("[data-copied-icon]");
  clearTimeout(resets.get(button));
  button.setAttribute("aria-busy", "true");
  status.textContent = "";
  status.classList.add("sr-only");
  try {
    await navigator.clipboard.writeText(button.dataset.copyText);
    icon.hidden = true;
    check.hidden = false;
    status.textContent = "Code copied.";
    resets.set(
      button,
      setTimeout(() => {
        icon.hidden = false;
        check.hidden = true;
      }, 2000),
    );
  } catch {
    icon.hidden = false;
    check.hidden = true;
    status.classList.remove("sr-only");
    status.textContent = "Could not copy. Select and copy the code manually.";
  } finally {
    button.removeAttribute("aria-busy");
  }
});
