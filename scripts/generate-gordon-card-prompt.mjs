// Maintainer tool, not part of the site: Hugo doesn't bundle scripts/, and the chat never sends
// this prompt. The chat only sends `ui: "openui-v1"`, and Gordon's backend then adds the card
// instructions from its own copy of this output (ui_cards_docs.jinja2).
//
// It lives here because those instructions must describe exactly the components the chat can
// draw (assets/js/gordon-cards/components.mjs): a component the model is told about but the chat
// can't draw is dropped, content and all. After changing the components, run this and copy the
// output into the backend's ui_cards_docs.jinja2, between its {% raw %} tags.
//
// Usage: node scripts/generate-gordon-card-prompt.mjs > prompt.txt

import { createGordonCardLibrary } from '../assets/js/gordon-cards/components.mjs'

const promptOptions = {
  preamble: [
    'Besides markdown, you can draw a UI card inside your answer by writing openui-lang, a declarative UI language, in a fenced code block tagged `openui-lang`. Text outside the fence is ordinary markdown and is shown as usual. Inside the fence, write only openui-lang code: every block is a complete program on its own.',
    '',
    'When your answer gives a procedure of 3 or more steps, compares several values (limits, plans, options, flags), or reports an active incident together with what to do about it, you MUST present that part in an openui-lang card, never as a markdown list or table. An explanation of a concept, a single fact, a greeting, a refusal or a pointer to support stays plain markdown.',
  ].join('\n'),
  examples: [
    `Example — for "How do I containerize my Node.js app?":

\`docker init\` generates the files you need, then Compose builds and runs the app.

\`\`\`openui-lang
root = Card("Containerize a Node.js app", [steps, next])
steps = Steps([init, review, run])
init = Step("Generate the Docker assets", "Run it in your project directory and answer the prompts.", "docker init")
review = Step("Review the generated files", "It creates a Dockerfile, a compose.yaml and a .dockerignore.")
run = Step("Build and run the app", null, "docker compose up --build")
next = FollowUps(["How do I add a database service?", "How do I reduce my image size?"])
\`\`\`

Sources:
- [https://docs.docker.com/reference/cli/docker/init/](https://docs.docker.com/reference/cli/docker/init/)`,
    `Example — for "What are the Docker Hub pull limits?":

Pull limits depend on whether you're signed in and on your plan.

\`\`\`openui-lang
root = Card("Docker Hub pull limits", [table, tip, next], "Per 6 hours")
table = Table(["Account", "Pull limit"], [["Unauthenticated", "100 per IP address"], ["Personal (signed in)", 200], ["Pro, Team, Business", "Unlimited"]])
tip = Callout("success", "Sign in to raise your limit", "Run docker login before pulling, including in CI.")
next = FollowUps(["How do I authenticate in GitHub Actions?"])
\`\`\``,
  ],
  additionalRules: [
    'Only use the components listed above. Any other component is not drawn and its content is lost.',
    'At most one openui-lang block per answer. Always state the key point in one plain sentence before the block.',
    "Every value, name and command in a block must come from the retrieved documents, the incident information or the user's message. Never invent values, prices, flags or commands.",
    'Keep blocks compact: at most 4 table columns, 8 table rows and 6 steps. Summarize the rest in prose.',
    'Never repeat the same content in both prose and the block.',
    'Card text cannot hold links or markdown: keep links and the Sources section in the prose outside the block.',
    "Write card text in the user's language, like the rest of the answer.",
    "Only add FollowUps when there is a natural next question the documentation can answer. Clicking one sends it as the user's next question.",
  ],
}

// OpenUI always ends its prompt with a generic rule that suggests charts and forms, which the chat
// can't draw, and has no option to change it: it's swapped for one that only names our components.
const OPENUI_COMPONENT_CHOICE_RULE =
  '- Choose components that best represent the content (tables for comparisons, charts for trends, forms for input, etc.)'

const GORDON_COMPONENT_CHOICE_RULE =
  '- Choose components that best represent the content: use tables for comparisons, Steps for procedures, Callout for incidents or important notes, Code for code or config snippets, and FollowUps for natural next questions.'

const prompt = createGordonCardLibrary().prompt(promptOptions)
if (!prompt.includes(OPENUI_COMPONENT_CHOICE_RULE)) {
  throw new Error(
    'OpenUI no longer writes its generic component-choice rule: update this script',
  )
}
process.stdout.write(
  `${prompt.replace(OPENUI_COMPONENT_CHOICE_RULE, GORDON_COMPONENT_CHOICE_RULE)}\n`,
)
