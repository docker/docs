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

// Add $markdown magic for rendering markdown with syntax highlighting
Alpine.magic('markdown', () => {
  return (content) => {
    if (!content) return ''
    const html = marked(content)

    // Parse and highlight code blocks
    const div = document.createElement('div')
    div.innerHTML = html

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
  isOpen: Alpine.$persist(false).using(sessionStorage).as('gordon-isOpen'),
  query: '',
  toggle() {
    this.isOpen = !this.isOpen
  },
  open(query) {
    this.isOpen = true
    if (query) this.query = query
  },
  close() {
    this.isOpen = false
  }
})

Alpine.start()
