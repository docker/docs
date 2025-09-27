---
title: cagent examples
description: Get inspiration from agent examples
keywords: [ai, agent, cagent]
weight: 10
---

Get inspiration from the following agent examples.

## Agentic development team

```yaml {title="dev-team.yaml"}
agents:
  root:
    model: claude
    description: Technical lead coordinating development
    instruction: |
      You are a technical lead managing a development team.
      Coordinate tasks between developers and ensure quality.
    sub_agents: [developer, reviewer, tester]

  developer:
    model: claude
    description: Expert software developer
    instruction: |
      You are an expert developer. Write clean, efficient code
      and follow best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think

  reviewer:
    model: gpt4
    description: Code review specialist
    instruction: |
      You are a code review expert. Focus on code quality,
      security, and maintainability.
    toolsets:
      - type: filesystem

  tester:
    model: gpt4
    description: Quality assurance engineer
    instruction: |
      You are a QA engineer. Write tests and ensure
      software quality.
    toolsets:
      - type: shell
      - type: todo

models:
  gpt4:
    provider: openai
    model: gpt-4o

  claude:
    provider: anthropic
    model: claude-sonnet-4-0
    max_tokens: 64000
```

## Research assistant

```yaml {title="research-assistant.yaml"}
agents:
  root:
    model: claude
    description: Research assistant with web access
    instruction: |
      You are a research assistant. Help users find information,
      analyze data, and provide insights.
    toolsets:
      - type: mcp
        command: mcp-web-search
        args: ["--provider", "duckduckgo"]
      - type: todo
      - type: memory
        path: "./research_memory.db"

models:
  claude:
    provider: anthropic
    model: claude-sonnet-4-0
    max_tokens: 64000
```

## Technical blog writer

```yaml {title="tech-blog-writer.yaml"}
#!/usr/bin/env cagent run
version: "1"

agents:
  root:
    model: anthropic
    description: Writes technical blog posts
    instruction: |
      You are the leader of a team of AI agents for a technical blog writing workflow.

      Here are the members in your team:
      <team_members>
      - web_search_agent: Searches the web
      - writer: Writes a 750-word technical blog post based on the chosen prompt
      </team_members>

      <WORKFLOW>
        1. Call the `web_search_agent` agent to search the web to get
           important information about the task that is asked

        2. Call the `writer` agent to write a 750-word technical blog
           post based on the research done by the web_search_agent
      </WORKFLOW>

      - Use the transfer_to_agent tool to call the right agent at the right
        time to complete the workflow.
      - DO NOT transfer to multiple members at once
      - ONLY CALL ONE AGENT AT A TIME
      - When using the `transfer_to_agent` tool, make exactly one call
        and wait for the result before making another. Do not batch or
        parallelize tool calls.
    sub_agents:
      - web_search_agent
      - writer
    toolsets:
      - type: think

  web_search_agent:
    model: anthropic
    add_date: true
    description: Search the web for information
    instruction: |
      Search the web for information

      Always include sources
    toolsets:
      - type: mcp
        command: uvx
        args: ["duckduckgo-mcp-server"]

  writer:
    model: anthropic
    description: Writes a 750-word technical blog post based on the chosen prompt.
    instruction: |
      You are an agent that receives a single technical writing prompt
      and generates a detailed, informative, and well-structured technical blog post.

      - Ensure the content is technically accurate and includes relevant
        code examples, diagrams, or technical explanations where appropriate.
      - Structure the blog post with clear sections, including an introduction,
        main content, and conclusion.
      - Use technical terminology appropriately and explain complex concepts clearly.
      - Include practical examples and real-world applications where relevant.
      - Make sure the content is engaging for a technical audience while
        maintaining professional standards.

      Constraints:
      - DO NOT use lists

models:
  anthropic:
    provider: anthropic
    model: claude-3-5-sonnet-latest
```

See more examples in the [repository](https://github.com/docker/cagent/tree/main/examples).