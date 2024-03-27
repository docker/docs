---
description: components and formatting examples used in Docker's docs
title: Videos
toc_max: 3
---

## iframe

To embed a video on a docs page, use an `<iframe>` element:

```html
<iframe
  class="border-0 w-full aspect-video mb-8"
  allow="fullscreen"
  title=""
  src=""
  ></iframe>
```

## asciinema

`asciinema` is a command line tool for recording terminal sessions. The
recordings can be embedded on the documentation site. These are similar to
`console` code blocks, but since they're playable and scrubbable videos, they
add another level of usefulness over static codeblocks in some cases. Text in
an `asciinema` "video" can also be copied, which makes them more useful.

Consider using an `asciinema` recording if:

- The input/output of the terminal commands are too long for a static example
  (you could also consider truncating the output)
- The steps you want to show are easily demonstrated in a few commands
- Where the it's useful to see both inputs and outputs of commands

To create an `asciinema` recording and add it to docs:

1. [Install](https://docs.asciinema.org/getting-started/) the `asciinema` CLI
2. Run `asciinema auth` to configure your client and create an account
3. Start a new recording with `asciinema rec`
4. Run the commands for your demo and stop the recording with `<C-d>` or `exit`
5. Upload the recording to <asciinema.org>
6. Embed the player with a `<script>` tag using the **Share** button on <asciinema.org>
