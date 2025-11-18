---
title: Customize a code quality check workflow
linkTitle: Customize workflow
summary: Adapt your GitHub and SonarQube workflow to focus on specific quality issues, integrate with CI/CD, and set custom thresholds.
description: Learn how to customize prompts for specific quality issues, filter by file patterns, set quality thresholds, and integrate your workflow with GitHub Actions for automated code quality checks.
weight: 20
---

Now that you understand the basics of automating code quality workflows with
GitHub and SonarQube in E2B sandboxes, you can customize the workflow
for your needs.

## Focus on specific quality issues

Modify the prompt to prioritize certain issue types:

{{< tabs group="language" >}}
{{< tab name="TypeScript" >}}

```typescript
const prompt = `Using SonarQube and GitHub MCP tools:

Focus only on:
- Security vulnerabilities (CRITICAL priority)
- Bugs (HIGH priority)
- Skip code smells for this iteration

Analyze "${repoPath}" and fix the highest priority issues first.`;
```

{{< /tab >}}
{{< tab name="Python" >}}

```python
prompt = f"""Using SonarQube and GitHub MCP tools:

Focus only on:
- Security vulnerabilities (CRITICAL priority)
- Bugs (HIGH priority)
- Skip code smells for this iteration

Analyze "{repo_path}" and fix the highest priority issues first."""
```

{{< /tab >}}
{{< /tabs >}}

## Integrate with CI/CD

Add this workflow to GitHub Actions to run automatically on pull requests:

{{< tabs group="language" >}}
{{< tab name="TypeScript" >}}

```yaml
name: Automated quality checks
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: "18"
      - run: npm install
      - run: npx tsx 06-quality-gated-pr.ts
        env:
          E2B_API_KEY: ${{ secrets.E2B_API_KEY }}
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          SONARQUBE_TOKEN: ${{ secrets.SONARQUBE_TOKEN }}
          GITHUB_OWNER: ${{ github.repository_owner }}
          GITHUB_REPO: ${{ github.event.repository.name }}
          SONARQUBE_ORG: your-org-key
```

{{< /tab >}}
{{< tab name="Python" >}}

```yaml
name: Automated quality checks
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with:
          python-version: "3.8"
      - run: pip install e2b python-dotenv
      - run: python 06_quality_gated_pr.py
        env:
          E2B_API_KEY: ${{ secrets.E2B_API_KEY }}
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          SONARQUBE_TOKEN: ${{ secrets.SONARQUBE_TOKEN }}
          GITHUB_OWNER: ${{ github.repository_owner }}
          GITHUB_REPO: ${{ github.event.repository.name }}
          SONARQUBE_ORG: your-org-key
```

{{< /tab >}}
{{< /tabs >}}

## Filter by file patterns

Target specific parts of your codebase:

{{< tabs group="language" >}}
{{< tab name="TypeScript" >}}

```typescript
const prompt = `Analyze code quality but only consider:
- Files in src/**/*.js
- Exclude test files (*.test.js, *.spec.js)
- Exclude build artifacts in dist/

Focus on production code only.`;
```

{{< /tab >}}
{{< tab name="Python" >}}

```python
prompt = """Analyze code quality but only consider:
- Files in src/**/*.js
- Exclude test files (*.test.js, *.spec.js)
- Exclude build artifacts in dist/

Focus on production code only."""
```

{{< /tab >}}
{{< /tabs >}}

## Set quality thresholds

Define when PRs should be created:

{{< tabs group="language" >}}
{{< tab name="TypeScript" >}}

```typescript
const prompt = `Quality gate thresholds:
- Only create PR if:
  * Bug count decreases by at least 1
  * No new security vulnerabilities introduced
  * Code coverage does not decrease
  * Technical debt reduces by at least 15 minutes

If changes do not meet these thresholds, explain why and skip PR creation.`;
```

{{< /tab >}}
{{< tab name="Python" >}}

```python
prompt = """Quality gate thresholds:
- Only create PR if:
  * Bug count decreases by at least 1
  * No new security vulnerabilities introduced
  * Code coverage does not decrease
  * Technical debt reduces by at least 15 minutes

If changes do not meet these thresholds, explain why and skip PR creation."""
```

{{< /tab >}}
{{< /tabs >}}

## Next steps

Learn how to troubleshoot common issues.
