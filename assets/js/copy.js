// copy code icon markup
const copyIcon = `<svg aria-hidden="true" data-testid="geist-icon" fill="none" height="18" shape-rendering="geometricPrecision" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" viewBox="0 0 24 24" width="18" style="color: currentcolor;"><path d="M8 17.929H6c-1.105 0-2-.912-2-2.036V5.036C4 3.91 4.895 3 6 3h8c1.105 0 2 .911 2 2.036v1.866m-6 .17h8c1.105 0 2 .91 2 2.035v10.857C20 21.09 19.105 22 18 22h-8c-1.105 0-2-.911-2-2.036V9.107c0-1.124.895-2.036 2-2.036z"></path></svg>`

// insert copy buttons for code blocks
const codeBlocks = document.querySelectorAll("div.highlighter-rouge")
codeBlocks.forEach((codeBlock) => {
  codeBlock.insertAdjacentHTML(
    "afterbegin",
    `<button class="copy" aria-label="Copy code to clipboard">${copyIcon}</button>`
  )
})

// handler that saves the code block innerText to clipboard
function copyCodeBlock(event) {
  const copyButton = event.currentTarget
  const codeBlock = copyButton.parentElement.querySelector("pre.highlight code")
  const code = codeBlock.innerText.trim()
  // remove "$ " prompt at start of lines in code
  const strippedCode = code.replace(/^[\s]?\$\s+/gm, "")
  window.navigator.clipboard.writeText(strippedCode)

  // change the button text temporarily
  copyButton.textContent = "Copied!"
  setTimeout(() => copyButton.innerHTML = copyIcon, 3000)
}

// register event listeners for copy buttons
const copyButtons = document.querySelectorAll("button.copy")
copyButtons.forEach((btn) => {
  btn.addEventListener("click", copyCodeBlock)
})
