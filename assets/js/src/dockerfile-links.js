const keywords = [
  "ADD",
  "ARG",
  "CMD",
  "COPY",
  "ENTRYPOINT",
  "ENV",
  "EXPOSE",
  "FROM",
  "HEALTHCHECK",
  "LABEL",
  // "MAINTAINER",
  "ONBUILD",
  "RUN",
  "SHELL",
  "STOPSIGNAL",
  "USER",
  "VOLUME",
  "WORKDIR",
]
const cmds = Array.from(document.querySelectorAll(".language-dockerfile span.k"))
  .filter((el) => keywords.some(kwd => el.textContent.includes(kwd)));

for (const cmd of cmds) {
  const name = cmd.textContent;
  const a = document.createElement("a")
  a.classList.add("underline","underline-offset-4","decoration-dashed","cursor-pointer")
  a.title = `Learn more about the ${name} instruction`
  a.href = `/reference/dockerfile/#${name.toLowerCase()}`
  a.innerHTML = cmd.outerHTML
  cmd.insertAdjacentElement("beforebegin", a)
  cmd.remove()
}
