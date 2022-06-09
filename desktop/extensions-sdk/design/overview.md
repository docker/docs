---
title: UI styling overview
description: Docker extension design
keywords: Docker, extensions, design
---

Our Design System is a constantly evolving set of specifications that aim to ensure visual consistency across Docker products, and meet [level AA accessibility standards](https://www.w3.org/WAI/WCAG2AA-Conformance). We've opened parts of it to extension authors, documenting basic styles (color, typography) and components. See: [Docker Extensions Styleguide](https://www.figma.com/file/U7pLWfEf6IQKUHLhdateBI/Docker-Design-Guidelines?node-id=1%3A28771).

We require extensions to match the wider Docker Desktop UI to a certain degree, and reserve the right to make this stricter in the future.

To get started on your UI, follow the below steps.

## 1. Choose your framework

### Preferred: React+MUI, using our theme

Docker Desktop's UI is written in React and [MUI](https://mui.com/) (using Material UI to specific). This is the only officially supported framework for building extensions, and the one that our `init` command automatically configures for you. Using it brings significant benefits to authors:

- You can use our [Material UI theme](https://www.npmjs.com/package/@docker/docker-mui-theme) to automatically replicate Docker Desktop's look & feel.
- In future, we'll release utilities and components specifically targeting this combination (e.g. custom MUI components, or React hooks for interacting with Docker).

Read our [MUI best practices](mui-best-practices.md) guide to learn future-proof ways to use MUI with Docker Desktop.

### Dispreferred: Some other framework

You may prefer to use another framework, perhaps because you or your team are more familiar with it or because you have existing assets you want to reuse. This is possible, but highly discouraged. It means that:

- You'll need to manually replicate the look and feel of Docker Desktop. This will take a lot of effort, and if you don't match our theme closely enough, users will find your extension jarring and we may ask you to make changes during a review process.
- You'll have a higher maintenance burden. Whenever Docker Desktop's theme changes (which could happen in any release), you'll need to manually change your extension to match it.
- If your extension is open-source, deliberately avoiding common conventions will make it harder for the community to contribute to it.

## 2. Follow the below recommendations

### Follow our MUI best practices (if applicable)

See our [MUI best practices](mui-best-practices.md) article.

### Only use colors from our palette

With minor exceptions, displaying your logo for example, you should only use colors from our palette. These can be found in our [style guide document](https://www.figma.com/file/U7pLWfEf6IQKUHLhdateBI/Docker-Design-Guidelines?node-id=1%3A28771), and will also soon be available in our MUI theme and via CSS variables.

### Use counterpart colors in light/dark mode

Our colors have been chosen so that the counterpart colors in each variant of the palette should have the same essential characteristics. Anywhere you use `red-300` in light mode, you should use `red-300` in dark mode too.
