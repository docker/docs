---
title: Formatting guide
description: Formatting guidelines for technical documentation
keywords: formatting, style guide, contribute
toc_max: 2
---

## Headings and subheadings

As readers skip, skim and scan our pages, they often only notice the first two to three words of any of the paragraphs on the page. They also pay fractionally more attention to headings, bulleted lists and links, so it's important to ensure the first two to three words in those items "front load" information as much as possible.

The best headings should:

- Let the reader know what they will find on the page.
- Describe succinctly and accurately what the content is about.
- Be creative, interesting, short (no more than eight words), to the point and written in plain, active language.
- Avoid puns, teasers and cultural references.
- Skip leading articles (a, the, etc)

## Page title

Try to make sure the title of your page is action orientated:
    - enable SCIM
    - Install Docker Desktop
Make sure the title of your page and the TOC entry matches
If you want to use a ‘:’ in a page title in the table of contents, you must wrap the entire title in “” to avoid breaking the build.

## Images

- Images, including screenshots, can help a reader better understand a concept. However, they should be used sparingly as they tend to go out-of-date quickly.

### Best practice:
- When you take screenshots:
    - Don’t use lorem ipsum text. Try to replicate how the feature would be used in a real-world scenario, and use realistic text.
    - Capture only the relevant UI. Don’t include unnecessary white space or areas of the UI that don’t help illustrate the point.
    - Keep it small. If you don’t need to show the full width of the screen, don’t.
    - Review how the image renders on the page. Make sure the image isn’t blurry or overwhelming.
    - Be consistent. Coordinate screenshots with the other screenshots already on a documentation page for a consistent reading experience.
    - To keep the Git repository light, compress the images (losslessly). Be sure to compress the images **before** adding them to the repository. Compressing images after adding them to the repository actually worsens the impact on the Git repo (however, ut still optimizes the bandwidth during browsing).
    - Don't forget to remove images that are no longer used.

## Links

Be careful not to create too many links, especially within body copy. Excess links can be distracting or send readers away from the current page.

When people follow links, they can often lose their train of thought and fail to return to the original page, despite not having absorbed all the information it contains.

The best links offer readers another way to scan information. 

### Best practice:

- Use plain language that does not mislead or promise too much.
- Be short and descriptive (around five words is best).
- Allow people to predict (with a fair degree of confidence) what they will get if they click mirror the heading text on the destination page in links whenever possible.
- Follow conventions for naming common features (Read more, Next/Previous).
- Front-load user-and-action-oriented terms at the beginning of the link (Download our app).
- Avoid generic calls to action (such as click here, find out more).
- Do not include any end punctuation when you hyperlink text (for example, periods, question marks, or exclamation points).
- Do not make link text italics or bold, unless it would be so as normal body copy.

## Code and code samples

- Format the following as code: docker commands, instructions and filenames (package names).
To apply inline code style, wrap the text in a single backtick (`). For example, this is inline code style

To apply code block style, wrap the text in triple backticks (three `) and add a syntax highlighting hint. For example:

```plaintext
This is codeblock style
```

When using code block style:

Use quadruple backticks (four `) to apply code block style when the code block you are styling has triple backticks in it. For example, when illustrating code block style.
Add a blank line above and below code blocks.

## Callouts

Use callouts to emphasize selected information in a text. 

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
|  | ✅ Example: 
Note 
Deleting a repository deletes all the images it contains and its build settings. This action cannot be undone. |  |
|  |  |  |

## Navigation

When documenting how to navigate through the Docker Desktop or Docker Hub UI:

Always use location, then action.
From the Visibility dropdown list (location), select Public (action).
Be brief and specific. For example:
Do: Select Save.
Do not: Select Save for the changes to take effect.
If a step must include a reason, start the step with it. This helps the user scan more quickly.
Do: To view the changes, in the merge request, select the link.
Do not: Select the link in the merge request to view the changes

## Optional steps

If a step is optional, start the step with the word Optional followed by a period.

For example:

1. Optional. Enter a description for the job.

## Placeholder text 

You might want to provide a command or configuration that uses specific values.

In these cases, use < and > to call out where a reader must replace text with their own value.

## Tables

Tables should be used to describe complex information in a straightforward manner. Note that in many cases, an unordered list is sufficient to describe a list of items with a single, simple description per item. But, if you have data that’s best described by a matrix, tables are the best choice.

### Best practice:

- Use sentence case for table headings.
- To keep tables accessible and scannable, tables should not have any empty cells. If there is no otherwise meaningful value for a cell, consider entering N/A for ‘not applicable’ or None.
