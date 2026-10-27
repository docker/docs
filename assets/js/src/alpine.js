import Alpine from 'alpinejs'
import collapse from '@alpinejs/collapse'
import persist from '@alpinejs/persist'
import focus from '@alpinejs/focus'
import { marked } from 'marked'
import hljs from 'highlight.js/lib/core'
// Import languages relevant to Docker docs
import bash from 'highlight.js/lib/languages/bash'
import dockerfile from 'highlight.js/lib/languages/dockerfile'
import yaml from 'highlight.js/lib/languages/yaml'
import json from 'highlight.js/lib/languages/json'
import javascript from 'highlight.js/lib/languages/javascript'
import python from 'highlight.js/lib/languages/python'
import go from 'highlight.js/lib/languages/go'

window.Alpine = Alpine

Alpine.plugin(collapse)
Alpine.plugin(persist)
Alpine.plugin(focus)

// Register highlight.js languages
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('sh', bash)
hljs.registerLanguage('shell', bash)
hljs.registerLanguage('console', bash)
hljs.registerLanguage('dockerfile', dockerfile)
hljs.registerLanguage('yaml', yaml)
hljs.registerLanguage('yml', yaml)
hljs.registerLanguage('json', json)
hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('js', javascript)
hljs.registerLanguage('python', python)
hljs.registerLanguage('py', python)
hljs.registerLanguage('go', go)
hljs.registerLanguage('golang', go)

// Configure marked to escape HTML in text tokens only (not code blocks)
marked.use({
  walkTokens(token) {
    // Escape HTML in text and HTML tokens, preserve code blocks
    if (token.type === 'text' || token.type === 'html') {
      const text = token.text || token.raw
      const escaped = text
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
      if (token.text) token.text = escaped
      if (token.raw) token.raw = escaped
    }
  }
})

// Gordon's answers can contain ```openui-lang cards. The code that draws them is a separate
// bundle (assets/js/gordon-cards), loaded the first time a card shows up.
Alpine.store('gordonCards', { ready: false, failed: false })

let gordonCardsLoading

function loadGordonCards() {
  gordonCardsLoading ??= new Promise((resolve, reject) => {
    const script = document.createElement('script')
    script.src = window.GORDON_CARDS_SRC
    script.onload = () => {
      Alpine.store('gordonCards').ready = true
      resolve(window.GordonCards)
    }
    script.onerror = () => {
      // Cards then stay plain code blocks, but a later call (such as Copy) can try again.
      gordonCardsLoading = undefined
      script.remove()
      Alpine.store('gordonCards').failed = true
      reject(new Error('Failed to load gordon-cards.js'))
    }
    document.head.append(script)
  })
  return gordonCardsLoading
}

Alpine.magic('gordonCards', () => loadGordonCards)

// Swaps each ```openui-lang block for the card it describes. Until the card code has loaded, or
// while a block is still streaming in and doesn't parse yet, it shows a placeholder. If the card
// code can't load, or a finished block isn't a card, the block stays a plain code block.
function replaceGordonCards(div, { streaming }) {
  const blocks = div.querySelectorAll('pre > code.language-openui-lang')
  if (blocks.length === 0) return

  const cards = Alpine.store('gordonCards')
  if (!cards.ready && !cards.failed) {
    loadGordonCards().catch((err) => console.error('Gordon cards:', err))
  }
  blocks.forEach((code) => {
    const pre = code.parentElement
    if (cards.ready) {
      const card = window.GordonCards.renderGordonCard(code.textContent)
      if (card) {
        pre.replaceWith(card)
      } else if (streaming) {
        pre.replaceWith(gordonCardPlaceholder())
      }
    } else if (!cards.failed) {
      pre.replaceWith(gordonCardPlaceholder())
    }
  })
}

function gordonCardPlaceholder() {
  const placeholder = document.createElement('div')
  placeholder.className = 'gordon-card-loading not-prose'
  placeholder.setAttribute('role', 'status')
  placeholder.setAttribute('aria-label', 'Drawing card')
  placeholder.append(document.createElement('span'), document.createElement('span'))
  return placeholder
}

// Add $markdown magic for rendering markdown with syntax highlighting
Alpine.magic('markdown', () => {
  return (content, { streaming = false } = {}) => {
    if (!content) return ''
    const html = marked(content)

    // Parse and highlight code blocks
    const div = document.createElement('div')
    div.innerHTML = html

    replaceGordonCards(div, { streaming })

    // Handle code blocks (pre > code)
    div.querySelectorAll('pre').forEach((pre) => {
      // Add not-prose to prevent Tailwind Typography styling
      pre.classList.add('not-prose')
      const code = pre.querySelector('code')
      if (code) {
        // Preserve the original text with newlines
        const codeText = code.textContent

        // Clear and set as plain text first to preserve structure
        code.textContent = codeText

        // Now apply highlight.js which will work with the text nodes
        hljs.highlightElement(code)
      }
    })

    // Handle inline code elements (not in pre blocks)
    div.querySelectorAll('code:not(pre code)').forEach((code) => {
      code.classList.add('not-prose')
    })

    return div.innerHTML
  }
})

// Stores
Alpine.store("showSidebar", false)
Alpine.store('gordon', {
  isOpen: false,
  query: '',
  autoSubmit: false,
  toggle() {
    this.isOpen = !this.isOpen
  },
  open(query, autoSubmit = false) {
    this.isOpen = true
    if (query) {
      this.query = query
      this.autoSubmit = autoSubmit
    }
  },
  close() {
    this.isOpen = false
  }
})

Alpine.start()
