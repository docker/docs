---
title: Get started
keywords: Docker, get started, containers, AI agents, sandboxes
description: Choose a 10-minute path to run an application with Docker or give an AI agent an isolated place to work.
layout: get-started
params:
  notoc: true
  sidebar:
    hidePage: true
    include:
      - Run an application with Docker
      - Run an AI agent in a sandbox
  journeys:
    - label: Containerize an application
      title: Run an app without setting up its stack
      description: Start a frontend, API, database, and proxy with one command. Then change the code and see the result.
      link: /get-started/run-an-app/
      icon: rocket-launch
      command: docker compose watch
      time: 10 minutes
      action: Run an application
    - label: Sandbox an AI agent
      title: Give an agent room to work without handing over your machine
      description: Let an agent change and test code in its own microVM. You decide what comes back to your working tree.
      link: /get-started/run-an-agent/
      icon: shield-check
      command: sbx run --clone claude
      time: 10 minutes
      action: Run an agent in a sandbox
aliases:
  - /engine/get-started/
  - /engine/tutorials/usingdocker/
  - /guides/getting-started/get-docker-desktop/
---
