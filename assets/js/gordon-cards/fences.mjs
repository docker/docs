// Markdown only opens a code block at the start of a line, but the model sometimes writes a card's
// opening fence right after a sentence. The card then shows as raw text, and its closing fence
// opens a code block that swallows the rest of the answer. Kept apart from the card bundle, as the
// chat renders every answer through it.
export function separateCardFences(text) {
  return text.replace(/([^\s`])[ \t]*(```openui-lang)/g, '$1\n\n$2')
}
