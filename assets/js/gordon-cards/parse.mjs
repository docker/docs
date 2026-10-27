import { createParser } from '@openuidev/lang-core'

import { createGordonCardLibrary, GORDON_CARD_ROOT } from './components.mjs'

let parser

// The card in a ```openui-lang block, as an element tree ({ type: 'element', typeName, props }),
// or null if nothing in it parses into a Card (yet, while it streams). Components that aren't in
// components.mjs, unresolved references and unused statements are left out.
export function parseGordonCard(code) {
  parser ??= createParser(
    createGordonCardLibrary().toJSONSchema(),
    GORDON_CARD_ROOT,
  )
  return parser.parse(code).root
}

// The props of a parsed child component, or undefined if `value` isn't one (yet). Props are
// model output, possibly half-streamed, so nothing about their shape can be assumed.
export function readElement(value) {
  if (
    typeof value !== 'object' ||
    value === null ||
    value.type !== 'element' ||
    typeof value.props !== 'object' ||
    value.props === null
  ) {
    return undefined
  }
  return { typeName: value.typeName, props: value.props }
}

export function readArray(value) {
  return Array.isArray(value) ? value : []
}

export function readString(value) {
  return typeof value === 'string' ? value : undefined
}

export function readCell(value) {
  return typeof value === 'string' || typeof value === 'number'
    ? String(value)
    : ''
}

export const CALLOUT_VARIANTS = ['info', 'success', 'warning', 'destructive']

export function readCalloutVariant(value) {
  return CALLOUT_VARIANTS.includes(value) ? value : 'info'
}
