document.querySelectorAll("[data-api-version]").forEach((select) =>
  select.addEventListener("change", () => {
    location.href = select.value;
  }),
);
document.querySelectorAll("[data-api-filter]").forEach((input) =>
  input.addEventListener("input", () => {
    const term = input.value.toLocaleLowerCase();
    let visible = 0;
    document.querySelectorAll("[data-api-filter-item]").forEach((row) => {
      row.hidden = !row.textContent.toLocaleLowerCase().includes(term);
      if (!row.hidden) visible++;
    });
    document.querySelector("[data-api-filter-status]").textContent =
      `${visible} ${visible === 1 ? "operation" : "operations"} found.`;
  }),
);
document.querySelectorAll("[data-api-copy]").forEach((button) =>
  button.addEventListener("click", async () => {
    const status = button.parentElement.querySelector("[data-api-copy-status]");
    status.textContent = "";
    try {
      await navigator.clipboard.writeText(
        button.parentElement.querySelector("[data-api-copy-source]")
          .textContent,
      );
      button.textContent = "Copied";
      status.textContent = "Request copied.";
    } catch {
      button.textContent = "Select and copy the request";
      status.textContent =
        "Could not copy. Select and copy the request manually.";
    }
  }),
);
document.querySelectorAll("[data-api-example-select]").forEach((select) => {
  const update = () =>
    select
      .closest("[data-api-examples]")
      .querySelectorAll("[data-api-example]")
      .forEach((example) => {
        example.hidden = example.dataset.apiExample !== select.value;
      });
  select.addEventListener("change", update);
  update();
});

document.querySelectorAll("[data-api-media-select]").forEach((select) =>
  select.addEventListener("change", () => {
    document.querySelectorAll("[data-api-media]").forEach((variant) => {
      variant.hidden = Boolean(
        select.value &&
        variant.dataset.apiMedia &&
        select.value !== variant.dataset.apiMedia,
      );
    });
  }),
);
