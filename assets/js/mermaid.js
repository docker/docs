import mermaid from 'mermaid'

const isDark = () => document.documentElement.classList.contains('dark')

mermaid.initialize({
  startOnLoad: false,
  securityLevel: 'strict',
  theme: isDark() ? 'dark' : 'default',
})

const render = () => {
  mermaid.run({ querySelector: 'pre.mermaid' })
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', render)
} else {
  render()
}
