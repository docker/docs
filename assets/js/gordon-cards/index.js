// Entry point of the gordon-cards.js bundle. It's built apart from scripts.js and only loaded by
// the Gordon chat when an answer first contains a card (see loadGordonCards in src/alpine.js), so
// the OpenUI parser doesn't weigh on every docs page.

import { renderGordonCard } from './render.mjs'
import { gordonAnswerToMarkdown } from './to-markdown.mjs'

window.GordonCards = { renderGordonCard, gordonAnswerToMarkdown }
