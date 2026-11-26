# Accessibility review

This review covers the Docs rebrand's shared templates and interactive
components. It does not establish whole-site WCAG conformance. Passing browser
tests does not replace testing with a screen reader or reviewing the meaning
of documentation and media.

## Automated and browser checks

Run `docker buildx bake test-site`, or follow the local setup in
[CONTRIBUTING.md](CONTRIBUTING.md#testing).

The shared browser suite scans the homepage, Build landing page, Desktop
installation page, and Buildx command reference with axe. It checks WCAG A/AA
rules through WCAG 2.2 and axe best practices in desktop light mode and mobile
dark mode. Keyboard checks cover sidebar disclosure, installation tabs, search
and chat focus restoration, and code copying. Tests block external requests.

Axe's incomplete findings, especially contrast over gradients, need manual
review. They are attached to each test's results under `tmp/test-results`.
Automated scans cover representative pages and do not establish whole-site
conformance or replace assistive-technology testing.

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
