// copy code icon markup
const copyIcon = `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
<mask id="mask0_1191_12724" style="mask-type:alpha" maskUnits="userSpaceOnUse" x="0" y="0" width="24" height="24">
<g clip-path="url(#clip0_1191_12724)">
<path d="M15 1H4C2.9 1 2 1.9 2 3V16C2 16.55 2.45 17 3 17C3.55 17 4 16.55 4 16V4C4 3.45 4.45 3 5 3H15C15.55 3 16 2.55 16 2C16 1.45 15.55 1 15 1ZM15.59 5.59L20.42 10.42C20.79 10.79 21 11.3 21 11.83V21C21 22.1 20.1 23 19 23H7.99C6.89 23 6 22.1 6 21L6.01 7C6.01 5.9 6.9 5 8 5H14.17C14.7 5 15.21 5.21 15.59 5.59ZM15 12H19.5L14 6.5V11C14 11.55 14.45 12 15 12Z" fill="black"/>
</g>
</mask>
<g mask="url(#mask0_1191_12724)">
<rect width="24" height="24" fill="currentColor"/>
</g>
<defs>
<clipPath id="clip0_1191_12724">
<rect width="24" height="24" fill="white"/>
</clipPath>
</defs>
</svg>`

const copiedIcon = `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
<mask id="mask0_1191_12493" style="mask-type:alpha" maskUnits="userSpaceOnUse" x="0" y="0" width="24" height="24">
<g clip-path="url(#clip0_1191_12493)">
<path d="M19 3H5C3.9 3 3 3.9 3 5V19C3 20.1 3.9 21 5 21H19C20.1 21 21 20.1 21 19V5C21 3.9 20.1 3 19 3ZM10.71 16.29C10.32 16.68 9.69 16.68 9.3 16.29L5.71 12.7C5.32 12.31 5.32 11.68 5.71 11.29C6.1 10.9 6.73 10.9 7.12 11.29L10 14.17L16.88 7.29C17.27 6.9 17.9 6.9 18.29 7.29C18.68 7.68 18.68 8.31 18.29 8.7L10.71 16.29Z" fill="black"/>
</g>
</mask>
<g mask="url(#mask0_1191_12493)">
<rect width="24" height="24" fill="currentColor"/>
</g>
<defs>
<clipPath id="clip0_1191_12493">
<rect width="24" height="24" fill="white"/>
</clipPath>
</defs>
</svg>`

// insert copy buttons for code blocks
const codeBlocks = document.querySelectorAll("div.highlighter-rouge")
codeBlocks.forEach((codeBlock) => {
  const title = codeBlock.getAttribute("title")
  codeBlock.insertAdjacentHTML(
    "afterbegin",
    `
    <header>
      <div class="hl-title"><span>${title ? `${title}` : ''}</span></div>
      <button class="hl-copy" aria-label="Copy code to clipboard"><span>Copy</span>${copyIcon}</button>
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
  copyButton.innerHTML = `<span>Copy</span>${copiedIcon}`
  setTimeout(() => (copyButton.innerHTML = `<span>Copy</span>${copyIcon}`), 3000)
}

// register event listeners for copy buttons
const copyButtons = document.querySelectorAll("button.hl-copy")
copyButtons.forEach((btn) => {
  btn.addEventListener("click", copyCodeBlock)
})
