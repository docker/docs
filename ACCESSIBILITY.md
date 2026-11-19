# Accessibility review

This review covers the Docs rebrand's shared templates and interactive
components. It does not establish whole-site WCAG conformance. Passing browser
tests does not replace testing with a screen reader or reviewing the meaning
of documentation and media.

## Automated and browser checks

Run against a built site with Pagefind, following the preparation steps in
[hack/redesign/README.md](hack/redesign/README.md#verification):

```console
$ DOCS_TEST_URL=http://localhost:1316 npm run test:a11y
$ DOCS_TEST_URL=http://localhost:1316 npm run test:redesign
```

Set `CHROMIUM_PATH` when using an existing Chromium installation. Axe is a
development dependency; the website does not load it.

The accessibility suite checks the homepage, category and product landings,
installation instructions, Dockerfile and CLI reference, interactive diagrams,
image zoom, search, chat responses, and chat errors. API filtering and request
copying also have checks for status announcements. It runs in both themes at
1440px and 320px, with additional checks at a viewport equivalent to 200%
browser zoom. It blocks external requests and mocks chat responses and errors.

The suite checks WCAG A/AA rules, including WCAG 2.2, and axe best practices.
It also exercises keyboard behavior that those rules cannot establish alone.
Axe's incomplete results, especially contrast over gradients, require review.
The suite writes those findings, including page and element details, to
`tmp/redesign/accessibility-incomplete.json`, including when a check fails.

| Checklist area | Coverage and findings |
| --- | --- |
| Landmarks, language, skip link | Shared shell declares English, provides header, main, footer, and named navigation. The skip link moves focus to the main content. |
| Keyboard navigation | Tests exercise tab arrows, Home and End, sidebar disclosures, image zoom, search, and chat. |
| Dialogs | Search and image zoom use native dialogs. Chat uses Alpine's focus trap. Tests check dismissal, focus restoration, and chat focus containment in both directions. |
| Tabs | Tablist, tab, and tabpanel relationships are present. Selection uses a single tab stop. Unavailable synchronized selections cannot hide every panel. |
| Forms and errors | Inputs have accessible names. Search status and completed chat answers use live regions; feedback errors use alerts. Chat announcements wait for a completed answer instead of announcing every streamed fragment. The shared shell has no personal-information fields requiring autocomplete. |
| Contrast and states | Sampled text over the rendered category washes passes AA in both themes. Hover, active, focus, and disabled controls received additional inspection. Disabled controls retain readable text instead of fading the entire control. Syntax colors have contrast assertions in the existing browser suite. |
| Focus visibility | Shared focus rings use Trident colors, with a darker light-mode ring over washes. Scroll padding reserves space for the floating Ask input. Search results reserve space for focus outlines. |
| Targets | Axe checks WCAG 2.2 target sizing and spacing. Sidebar toggles and tooltip buttons have explicit target sizes. Inline text links use the WCAG spacing and inline exceptions. |
| Zoom and reflow | Page sweeps reject page-level horizontal overflow. Tables, code, and wide diagrams remain independently scrollable by keyboard. Search and populated chat are also checked on mobile. |
| Cards | Whole-card links have one link target. Grid cards use level-two headings by default; authors can select a level that fits their outline. |
| Images | Markdown images retain author-provided alt text. Image zoom supports keyboard activation, a named close control, Escape, and focus restoration. Two figure shortcodes missing alt text were corrected. |
| Tables | Markdown headers declare column scope. Generated CLI tables have names and row or column scopes. Chat tables retain table semantics and support keyboard scrolling. Authored table names and complex header relationships still need content review. |
| Motion | Reduced-motion styles shorten transitions and stop repeating CSS animations. Chat and code expansion respect the preference when scrolling. Interactive diagrams skip token travel animations. Markdown GIFs start inside closed native disclosures so readers choose when to show or hide them. |

In the screenshot-based contrast sample, the lowest ordinary text contrast over
light-mode washes was 5.21:1, above the 4.5:1 AA threshold. This sample covered
the homepage, AI and application development categories, Sandboxes, an Ubuntu
installation page, and Dockerfile reference. It does not measure every possible
content, state, or background combination.

## Outstanding content and assistive-technology review

The rendered-content inventory on 2026-11-18 found the following work outside
the shared component fixes:

- 896 tables lack a caption or explicit accessible name. Some have sufficient
  surrounding context for WCAG, but they do not satisfy the checklist's stricter
  requirement for an individually named table. Add names that describe each
  table's purpose; do not generate generic labels to silence checks.
- 45 heading-level jumps remain in authored or vendored content. Review the
  outline before changing levels. Examples include Compose FAQs, Compose
  production guidance, Offload Sandboxes, and vendored API descriptions. These
  are inventory flags, rather than 45 independently confirmed WCAG failures.
- Check caption and transcript availability for embedded videos, including the
  Docker Scout GitLab recording. The Builds export video has no audio track
  and is accompanied by written instructions; check that those steps cover
  its relevant visual information.
- Review animated assets for flashing before publishing them. The Docker Agent
  TUI GIF loops for about 15 seconds; it and the Engine 27 device-login GIF are
  hidden until a reader expands their disclosure. CSS reduced-motion rules
  alone cannot stop GIF playback.
- Have content owners review whether alt text explains each image, link text
  makes sense in isolation, terminology suits the audience, and acronyms are
  introduced. Attribute presence cannot establish those properties.
- Test the primary flows with VoiceOver/Safari and NVDA/Firefox or Chrome.
  Keyboard and accessibility-tree checks were performed here; spoken output
  with those screen readers was not tested.
- Review production cookie-consent UI and third-party media separately. The
  repeatable local tests block external requests.

For the screen-reader pass, start at the skip link, navigate the heading and
landmark lists, switch installation tabs, expand a sidebar section, zoom an
image, search for a term, and open chat. Confirm that results, completed answers,
and errors are announced once, that inactive content stays out of navigation,
and that closing each dialog returns to a useful focus target. Repeat with the
floating Ask input and with a narrow viewport.

Vendored content must be corrected in its source repository, as described in
[AGENTS.md](AGENTS.md#vendored-content-do-not-edit).
