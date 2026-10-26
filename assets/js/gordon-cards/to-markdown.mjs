import { GORDON_CARD_FENCE_LANGUAGE } from './components.mjs'
import {
  parseGordonCard,
  readArray,
  readCalloutVariant,
  readCell,
  readElement,
  readString,
} from './parse.mjs'

// Gordon's answer with each ```openui-lang card written out as plain markdown, for where the
// answer is read as text rather than drawn: what the Copy button copies. Follow-up chips are left
// out, as they only make sense in the chat. A block that doesn't parse into a card is kept as is.
export function gordonAnswerToMarkdown(answer) {
  return answer.replace(CARD_FENCE, (block, code) => {
    const root = parseGordonCard(code)
    return root === null ? block : cardToMarkdown(root.props)
  })
}

const CARD_FENCE = new RegExp(
  `^\`\`\`${GORDON_CARD_FENCE_LANGUAGE}[^\\S\\n]*\\n([\\s\\S]*?)^\`\`\`[^\\S\\n]*$`,
  'gm',
)

const CALLOUT_ALERTS = {
  info: 'NOTE',
  success: 'TIP',
  warning: 'WARNING',
  destructive: 'CAUTION',
}

function cardToMarkdown(props) {
  const title = readString(props.title)
  const description = readString(props.description)
  const children = readArray(props.children)
    .map(componentToMarkdown)
    .filter((markdown) => markdown !== undefined)
  return [
    title === undefined ? undefined : `**${title}**`,
    description,
    ...children,
  ]
    .filter((part) => part !== undefined)
    .join('\n\n')
}

function componentToMarkdown(node) {
  const child = readElement(node)
  if (child === undefined) {
    return undefined
  }
  switch (child.typeName) {
    case 'Steps':
      return stepsToMarkdown(child.props)
    case 'Callout':
      return calloutToMarkdown(child.props)
    case 'Table':
      return tableToMarkdown(child.props)
    case 'Code':
      return codeToMarkdown(
        readString(child.props.code) ?? '',
        readString(child.props.language),
      )
    default:
      return undefined
  }
}

function stepsToMarkdown(props) {
  return readArray(props.items)
    .map(readElement)
    .filter((step) => step !== undefined)
    .map(({ props: step }, index) => {
      const detail = readString(step.detail)
      const command = readString(step.command)
      return [
        `${index + 1}. ${readString(step.title) ?? ''}`,
        ...(detail === undefined ? [] : [`   ${detail}`]),
        ...(command === undefined
          ? []
          : [indent(codeToMarkdown(command, 'bash'), '   ')]),
      ].join('\n')
    })
    .join('\n')
}

function calloutToMarkdown(props) {
  const title = readString(props.title)
  const text = readString(props.text)
  return [
    `> [!${CALLOUT_ALERTS[readCalloutVariant(props.variant)]}]`,
    ...(title === undefined ? [] : [`> **${title}**`]),
    ...(text === undefined ? [] : [`> ${text}`]),
  ].join('\n')
}

function tableToMarkdown(props) {
  const columns = readArray(props.columns).map(toTableCell)
  const rows = readArray(props.rows).map((row) =>
    readArray(row).map(toTableCell),
  )
  return [
    toTableRow(columns),
    toTableRow(columns.map(() => '---')),
    ...rows.map(toTableRow),
  ].join('\n')
}

function codeToMarkdown(code, language) {
  return [`\`\`\`${language ?? ''}`, code, '```'].join('\n')
}

function toTableRow(cells) {
  return `| ${cells.join(' | ')} |`
}

function toTableCell(value) {
  return readCell(value).replaceAll('|', '\\|').replaceAll('\n', ' ')
}

function indent(text, prefix) {
  return text
    .split('\n')
    .map((line) => prefix + line)
    .join('\n')
}
