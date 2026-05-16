---
title: "Lab: Building an AI Product Reviewer"
linkTitle: "Lab: Building an AI Product Reviewer"
description: |
  Build a complete AI-powered feedback analysis pipeline — sentiment analysis,
  semantic clustering with embeddings, and response generation — all running
  locally via Docker Model Runner. No API keys required.
summary: |
  Hands-on lab: Build an AI pipeline that classifies product reviews by
  sentiment, clusters them by topic using embeddings, extracts actionable
  features, and generates context-aware responses — all running locally.
keywords: AI, Docker, Model Runner, sentiment analysis, embeddings, RAG, lab, labspace
aliases:
  - /labs/docker-for-ai/creating-ai-product-reviewer/
params:
  tags: [ai, labs]
  time: 60 minutes
  resource_links:
    - title: Docker Model Runner docs
      url: /ai/model-runner/
    - title: Labspace repository
      url: https://github.com/dockersamples/labspace-creating-ai-product-reviewer
---

Build a complete feedback analysis pipeline for a fictional AI product called
Jarvis. You'll write Node.js code that runs local LLMs and embedding models via
Docker Model Runner — no API keys, no cloud subscriptions, no data leaving your
machine.

## Launch the lab

{{< labspace-launch image="dockersamples/labspace-creating-ai-product-reviewer" >}}

## What you'll learn

By the end of this Labspace, you will have completed the following:

- Run LLMs locally via Docker Model Runner's OpenAI-compatible API
- Connect a Node.js app to Docker Model Runner using the OpenAI SDK and the Compose `models:` integration
- Perform sentiment analysis using low-temperature LLM classification
- Use embeddings and cosine similarity to cluster semantically related feedback
- Extract structured data from an LLM using `response_format: { type: 'json_object' }`
- Generate context-aware responses to reviews informed by extracted product features

## Modules

| #   | Module                              | Description                                                                        |
| --- | ----------------------------------- | ---------------------------------------------------------------------------------- |
| 1   | Introduction                        | Overview of the pipeline and Docker Model Runner setup                             |
| 2   | Project Setup & Docker Model Runner | Explore the starter project and wire up Compose model integration                  |
| 3   | Generating Synthetic Feedback       | Use the LLM to generate realistic product reviews as test data                     |
| 4   | Sentiment Analysis                  | Classify reviews as positive, negative, or neutral with low-temperature generation |
| 5   | Embeddings & Semantic Clustering    | Group related reviews using vector embeddings and cosine similarity                |
| 6   | Features & Responses                | Extract actionable features and generate context-aware review responses            |
| 7   | Wrap-up                             | Summary of techniques and ideas for extending the pipeline                         |
