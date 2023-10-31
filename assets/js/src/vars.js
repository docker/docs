const codeblocks = document.querySelectorAll(".highlight")

for (const codeblock of codeblocks) {
  codeblock.innerHTML = codeblock.innerHTML
    .replaceAll(
      /&lt;([A-Z_]+)&gt;/g,
      `<var class="text-violet-light dark:text-violet-dark"
        >$1</var>`
    )
}
