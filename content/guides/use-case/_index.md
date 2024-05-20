---
title: Use-case guides
keywords: guide, use-case
description: Learn how to accomplish tasks with Docker.
notoc: true

grid_genai:
- title: PDF analysis and chat
  description: Learn how to containerize a GenAI app that analyzes PDFs.
  icon: article
  link: /guides/use-case/genai-pdf-bot/
- title: Video transcription and chat
  description: Explore a video transcription app that lets you ask questions about videos.
  icon: play_circle
  link: /guides/use-case/genai-video-bot/
- title: Face detection with TensorFlow.js
  description: Explore a face detection app built with TensorFlow.js and Docker.
  icon: face
  link: /guides/use-case/tensorflowjs/
- title: Language translation
  description: Build and run an app to translate languages.
  link: /guides/use-case/nlp/language-translation/
  icon: translate
- title: Named entity recognition
  description: Explore an app that can identify named entities, like people, organizations, or locations.
  link: /guides/use-case/nlp/named-entity-recognition/
  icon: badge
- title: Sentiment analysis
  description: Learn how to analyze text for sentiment, like positive, negative, or neutral.
  link: /guides/use-case/nlp/sentiment-analysis/
  icon: sentiment_very_satisfied
- title: Text classification
  description: Build a text classification app and analyze the precision of your model.
  link: /guides/use-case/nlp/text-classification/
  icon: text_snippet
- title: Text summarization
  description: Explore an app that can summarize text.
  link: /guides/use-case/nlp/text-summarization/
  icon: summarize
---

Explore this collection of use-case guides designed to help you leverage Docker
for a variety of advanced applications. From generative AI to security, discover
how Docker can streamline your projects and workflows.

## Data science

{{< card
  title="Data science with JupyterLab"
  icon=analytics
  link=/guides/use-case/jupyter/
  description="Set up a JupyterLab environment to conduct and share data science research, making your projects more reproducible and scalable."
>}}

## Databases & storage

{{< card
  title="Use containerized databases"
  icon=storage
  link=/guides/use-case/databases/
  description="Learn how to run and manage containerized databases."
>}}

## Machine learning & AI

{{< grid items="grid_genai" >}}


## Security

{{< card
  title="Suppress CVEs with VEX"
  icon=verified_user
  link=/scout/guides/vex/
  description="Learn how to suppress non-applicable or fixed vulnerabilities found in your images."
>}}

