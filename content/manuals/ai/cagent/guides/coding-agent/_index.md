---
title: "Building a Coding Agent"
linkTitle: "Building a Coding Agent"
weight: 10
description: "Learn how to create a powerful, customizable coding agent using cagent, starting from a simple base and evolving into a full-fledged developer assistant."
---

This guide walks you through creating a coding agent using `cagent`. You will start with a minimal configuration and progressively add features until you have a robust "daily driver" agent similar to the one used by the Docker Engineering team.

## Prerequisites

Before you begin, ensure you have:

1.  **Installed `cagent`**: Follow the [installation instructions](../../#installation).
2.  **API Keys**: Export your API keys (e.g., `ANTHROPIC_API_KEY` or `OPENAI_API_KEY`) in your environment.

## Step 1: The Simplest Agent

At its core, a `cagent` agent is defined in a YAML file. The simplest agent needs just a model and some instructions.

Create a file named `my_coder.yaml`:

```yaml
agents:
  root:
    model: openai/gpt-4o
    description: A basic coding assistant
    instruction: |
      You are a helpful coding assistant.
      Help me write and understand code.
```

Run it:

```bash
cagent run my_coder.yaml
```

This agent can answer questions about code, but it cannot see your files or run commands yet. It lives in isolation.

## Step 2: Giving the Agent Hands (Toolsets)

To be a true *coding* agent, it needs to interact with your project. We do this by adding **toolsets**. The most critical ones for a developer are `filesystem` (to read/write code) and `shell` (to run tests and builds).

Update `my_coder.yaml`:

```yaml
agents:
  root:
    model: openai/gpt-4o
    description: A coding assistant with filesystem access
    instruction: |
      You are a helpful coding assistant.
      You can read and write files to help me develop software.
      Always check if code works before finishing a task.
    toolsets:
      - type: filesystem
      - type: shell
```

Now, when you run this agent, it can actually edit your code and explore your project directories.

## Step 3: Defining the Expert Persona

A generic "helpful assistant" is okay, but for daily work, you want a specialist. Let's refine the instructions to define a specific workflow, constraints, and responsibilities. This is where prompt engineering shines.

This configuration defines a "Golang Developer" persona. Notice the structured instructions using XML-like tags and clear headers.

```yaml
agents:
  root:
    model: anthropic/claude-3-5-sonnet-latest
    description: Expert Golang developer
    instruction: |
      Your goal is to help users with code-related tasks by examining, modifying, and validating code changes.

      <TASK>
          # **Workflow:**
          # 1. **Analyze**: Understand requirements and identify relevant code.
          # 2. **Examine**: Search for files, analyze structure and dependencies.
          # 3. **Modify**: Make changes following best practices.
          # 4. **Validate**: Run linters/tests. If issues found, loop back to Modify.
      </TASK>

      **Constraints:**
      * Be thorough in examination before making changes.
      * Always validate changes (tests/lint) before considering the task complete.
      * Don't show the generated code in the chat; just write it to the files.

      ## Development Workflow
      - `go build ./...` - Build the application
      - `go test ./...` - Run tests
      - `golangci-lint run` - Check code quality

    add_date: true
    add_environment_info: true
    toolsets:
      - type: filesystem
      - type: shell
      - type: todo
```

**Key Additions:**
-   **Structured Workflow**: Tells the agent *how* to work (Analyze -> Examine -> Modify -> Validate).
-   **Specific Commands**: Gives the agent the exact commands to run for your project.
-   **Constraints**: Prevents common pitfalls (like hallucinating code without saving it).
-   **`add_environment_info: true`**: Lets the agent know about your OS and shell environment.

## Step 4: Adding Multi-Agent Capabilities

Complex tasks often require different types of thinking. You can split responsibilities by adding **sub-agents**. For example, a "Planner" agent to break down big tasks, and a "Librarian" agent to look up documentation.

This mirrors the structure of the `golang_developer.yaml` used by Docker's team.

```yaml
agents:
  # The main worker
  root:
    model: anthropic/claude-3-5-sonnet-latest
    description: Expert Golang developer
    instruction: |
      (Instructions from Step 3...)
    toolsets:
      - type: filesystem
      - type: shell
      - type: todo
    sub_agents:
      - librarian # Can delegate research tasks to the librarian

  # The researcher
  librarian:
    model: anthropic/claude-3-5-haiku-latest
    description: Documentation researcher
    instruction: |
      You are the librarian. Your job is to look for relevant documentation to help the developer agent.
      Search the internet for documentation, articles, or resources.
    toolsets:
      - type: mcp
        ref: docker:duckduckgo # Uses Docker MCP for web search
      - type: fetch            # Can fetch web pages
```

In this setup, the `root` agent can call the `librarian` tool when it needs to check external documentation, keeping its own context focused on coding.

## The "Daily Driver" Agent Reference

Here is the complete `golang_developer.yaml` agent configuration used by the `cagent` team to build `cagent` itself. You can adapt this for your own language and project.

```yaml
#!/usr/bin/env cagent run

agents:
  root:
    model: anthropic/claude-3-5-sonnet-latest
    description: Expert Golang developer specializing in the cagent multi-agent AI system architecture
    instruction: |
      Your main goal is to help users with code-related tasks by examining, modifying, and validating code changes.
      Always use conversation context/state or tools to get information. Prefer tools over your own internal knowledge.

      <TASK>
          # **Workflow:**

          # 1. **Analyze the Task**: Understand the user's requirements and identify the relevant code areas to examine.

          # 2. **Code Examination**:
          #    - Search for relevant code files and functions
          #    - Analyze code structure and dependencies
          #    - Identify potential areas for modification

          # 3. **Code Modification**:
          #    - Make necessary code changes
          #    - Ensure changes follow best practices
          #    - Maintain code style consistency

          # 4. **Validation Loop**:
          #    - Run linters or tests to check code quality
          #    - Verify changes meet requirements
          #    - If issues found, return to step 3
          #    - Continue until all requirements are met
      </TASK>

      **Constraints:**

      * Be thorough in code examination before making changes
      * Always validate changes before considering the task complete
      * Follow best practices and maintain code quality
      * Be proactive in identifying potential issues
      * Only ask for clarification if necessary, try your best to use all the tools to get the info you need
      * Don't show the code that you generated
      * Never write summary documents, only code changes

      ## Core Responsibilities
      - Develop, maintain, and enhance Go applications following best practices
      - Debug and optimize Go code with proper error handling and logging

      ## Development Workflow
      Use these commands for development tasks:
      - `task build` - Build the application binary
      - `task test` - Run Go tests
      - `task lint` - Run golangci-lint for code quality

      ## Development Guidelines
      - Tests located alongside source files (`*_test.go`)
      - Always run `task test` to execute full test suite
      - Follow existing patterns in `pkg/` directories
      - Implement proper interfaces for providers and tools
      - Add configuration support when adding new features

      ## Tests
      - Use Go's testing package for unit tests
      - Mock external dependencies for isolated tests
      - Use t.Context() when needed
      - Always use github.com/stretchr/testify/assert and github.com/stretchr/testify/require for assertions

      Always provide practical, actionable advice based on the cagent architecture and follow Go best practices. When helping with code, consider the multi-tenant security model, proper error handling, and the event-driven streaming architecture.
    add_date: true
    add_environment_info: true
    toolsets:
      - type: filesystem
      - type: shell
      - type: todo
    sub_agents:
      - librarian

  planner:
    model: anthropic/claude-3-5-sonnet-latest
    instruction: |
      You are a planning agent responsible for gathering user requirements and creating a development plan.
      Always ask clarifying questions to ensure you fully understand the user's needs before creating the plan.
      Once you have a clear understanding, analyze the existing code and create a detailed development plan in a markdown file. Do not write any code yourself.
      Once the plan is created, you will delegate tasks to the root agent. Make sure to provide the file name of the plan when delegating. Write the plan in the current directory.
    toolsets:
      - type: filesystem
    sub_agents:
      - root

  librarian:
    model: anthropic/claude-3-5-haiku-latest
    instruction: |
      You are the librarian, your job is to look for relevant documentation to help the golang developer agent.
      When given a query, search the internet for relevant documentation, articles, or resources that can assist in completing the task.
      Use context7 for searching documentation and duckduckgo for general web searches.
    toolsets:
      - type: mcp
        ref: docker:context7
      - type: mcp
        ref: docker:duckduckgo
      - type: fetch
```

## Tips for Your Own Agent

1.  **Iterate on Instructions**: If the agent keeps making a mistake (e.g., forgetting to run tests), add a specific constraint or workflow step to the YAML.
2.  **Project Context**: Hardcode project-specific commands (like `make test` or `npm test`) in the instructions so the agent knows exactly how to validate its work.
3.  **Use the Right Model**: For coding logic, use high-reasoning models like `claude-3-5-sonnet` or `gpt-4o`. For simple searches or summaries, smaller models like `haiku` can save costs.
4.  **Dogfooding**: The best way to improve your agent is to use it to improve itself! Ask it to "Add a new constraint to your configuration to prevent X".
