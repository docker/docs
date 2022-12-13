// insert copy buttons for code blocks
const codeBlocks = document.querySelectorAll("div.highlighter-rouge")
codeBlocks.forEach((codeBlock) => {
  const title = codeBlock.getAttribute("title")
  codeBlock.insertAdjacentHTML(
    "afterbegin",
    `
    <header>
      <div class="hl-title"><span>${title ? `${title}` : ""}</span></div>
      <button class="hl-copy" aria-label="Copy code to clipboard"><div title="copy icon" class="copy-icon"></button>
    </header>
    `
  )
})

// handler that saves the code block innerText to clipboard
function copyCodeBlock(event) {
  const copyButton = event.currentTarget
  const codeBlock = copyButton
    .closest(".highlighter-rouge")
    .querySelector("pre.highlight code")
  const code = codeBlock.innerText.trim()
  // remove "$ " prompt at start of lines in code
  const strippedCode = code.replace(/^[\s]?(\$|PS>)\s+/gm, "")
  window.navigator.clipboard.writeText(strippedCode)

  // change the button text temporarily
  copyButton.classList.add("copied")
  setTimeout(() => copyButton.classList.remove("copied"), 3000)
}

// register event listeners for copy buttons
const copyButtons = document.querySelectorAll("button.hl-copy")
copyButtons.forEach((btn) => {
  btn.addEventListener("click", copyCodeBlock)
})
