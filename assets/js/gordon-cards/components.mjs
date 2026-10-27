// The components Gordon may draw in a ```openui-lang card in a docs chat answer (OpenUI Lang, see
// https://www.openui.com/docs/openui-lang): each one's name, what the model is told it's for, and
// its props. render.js draws each of them. Gordon's backend tells the model about these same
// components, so when adding or changing one here, the backend's card prompt must change with it:
// a component the model writes but the chat can't draw is dropped, content and all.
//
// OpenUI Lang passes props positionally, in each schema's key order. Reordering or removing a key
// changes the meaning of cards already in a conversation: only ever append optional keys.
//
// Children come before the parents that reference them: OpenUI names a parent's child types by
// looking their schemas up among the components registered with it.

import { createLibrary, defineComponent } from '@openuidev/lang-core'
import { z } from 'zod/v4'

// The fence language that marks a markdown code block in Gordon's answer as a card.
export const GORDON_CARD_FENCE_LANGUAGE = 'openui-lang'

// Sent to the Gordon backend as `ui` in the request context to say the chat can draw cards. The
// backend only teaches the model to write them when it's present.
export const GORDON_CARD_UI_CAPABILITY = 'openui-v1'

// Every card's `root` is a `Card`.
export const GORDON_CARD_ROOT = 'Card'

const tableProps = z.object({
  columns: z.array(z.string()),
  rows: z.array(z.array(z.union([z.string(), z.number()]))),
})

const stepProps = z.object({
  title: z.string(),
  detail: z.string().optional(),
  command: z.string().optional(),
})

const stepsProps = z.object({
  items: z.array(stepProps),
})

const calloutProps = z.object({
  variant: z.enum(['info', 'success', 'warning', 'destructive']),
  title: z.string(),
  text: z.string().optional(),
})

const codeProps = z.object({
  code: z.string(),
  language: z.string().optional(),
})

const followUpsProps = z.object({
  items: z.array(z.string()),
})

const cardProps = z.object({
  title: z.string(),
  children: z.array(
    z.union([stepsProps, calloutProps, tableProps, codeProps, followUpsProps]),
  ),
  description: z.string().optional(),
})

export const gordonCardComponents = [
  {
    name: 'Table',
    description:
      'A table. rows is a list of rows; each row has one cell per column, in column order. A cell is a string or a number.',
    props: tableProps,
  },
  {
    name: 'Step',
    description:
      'One step of a procedure: a short title, an optional detail sentence, and an optional shell command to run.',
    props: stepProps,
  },
  {
    name: 'Steps',
    description: 'A numbered procedure.',
    props: stepsProps,
  },
  {
    name: 'Callout',
    description:
      'A highlighted note: a cause, a warning, an active incident or a result worth noticing.',
    props: calloutProps,
  },
  {
    name: 'Code',
    description:
      'A code or config snippet, e.g. a Dockerfile, a compose.yaml or a daemon.json excerpt.',
    props: codeProps,
  },
  {
    name: 'FollowUps',
    description:
      'Up to 3 short follow-up questions the user is likely to ask next, as clickable chips. Last child of the Card.',
    props: followUpsProps,
  },
  {
    name: 'Card',
    description:
      'The root of every block: a titled container for the other components, stacked vertically.',
    props: cardProps,
  },
]

// The OpenUI library for the components above. Only their schemas matter: render.js draws the
// parsed tree itself, so no component gets a renderer here.
export function createGordonCardLibrary() {
  return createLibrary({
    root: GORDON_CARD_ROOT,
    components: gordonCardComponents.map((component) =>
      defineComponent({ ...component, component: renderNothing }),
    ),
  })
}

function renderNothing() {
  return null
}
