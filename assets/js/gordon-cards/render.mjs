import {
  parseGordonCard,
  readArray,
  readCalloutVariant,
  readCell,
  readElement,
  readString,
} from './parse.mjs'

// Draws a ```openui-lang block from Gordon's answer as a card element, or returns null if nothing
// in it parses into a card (yet). Text from the model only ever goes in through textContent and
// setAttribute, never as HTML. Code blocks come out as <pre><code> for the chat's usual syntax
// highlighting, and follow-up chips as buttons carrying data-gordon-follow-up, which the chat
// turns into the user's next question.
export function renderGordonCard(code) {
  const root = parseGordonCard(code)
  return root === null ? null : renderCard(root.props)
}

function renderCard(props) {
  const card = element('section', 'gordon-card not-prose')
  const title = readString(props.title)
  const description = readString(props.description)
  if (title !== undefined) {
    card.append(element('h4', 'gordon-card-title', title))
  }
  if (description !== undefined) {
    card.append(element('p', 'gordon-card-description', description))
  }
  for (const child of readArray(props.children)) {
    const rendered = renderChild(child)
    if (rendered !== undefined) {
      card.append(rendered)
    }
  }
  return card
}

function renderChild(node) {
  const child = readElement(node)
  if (child === undefined) {
    return undefined
  }
  switch (child.typeName) {
    case 'Steps':
      return renderSteps(child.props)
    case 'Callout':
      return renderCallout(child.props)
    case 'Table':
      return renderTable(child.props)
    case 'Code':
      return renderCode(
        readString(child.props.code) ?? '',
        readString(child.props.language),
      )
    case 'FollowUps':
      return renderFollowUps(child.props)
    default:
      return undefined
  }
}

function renderSteps(props) {
  const list = element('ol', 'gordon-card-steps')
  for (const item of readArray(props.items)) {
    const step = readElement(item)
    if (step === undefined) {
      continue
    }
    const title = readString(step.props.title)
    const detail = readString(step.props.detail)
    const command = readString(step.props.command)
    const entry = element('li', 'gordon-card-step')
    if (title !== undefined) {
      entry.append(element('p', 'gordon-card-step-title', title))
    }
    if (detail !== undefined) {
      entry.append(element('p', 'gordon-card-step-detail', detail))
    }
    if (command !== undefined) {
      entry.append(renderCode(command, 'bash'))
    }
    list.append(entry)
  }
  return list
}

function renderCallout(props) {
  const variant = readCalloutVariant(props.variant)
  const callout = element(
    'div',
    `gordon-card-callout gordon-card-callout-${variant}`,
  )
  const title = readString(props.title)
  const text = readString(props.text)
  if (title !== undefined) {
    callout.append(element('p', 'gordon-card-callout-title', title))
  }
  if (text !== undefined) {
    callout.append(element('p', 'gordon-card-callout-text', text))
  }
  return callout
}

function renderTable(props) {
  const wrapper = element('div', 'gordon-card-table')
  const table = element('table')
  const head = element('thead')
  const headRow = element('tr')
  for (const column of readArray(props.columns)) {
    headRow.append(element('th', undefined, readCell(column)))
  }
  head.append(headRow)
  const body = element('tbody')
  for (const row of readArray(props.rows)) {
    const bodyRow = element('tr')
    for (const cell of readArray(row)) {
      bodyRow.append(element('td', undefined, readCell(cell)))
    }
    body.append(bodyRow)
  }
  table.append(head, body)
  wrapper.append(table)
  return wrapper
}

function renderCode(code, language) {
  const pre = element('pre', 'not-prose')
  const codeElement = element('code', undefined, code)
  if (language !== undefined && /^[\w+-]+$/.test(language)) {
    codeElement.className = `language-${language}`
  }
  pre.append(codeElement)
  return pre
}

function renderFollowUps(props) {
  const list = element('div', 'gordon-card-follow-ups')
  const questions = readArray(props.items)
    .map(readString)
    .filter((question) => question !== undefined && question.trim() !== '')
    .slice(0, 3)
  for (const question of questions) {
    const chip = element('button', 'gordon-card-follow-up', question)
    chip.type = 'button'
    chip.setAttribute('data-gordon-follow-up', question)
    list.append(chip)
  }
  return list
}

function element(tag, className, text) {
  const node = document.createElement(tag)
  if (className !== undefined) {
    node.className = className
  }
  if (text !== undefined) {
    node.textContent = text
  }
  return node
}
