document.querySelectorAll("[data-api-version]").forEach((select) =>
  select.addEventListener("change", () => {
    location.href = select.value;
  }),
);
document.querySelectorAll("[data-api-filter]").forEach((input) =>
  input.addEventListener("input", () => {
    const term = input.value.toLocaleLowerCase();
    document.querySelectorAll("[data-api-filter-item]").forEach((row) => {
      row.hidden = !row.textContent.toLocaleLowerCase().includes(term);
    });
  }),
);
document.querySelectorAll("[data-api-copy]").forEach((button) =>
  button.addEventListener("click", async () => {
    try {
      await navigator.clipboard.writeText(
        button.parentElement.querySelector("[data-api-copy-source]")
          .textContent,
      );
      button.textContent = "Copied";
    } catch {
      button.textContent = "Select and copy the request";
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
