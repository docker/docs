const cmds = document.querySelectorAll(".language-dockerfile span.k");

for (const cmd of cmds) {
  const name = cmd.textContent;
  const a = document.createElement("a")
  a.classList.add("underline","underline-offset-4","decoration-dashed","cursor-pointer")
  a.title = `Learn more about the ${name} instruction`
  a.href = `/engine/reference/builder/#${name.toLowerCase()}`
  a.innerHTML = cmd.outerHTML
  cmd.insertAdjacentElement("beforebegin", a)
  cmd.remove()
}
