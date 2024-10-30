---
title: Guidelines for writing Docker usage guides
linkTitle: Write a Docker guide
description: Learn how to write guides for learning about Docker, with Docker.
---

<!-- vale Docker.We = NO -->

This guide provides instructions and best practices for writing tutorial-style
usage guides that help users achieve specific goals using Docker. These guides
will be featured in the [guides section](/guides/_index.md) of the website,
alongside our product manuals and reference materials.

Our documentation is written in Markdown, with YAML front matter for metadata.
We use [Hugo](https://gohugo.io) to build the website. For more information
about how to write content for Docker docs, refer to our
[CONTRIBUTING.md](https://github.com/docker/docs/tree/main/CONTRIBUTING.md).

## Purpose of the guides

Our usage guides aim to:

- Educate users on how to leverage Docker in various contexts.
- Provide practical, hands-on experience through step-by-step tutorials.
- Help users achieve specific goals relevant to their interests or projects.

## Audience

- Experience levels: From beginners to advanced users.
- Roles: Developers, system administrators, DevOps engineers, and more.
- Technologies: Various programming languages and frameworks.

## Metadata for guides

Each guide must include a metadata section at the beginning of the document.
This metadata helps users discover and filter guides based on their interests
and needs.

### Example metadata format

```yaml
---
title: Deploy a machine learning model with Docker
linkTitle: Docker for ML deployment
description: Learn how to containerize and deploy a machine learning model using Docker.
summary: |
  This guide walks you through the steps to containerize a machine learning
  model and deploy it using Docker, enabling scalable and portable AI
  solutions.
tags: [machine-learning, deployment]
languages: [python]
params:
  time: 30 minutes
---
```

### Metadata fields

- `title` (required): The main title of the guide in sentence case.
- `linkTitle` (optional): A shorter title used in navigation menus.
- `description` (required): A concise description for SEO purposes.
- `summary` (required): A brief overview of the guide's content.
- `languages` \* (optional): List of programming languages used.
- `tags` \* (optional): Domains or subject areas covered.
- `params`
  - `time` (optional): Estimated reading or completion time.

\* Do apply at least one of the `languages` or `tags` taxonomies. The values
are used to associate the page with the filter options on the guides landing
page.

## Document structure

Our guides emphasize learning by doing. Depending on the type of guide, the
structure may vary to best suit the content and provide a smooth learning
experience.

All guides live directly under the `content/guides/` directory in the Docker
docs repository. Guides can either be a single page or multiple pages. In the
case of multi-page guides, every page is a step in a sequential workflow.

If you're creating a single-page guide, create a single markdown file in the
guides directory:

```bash
# Create the file
touch content/guides/my-docker-guide.md
# or if you have Hugo installed:
hugo new content/guides/my-docker-guide.md
```

To create a multi-page guide, create a directory where each page is a markdown
file, with an `_index.md` file representing the introduction to the guide.

```bash
# Create the index page for the guide
mkdir content/guides/my-docker-guide.md
touch content/guides/my-docker-guide/_index.md
# or if you have Hugo installed:
hugo new content/guides/my-docker-guide/_index.md
```

Then create the pages for the guide under `content/guides/<dir>/<page>.md` as
needed. The [metadata](#metadata-for-guides) lives on the `_index.md` page (you
don't need to assign metadata to individual pages).

### Guides for specific frameworks or languages

For guides that demonstrate how to use Docker with a particular framework or
programming language, consider the following outline:

1. **Introduction**
   - Briefly introduce the framework or language in the context of Docker.
   - Explain what the user will achieve by the end of the guide.
   - List required software, tools, and knowledge.
2. **Development setup**
   - Guide the user through setting up a development environment.
   - Include instructions on writing or obtaining sample code.
   - Show how to run containers for local development.
3. **Building the application**
   - Explain how to create a Dockerfile tailored to the application.
   - Provide step-by-step instructions for building the Docker image.
   - If applicable, show how to test the application using Docker.
4. **Deploying with Docker**
   - Show how to run the application in a Docker container.
   - Discuss configuration options and best practices.
5. **Best practices and conclusions**
   - Offer tips for optimizing Docker usage with the framework or language.
   - Summarize key takeaways, suggest next steps, and further reading.

### Use-case guides

For guides focused on accomplishing a specific goal or use case with Docker
(e.g., deploying a machine learning model), use a flexible outline that ensures
a logical flow.

The following outline is an example. The structure should be adjusted to best
fit the content and ensure a clear, logical progression. Depending on the
subject matter of your use-case guide, the exact structure will vary.

1. **Introduction**
   - Describe the problem or goal.
   - Explain the benefits of using Docker in this context.
2. **Prerequisites**
   - List required tools, technologies, and prior knowledge.
3. **Setup**
   - Provide instructions for setting up the environment.
   - Include any necessary configuration steps.
4. **Implementation**
   - Walk through the process step by step.
   - Use code snippets and explanations to illustrate key points.
5. **Running or deploying the application**
   - Guide the user on how to execute or deploy the solution.
   - Discuss any verification steps to ensure success.
6. **Conclusion**
   - Recap what was achieved.
   - Suggest further reading or advanced topics.

## Example code

If you create an example repository with source code to accompany your guide,
we strongly encourage you to publish that repository under the
[dockersamples](https://github.com/dockersamples) organization on GitHub.
Publishing your source code under this organization, rather than your personal
account, ensures that the source code remains accessible to the maintainers of
the documentation site after publishing. In the event of a bug or an issue in
the guide, the documentation team can more easily update the guide and its
corresponding example repository.

Hosting the examples in the official samples namespace also helps secure trust
with users who are reading the guide.

### Publish a repository under `dockersamples`

To publish your repository under the `dockersamples` organization, use the
[dockersamples template](https://github.com/dockersamples/sample-app-template)
to initiate a sample repository under your personal namespace. When you've
finished drafting your content and opened your pull request for the
documentation, we can transfer the repository to the dockersamples
organization.

## Writing style

Use [sentence case](/contribute/style/formatting.md#capitalization) for all
headings and titles. This means only the first word and proper nouns are
capitalized.

### Voice and tone

- **Clarity and conciseness**: Use simple language and short sentences to convey information effectively.
- **Active voice**: Write in the active voice to engage the reader (e.g., "Run the command" instead of "The command should be run").
- **Consistency**: Maintain a consistent tone and terminology throughout the guide.

For detailed guidelines, refer to our [voice and tone documentation](/contribute/style/voice-tone.md).

### Formatting conventions

- **Headings**: Use H2 for main sections and H3 for subsections, following sentence case.
- **Code examples**: Provide complete, working code snippets with syntax highlighting.
- **Lists and steps**: Use numbered lists for sequential steps and bullet points for non-sequential information.
- **Emphasis**: Use bold for UI elements (e.g., **Button**), and italics for emphasis.
- **Links**: Use descriptive link text (e.g., [Install Docker](/get-started/get-docker.md)).

For more details, see our [content formatting guidelines](/contribute/style/formatting.md)
and [grammar and style rules](/contribute/style/grammar.md).

## Best practices

- **Test all instructions**
  - Ensure all code and commands work as expected.
  - Verify that the guide can be followed successfully from start to finish.
- **Relevance**
  - Focus on real-world applications and scenarios.
  - Keep the content up-to-date with the latest Docker versions.
- **Troubleshooting tips**
  - Anticipate common issues and provide solutions or references.
- **Visual aids**
  - Include screenshots or diagrams where they enhance understanding.
  - Add captions and alt text for accessibility.
- **Third-party tools**
  - Avoid requiring the user to install third-party tools or modify their local development environment.
  - Prefer using containerized tools and methods where applicable (e.g. `docker exec`).
  - Some tools are reasonable to assume as installed or prerequisites for guides, such as Node.js and npm. Use your better judgement.

## Additional resources

- **Existing guides**
  - Refer to [Docker Guides](/guides/_index.md) for examples and inspiration.
- **Style guides**
  - [Voice and tone](/contribute/style/voice-tone.md)
  - [Content formatting](/contribute/style/formatting.md)
  - [Grammar and style](/contribute/style/grammar.md)

## Submission process

- **Proposal**
  - Raise an issue on the [Docker documentation GitHub repository](https://github.com/docker/docs)
    with a request to add a new guide.
  - Once the proposal has been accepted, start writing your guide by forking
    the repository and creating a branch for your work.

    > [!NOTE]
    > Avoid contributing from the `main` branch of your fork, since it makes it
    > more difficult for maintainers to help you fix any issues that may arise.

- **Review**
  - Proofread your guide for grammar, clarity, and adherence to the guidelines.
  - Once your draft is ready, raise a pull request, with a reference to the
    original issue.
  - The Docker documentation team and subject matter experts will review your
    submission and provide feedback on the pull request directly.
- **Publishing**
  - Your guide will be published on the Docker documentation website when the
    reviewers have approved and merged your pull request.

Thank you for contributing to the Docker community. Your expertise helps users
worldwide harness the power of Docker.
