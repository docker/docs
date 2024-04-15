---
description: Learn about best practices for videos in docs, and how about the component to add videos.
title: Videos
toc_max: 3
---

## Video guidelines

Videos are used rarely in Docker's documentation. When used, video should complement written text and not be the sole format of documentation. Videos can take longer to produce and be more difficult to maintain than written text or even screenshots, so consider the following before adding video:

- Can you demonstrate a clear customer need for using video ?
- If the video contains user interfaces that may change regularly, do you have a maintenance plan to keep the video fresh?
- Does the [voice and tone](../style/voice-tone.md) of the video align with the rest of the documentation?
- Have you considered other options, such as screenshots or clarifying existing documentation?
- Is the quality of the video similar to the rest of Docker's documentation?
- Can the video be linked or embedded from the site?

If all of the above criteria are met, you can reference the following best practices before creating a video to add to Docker documentation.

### Best practices

- Determine the audience for your video. Will the video be a broad overview for beginners, or will it be a deep dive into a technical process designed for an advanced user?
- Videos should be less than 5 minutes long. Keep in mind how long the video needs to be to properly explain the topic, and if the video needs to be longer than 5 minutes, consider text, diagrams, or screenshots instead. These are easier for a user to scan for relevant information.
- Videos should adhere to the same standards for accessibility as the rest of the documentation.
- Ensure the quality of your video by writing a script (if there's narration), making sure multiple browsers and URLs aren't visible, blurring or cropping out any sensitive information, and using smooth transitions between different browsers or screens.

Videos are not hosted in the Docker documentation repository. To add a video, you can use a [link](./links.md) to hosted content, or embed using an [iframe](#iframe).


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
