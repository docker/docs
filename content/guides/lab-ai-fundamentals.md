---
title: "Lab: AI Fundamentals for Developers"
linkTitle: "Lab: AI Fundamentals"
description: |
  Learn the core concepts of AI development — models, prompt engineering, tool
  calling, and RAG — through hands-on exercises in a live environment.
summary: |
  Hands-on lab: Learn the four core pillars of AI application development.
  Work with the Chat Completions API, prompt engineering, tool calling, and RAG
  through interactive exercises.
keywords: AI, Docker, Model Runner, prompt engineering, RAG, tool calling, lab, labspace
aliases:
  - /labs/docker-for-ai/ai-fundamentals/
params:
  tags: [ai, labs]
  time: 45 minutes
  resource_links:
    - title: Docker Model Runner docs
      url: /ai/model-runner/
    - title: Labspace repository
      url: https://github.com/dockersamples/labspace-ai-fundamentals
---

Get hands-on with the four core pillars of AI application development: models,
prompt engineering, tool calling, and RAG. This lab runs entirely on your
machine using Docker Model Runner — no API key or cloud account required.

## Launch the lab

{{< labspace-launch image="dockersamples/labspace-ai-fundamentals" model-download="true" >}}

## What you'll learn

By the end of this Labspace, you will have completed the following:

- Understand the Chat Completions API and how to structure messages for a model
- Use prompt engineering techniques including system prompts, few-shot examples, and structured output
- Implement tool calling and the agentic loop in code
- Build a RAG pipeline that grounds model responses in your own data

## Modules

| #   | Module                               | Description                                                       |
| --- | ------------------------------------ | ----------------------------------------------------------------- |
| 1   | Welcome & Setup                      | Introduction to the lab and verifying your environment            |
| 2   | Talking to Models                    | Chat Completions API, message roles, and stateless model behavior |
| 3   | Prompt Engineering                   | System prompts, few-shot examples, and structured output          |
| 4   | Tool Calling                         | Tool definitions, the agentic loop, and executing tools in code   |
| 5   | Retrieval Augmented Generation (RAG) | Retrieve, augment, and generate with your own knowledge base      |
| 6   | Wrap-up                              | Summary of concepts and next steps                                |
