# Docker Documentation Style Guide

This guide consolidates voice, grammar, formatting, and terminology
standards for Docker documentation. Follow these guidelines to create
clear, consistent, and helpful content. For instructions on how to use
components, shortcodes, and other features, see [COMPONENTS.md](COMPONENTS.md).

## Voice and tone

Write like a knowledgeable colleague explaining something useful. We're
developers writing for developers.

### Core principles: The 4Cs

1. **Correct** - Information is accurate
2. **Concise** - Succinct without unnecessary words
3. **Complete** - Includes enough detail to complete the task
4. **Clear** - Easy to understand

### Writing approach

- **Be honest.** Give all the facts. Don't use misdirection or
  ambiguous statements.
- **Be concise.** Don't bulk up communication with fluffy words or
  complex metaphors. Get to the point.
- **Be relaxed.** Casual but not lazy, smart but not arrogant, focused
  but not cold. Be welcoming and warm.
- **Be inclusive.** Every person is part of our community, regardless
  of experience level.

### Tone guidelines

- Use a natural, friendly, and respectful tone
- Use contractions to sound conversational (it's, you're, don't)
- Avoid overdoing politeness - skip "please" in most technical documentation
- Be clear over comical
- Use positive language - emphasize what users can do, not what they can't

**Positive language example:**

Instead of: "Features such as Single Sign-on (SSO), Image Access
Management, Registry Access Management are not available in Docker Team
subscription."

Use: "Features such as Single Sign-on (SSO), Image Access Management,
Registry Access Management are available in Docker Business subscription."

### Avoiding marketing language

Documentation should be factual and direct, not promotional.

**Avoid hedge words** that overstate ease or capability:

- ❌ simply, just, easily, seamlessly (implies ease when it may not be
  easy)
- ❌ robust, powerful (marketing language)
- ✅ Instead: describe what it actually does

**Avoid excessive enthusiasm:**

- ❌ "powerful feature", "game-changing", "revolutionary", "amazing"
- ✅ Instead: describe the feature and its benefits factually

**Avoid meta-commentary:**

These phrases add no value - state the fact directly:

- ❌ "It's worth noting that..."
- ❌ "It's important to understand that..."
- ❌ "Keep in mind that..."
- ✅ Instead: state the information directly

### Voice and perspective

- **Use "you" not "we"**: Focus on what the user can do, not what "we"
  created
  - ❌ "We provide a feature that helps you deploy"
  - ✅ "Deploy your applications with..."
- **Avoid "please"**: Don't use in normal explanations or "please note"
  phrases
- **Write timelessly**: Avoid "currently" or "as of this writing" - the
  documentation describes the product as it is today

### Scope preservation

When updating existing documentation, resist the urge to expand
unnecessarily. Users value brevity.

**Understand the current document:**
Read it fully to grasp its scope, length, and character. Is it a minimal
how-to or a comprehensive reference?

**Match the existing character:**
If the document is brief and direct (90 lines), keep it that way. Don't
transform a focused guide into an exhaustive tutorial.

**Add only what's genuinely missing:**
Fill gaps, don't elaborate. If the document already covers a topic
adequately, don't expand it.

**Value brevity:**
Say what needs to be said, then stop. Not every topic needs
prerequisites, troubleshooting, best practices, and examples sections.

**Respect the original intent:**
The document exists in its current form for a reason. Improve it, don't
remake it.

Good additions fill genuine gaps. Bad additions change the document's
character. When in doubt, add less rather than more.

## Grammar and style

Write in US English with US grammar.

### Acronyms and initialisms

- Spell out lesser-known acronyms on first use, then follow with the
  acronym in parentheses: "single sign-on (SSO)"
- Don't spell out common acronyms like URL, HTML, API
- Use all caps for file type acronyms: JPEG, PNG
- Don't use apostrophes for plurals: URLs not URL's
- Avoid introducing acronyms in headings - introduce them in the body
  text that follows

### Capitalization

Use sentence case for almost everything: headings, titles, links, buttons,
navigation labels.

- Capitalize Docker solutions: Docker Desktop, Docker Hub, Docker Extensions
- Capitalize job titles only when they immediately precede a name:
  "Chief Executive Officer Scott Johnston" but "Scott Johnston, chief
  executive officer"
- Follow company capitalization preferences: DISH, bluecrux
- Match UI capitalization when referring to specific interface elements

### Contractions

Use contractions to maintain a conversational tone, but don't overdo it.

- Stay consistent - don't switch between "you are" and "you're" in the
  same content
- Avoid negative contractions when possible (can't, don't, won't) -
  rewrite to be positive
- Never contract a noun with a verb (your container is ready, not your
  container's ready)

### Dangling modifiers

Avoid unclear subjects:

- ❌ After enabling auto-log-out, your users are logged out
- ✅ When you enable auto-log-out, your users are logged out

### Dates and numbers

- Use US date format: June 26, 2021 or Jun. 26, 2021
- Spell out numbers one through nine, except in units: 4 MB
- Use numerals for 10 and above: 10, 625, 1000
- Use decimals instead of fractions: 0.5 not ½
- For numbers with five or more digits, use spaces instead of commas:
  14 586 not 14,586
- For version numbers: version 4.8.2, v1.0, Docker Hub API v2

### Punctuation

- **Commas:** Use the Oxford comma (serial comma) in lists
- **Colons:** Use to introduce lists or provide explanations
- **Semicolons:** Don't use - write two sentences instead
- **Em dashes:** Use sparingly with spaces on either side: "text - like
  this - text"
- **Hyphens:** Use for compound adjectives before nouns:
  "up-to-date documentation"
- **Exclamation marks:** Avoid
- **Parentheses:** Avoid in technical documentation - they reduce
  readability

### Conciseness and redundant phrases

Remove unnecessary words to make documentation clearer and more direct.

**Eliminate redundant phrases:**

- ❌ "in order to" → ✅ "to"
- ❌ "serves the purpose of" → ✅ state what it does directly
- ❌ "allows you to" → ✅ "lets you" or state what it does
- ❌ "enables you to" → ✅ "lets you" or state what it does

### Bold and italics

Use bold **only** for UI elements (buttons, menus, field labels). Never
use bold for emphasis, product names, or feature names.

- ✅ Select **Save** (UI button - use bold)
- ✅ Docker Hub provides storage (product name - no bold)
- ✅ This is important to understand (emphasis - no bold)
- ❌ **Docker Desktop** is a tool (product name - don't bold)
- ❌ This is **very important** (emphasis - don't bold)
- ❌ The **build** command (command name - don't bold)

**Italics:**
Use italics sparingly, as this formatting can be difficult to read in
digital experiences. Notable exceptions are titles of articles, blog
posts, or specification documents.

## Formatting

### Content types and detail balance

Different content types require different writing approaches. Match your
style and detail level to the content type.

**Getting Started / Tutorials:**

- Step-by-step instructions
- Assume beginner knowledge
- Explain _why_ at each step
- Include more context and explanation

**How-to Guides:**

- Task-focused and goal-oriented
- Assume intermediate knowledge
- Focus on _how_ efficiently
- Less explanation, more direct steps

**Reference Documentation:**

- Comprehensive and exhaustive
- Assume advanced knowledge
- Focus on _what_ precisely
- Complete parameter lists and options

**Concept Explanations:**

- Educational and foundational
- Any skill level
- Focus on _understanding_ over doing
- Theory before practice

**Match detail to context:**

- **First mention** of a concept: explain it
- **Subsequent mentions**: link to the explanation or use the term
  directly if recently explained
- **Common knowledge** (to your audience): state it, don't explain it
- **Edge cases**: mention them but don't let them dominate the main flow

### Headings

- Keep headings short (no more than eight words)
- Front-load the most important words
- Use descriptive, action-oriented language
- Skip leading articles (a, the)
- Avoid puns, teasers, and cultural references
- Page titles should be action-oriented: "Install Docker Desktop",
  "Enable SCIM"

### Page structure and flow

Every page should answer two questions in the first paragraph:

1. **What will I learn?** - State what the page covers
2. **Why does this matter?** - Explain the benefit or use case

**Good opening:**

```markdown
Docker Compose Watch automatically updates your running containers when
you change code. This eliminates the manual rebuild-restart cycle during
development.
```

**Weak opening:**

```markdown
Docker Compose Watch is a powerful feature that enables developers to
streamline their development workflow by providing automatic
synchronization capabilities.
```

**Transitions and flow:**
Connect ideas naturally. Each section should flow logically from the
previous one. Save detailed discussions for after showing basic usage -
don't front-load complexity.

**Good flow:**

1. Prerequisites
2. Basic usage (show the simple case first)
3. Advanced options (add complexity after basics are clear)

**Jarring flow:**

1. Prerequisites
2. Overview of all capabilities (too much before seeing it work)
3. Basic usage

**Avoid structural problems:**

- Don't start multiple sentences or sections with the same structure
  (varies pacing and improves readability)
- Don't over-explain obvious things (trust the reader's competence)
- Don't use "**Bold:** format" for subsection labels (use plain text
  with colon)

### Lists

- Limit bulleted lists to five items when possible
- Don't add commas or semicolons to list item ends
- Capitalize list items for ease of scanning
- Make all list items parallel in structure
- Start each item with a capital letter unless it's a parameter or
  command
- For nested sequential lists, use lowercase letters: 1. Top level,
  a. Child step

**Avoid marketing-style list formatting:**

Don't use "**Term** - Description" format, which reads like marketing
copy:

❌ **Bad example:**

- **Build** - Creates images from Dockerfiles
- **Run** - Starts containers from images
- **Push** - Uploads images to registries

✅ **Good alternatives:**

Use simple descriptive bullets:

- Build images from Dockerfiles
- Start containers from images
- Upload images to registries

Or use proper prose when appropriate:
"Docker lets you build images from Dockerfiles, start containers from
images, and upload images to registries."

### Code and commands

Format Docker commands, instructions, filenames, and package names as
inline code using backticks: `docker build`

**Code example pattern:**

Follow this three-step pattern for code examples:

1. **State** what you're doing (brief)
2. **Show** the command or code
3. **Explain** the result or key parts (if not obvious)

**Example:**

````markdown
Build your image:

```console
$ docker build -t my-app .
```

This creates an image tagged `my-app` using the Dockerfile in the
current directory.
````

**When to show command output:**

- Show output when it helps understanding
- Show output when users need to verify results
- Don't show output for commands with obvious results
- Don't show output when it's not relevant to the point

For code block syntax, language hints, variables, and advanced features
(line highlighting, collapsible blocks), see
[COMPONENTS.md](COMPONENTS.md#code-blocks).

### Links

- Use plain, descriptive language (around five words)
- Front-load important terms at the beginning
- Avoid generic calls to action like "click here" or "find out more"
- Don't include end punctuation in link text
- Use relative links with source filenames for internal links
- Don't make link text bold or italic unless it would be in normal
  body copy

### Images

**Best practices:**

- Capture only relevant UI - don't include unnecessary white space
- Use realistic text, not lorem ipsum
- Keep images small and focused
- Compress images before adding to repository
- Remove unused images from repository

For image syntax and parameters (sizing, borders), see
[COMPONENTS.md](COMPONENTS.md#images).

### Callouts

Use callouts sparingly to emphasize important information. Use them only
when information truly deserves special emphasis - for warnings, critical
notes, or important context that affects how users approach a task.

Callout types: Note (informational), Tip (helpful suggestion), Important
(moderate issue), Warning (potential damage/loss), Caution (serious risk).

For syntax and detailed usage guidelines, see
[COMPONENTS.md](COMPONENTS.md#callouts).

### Tables

- Use sentence case for table headings
- Don't leave cells empty - use "N/A" or "None" if needed
- Align decimals on the decimal point

### Navigation instructions

- Use location, then action: "From the **Visibility** drop-down list, select Public"
- Be brief and specific: "Select **Save**"
- Start with the reason if a step needs one: "To view changes, select the link"

### Optional steps

Start optional steps with "Optional." followed by the instruction:

_1. Optional. Enter a description for the job._

### File types and units

- Use formal file type names: "a PNG file" not "a .png file"
- Use abbreviated units with spaces: "10 GB" not "10GB" or
  "10 gigabytes"

## Word list

### Recommended terms

| Use                                   | Don't use              |
| ------------------------------------- | ---------------------- |
| lets                                  | allows                 |
| turn on, toggle on                    | enable                 |
| turn off, toggle off                  | disable                |
| select                                | click                  |
| following                             | below                  |
| previous                              | above                  |
| checkbox                              | check box              |
| username                              | account name           |
| sign in                               | sign on, log in, login |
| sign up, create account               | register               |
| for example, such as                  | e.g.                   |
| run                                   | execute                |
| want                                  | wish                   |
| Kubernetes                            | K8s                    |
| repository                            | repo                   |
| administrator (first use), admin (UI) | -                      |
| and                                   | & (ampersand)          |
| move through, navigate to             | scroll                 |
| (be more precise)                     | respectively           |
| versus                                | vs, vs.                |
| use                                   | utilize                |
| help                                  | facilitate             |

### Version numbers

- Use "earlier" not "lower": Docker Desktop 4.1 and earlier
- Use "later" not "higher" or "above": Docker Desktop 4.1 and later

### Quick transformations

Common phrases to transform for clearer, more direct writing:

| Instead of...            | Write...                     |
| ------------------------ | ---------------------------- |
| In order to build        | To build                     |
| This allows you to build | This lets you build / Build  |
| Simply run the command   | Run the command              |
| We provide a feature     | Docker provides / You can    |
| Utilize the API          | Use the API                  |
| This facilitates testing | This helps test / This tests |
| Click the button         | Select the button            |
| It's worth noting that X | X (state it directly)        |

### UI elements

- **tab vs view** - Use "view" for major UI sections, "tab" for
  sub-sections
- **toggle** - You "turn on" or "turn off" a toggle

## Docker terminology

### Products and features

- **Docker Compose** - The application or all functionality associated
  with it
- **`docker compose`** - Use code formatting for commands
- **Docker Compose CLI** - The family of Compose commands from
  Docker CLI
- **Compose plugin** - The add-on for Docker CLI that can be
  enabled/disabled
- **compose.yaml** - Current designation for the Compose file (use
  code formatting)

### Technical terms

- **Digest** - Long string automatically created when you push an
  image
- **Member** - A user of Docker Hub who is part of an organization
- **Namespace** - Organization or user name; every image needs a
  namespace
- **Node** - Physical or virtual machine running Docker Engine in
  swarm mode
  - Manager nodes: Perform swarm management and orchestration
  - Worker nodes: Execute tasks
- **Registry** - Online storage for Docker images
- **Repository** - Lets users share container images with team,
  customers, or community

### Platform terminology

- **Multi-platform** - Broad: Mac/Linux/Windows; Narrow: Linux/amd64
  and Linux/arm64
- **Multi-architecture (multi-arch)** - CPU architecture or
  hardware-architecture-based; don't use as synonym for multi-platform

## Quick reference

### Voice checklist

- ✅ Direct and practical
- ✅ Conversational with active voice
- ✅ Specific and actionable
- ❌ Corporate-speak
- ❌ Condescending
- ❌ Overly formal

### Structure checklist

- ✅ Answer "What will I learn?" in the first paragraph
- ✅ Answer "Why does this matter?" in the first paragraph
- ✅ Front-load important information in headings and lists
- ✅ Connect ideas naturally between sections
- ✅ Use examples to clarify concepts
- ✅ Match content type (tutorial vs how-to vs reference vs concept)
- ✅ Follow code example pattern: state → show → explain
- ✅ Preserve document scope (don't unnecessarily expand)

### Common mistakes

- ❌ Using "click" instead of "select"
- ❌ Using "below" instead of "following"
- ❌ Starting sentences with numbers (recast the sentence)
- ❌ Using apostrophes in plural acronyms
- ❌ Leaving table cells empty
- ❌ Using commas in large numbers (use spaces)
- ❌ Referring to file types by extension (.png file)
- ❌ Using bold for product names or emphasis (only for UI elements)
- ❌ Using "**Term** - Description" format in lists
- ❌ Using hedge words (simply, easily, just, seamlessly)
- ❌ Using meta-commentary (it's worth noting that...)
- ❌ Using "we" instead of "you" or "Docker"
- ❌ Showing command output when it's not needed
- ❌ Not explaining what users will learn in first paragraph
