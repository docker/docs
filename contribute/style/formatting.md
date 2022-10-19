---
title: Formatting guide
description: Formatting guidelines for technical documentation
keywords: formatting, style guide, contribute
toc_max: 2
---

## Headings and subheadings

Readers pay fractionally more attention to headings, bulleted lists and links, so it's important to ensure the first two to three words in those items "front load" information as much as possible.

### Best practice:

- Headings and subheadings should let the reader know what they will find on the page.
- They should describe succinctly and accurately what the content is about.
- Headings should be short (no more than eight words), to the point and written in plain, active language.
- You should avoid puns, teasers and cultural references.
- Skip leading articles (a, the, etc.)

## Page title

Page titles should be action orientated. For example:
    - _Enable SCIM_
    - _Install Docker Desktop_

### Best practice:

- Make sure the title of your page and the TOC entry matches
- If you want to use a ‘:’ in a page title in the table of contents (_toc.yaml), you must wrap the entire title in “” to avoid breaking the build.

## Images

Images, including screenshots, can help a reader better understand a concept. However, they should be used sparingly as they tend to go out-of-date quickly.

### Best practice:
- When you take screenshots:
    - Don’t use lorem ipsum text. Try to replicate how the feature would be used in a real-world scenario, and use realistic text.
    - Capture only the relevant UI. Don’t include unnecessary white space or areas of the UI that don’t help illustrate the point.
    - Keep it small. If you don’t need to show the full width of the screen, don’t.
    - Review how the image renders on the page. Make sure the image isn’t blurry or overwhelming.
    - Be consistent. Coordinate screenshots with the other screenshots already on a documentation page for a consistent reading experience.
    - To keep the Git repository light, compress the images (losslessly). Be sure to compress the images before adding them to the repository. Compressing images after adding them to the repository actually worsens the impact on the Git repo (however, it still optimizes the bandwidth during browsing).
    - Don't forget to remove images from the repo that are no longer used.

For information on how to add images to your content, see [Useful component and formatting examples](../components/images.md).

## Links

Be careful not to create too many links, especially within body copy. Excess links can be distracting or send readers away from the current page.

When people follow links, they can often lose their train of thought and fail to return to the original page, despite not having absorbed all the information it contains.

The best links offer readers another way to scan information.

### Best practice:

- Use plain language that does not mislead or promise too much.
- Be short and descriptive (around five words is best).
- Allow people to predict (with a fair degree of confidence) what they will get if they click. Mirror the heading text on the destination page in links whenever possible.
- Front-load user-and-action-oriented terms at the beginning of the link (Download our app).
- Avoid generic calls to action (such as click here, find out more).
- Do not include any end punctuation when you hyperlink text (for example, periods, question marks, or exclamation points).
- Do not make link text italics or bold, unless it would be so as normal body copy.
- Make sure your link opens in a new tab so it doesn't interrupt the user-flow.

For information on how to add links to your content, see [Useful component and formatting examples](../components/links.md).

## Code and code samples

Format the following as code: docker commands, instructions and filenames (package names).

To apply inline code style, wrap the text in a single backtick (`).

For information on how to add codeblocks to your content, see [Useful component and formatting examples](../components/code-blocks.md).

## Callouts

Use callouts to emphasize selected information in a page.

### Best practice:

- Format the word Warning, Important, or Note in bold. Do not bold the content within the callout.
- It's good practice to avoid placing a lot of text callouts on one page. They can create a cluttered appearance if used to excess, and you'll diminish their impact.

| Text callout  | Use case scenario | Color or callout box |
| --- | --- | --- |
| Warning | Use a Warning tag to signal to the reader where actions may cause damage to hardware or software loss of data.  | Red |
|  | ✅ Example: Warning: When you use the docker login command, your credentials are stored in your home directory in .docker/config.json. The password is base64-encoded in this file. |  |
| Important | Use an Important tag to signal to the reader where actions may cause issues of a lower magnitude. | Yellow |
|  | ✅ Example: Update to the Docker Desktop terms |  |
| Note | Use the Note tag for information that may not apply to all readers, or if the information is tangential to a topic. | Blue |
|  | ✅ Example: Deleting a repository deletes all the images it contains and its build settings. This action cannot be undone.|

For information on how to add call outs to your content, see [Useful component and formatting examples](../components/call-outs.md).

## Navigation

When documenting how to navigate through the Docker Desktop or Docker Hub UI:

- Always use location, then action. For example:
    _From the **Visibility** dropdown list (location), select Public (action)._
- Be brief and specific. For example:
    Do: _Select **Save**._
    Do not: _Select **Save** for the changes to take effect._
- If a step must include a reason, start the step with it. This helps the user scan more quickly.
    Do: _To view the changes, in the merge request, select the link._
    Do not: _Select the link in the merge request to view the changes_

## Optional steps

If a step is optional, start the step with the word Optional followed by a period.

For example:

_1. Optional. Enter a description for the job._

## Placeholder text

You might want to provide a command or configuration that uses specific values.

In these cases, use < and > to call out where a reader must replace text with their own value. For example:

`docker extension install <name-of-your-extension>`

## Tables

Tables should be used to describe complex information in a straightforward manner.

Note that in many cases, an unordered list is sufficient to describe a list of items with a single, simple description per item. But, if you have data that’s best described by a matrix, tables are the best choice.

### Best practice:

- Use sentence case for table headings.
- To keep tables accessible and scannable, tables should not have any empty cells. If there is no otherwise meaningful value for a cell, consider entering N/A for ‘not applicable’ or None.

For information on how to add tables to your content, see [Useful component and formatting examples](../components/tables.md).
