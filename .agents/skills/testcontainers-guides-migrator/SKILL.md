---
name: testcontainers-guide-migrator
description: >
  Migrate a Testcontainers guide from testcontainers.com into the Docker docs site (docs.docker.com).
  Converts AsciiDoc to Hugo Markdown, updates code to the latest Testcontainers API, splits into
  chapters with stepper navigation, verifies code compiles and tests pass, and validates against
  Docker docs style rules. Use when asked to migrate a testcontainers guide, add a TC guide, or
  port content from testcontainers.com to Docker docs.
---

# Migrate a Testcontainers Guide

You are migrating guides from https://testcontainers.com/guides/ into the Docker docs Hugo site.
Each guide lives in its own GitHub repo under `testcontainers/tc-guide-*`, written in AsciiDoc.
The source repos are listed in the testcontainers-site build.sh:
https://github.com/testcontainers/testcontainers-site/blob/main/build.sh#L23-L45

## Inputs

The user provides one or more guides to migrate. Resolve these from the inventory below:

- **REPO_NAME**: GitHub repo (e.g. `tc-guide-getting-started-with-testcontainers-for-java`)
- **SLUG**: guide slug inside `guide/` dir (e.g. `getting-started-with-testcontainers-for-java`)
- **LANG**: language identifier (go, java, dotnet, nodejs, python)
- **GUIDE_ID**: short kebab-case name (e.g. `getting-started`)

## Guide inventory

These are the 21 guides from testcontainers.com/guides/ and their source repos:

| # | Title | Repo | Lang | GUIDE_ID |
|---|-------|------|------|----------|
| 1 | Introduction to Testcontainers | tc-guide-introducing-testcontainers | (none) | introducing |
| 2 | Getting started for Java | tc-guide-getting-started-with-testcontainers-for-java | java | getting-started |
| 3 | Testing Spring Boot REST API | tc-guide-testing-spring-boot-rest-api | java | spring-boot-rest-api |
| 4 | Testcontainers lifecycle (JUnit 5) | tc-guide-testcontainers-lifecycle | java | lifecycle |
| 5 | Configuration of services in container | tc-guide-configuration-of-services-running-in-container | java | service-configuration |
| 6 | Replace H2 with real database | tc-guide-replace-h2-with-real-database-for-testing | java | replace-h2 |
| 7 | Testing ASP.NET Core web app | tc-guide-testing-aspnet-core | dotnet | aspnet-core |
| 8 | Testing Spring Boot Kafka Listener | tc-guide-testing-spring-boot-kafka-listener | java | spring-boot-kafka |
| 9 | REST API integrations with MockServer | tc-guide-testing-rest-api-integrations-using-mockserver | java | mockserver |
| 10 | Getting started for .NET | tc-guide-getting-started-with-testcontainers-for-dotnet | dotnet | getting-started |
| 11 | AWS integrations with LocalStack | tc-guide-testing-aws-service-integrations-using-localstack | java | aws-localstack |
| 12 | Testcontainers in Quarkus apps | tc-guide-testcontainers-in-quarkus-applications | java | quarkus |
| 13 | Getting started for Go | tc-guide-getting-started-with-testcontainers-for-go | go | getting-started |
| 14 | jOOQ and Flyway with Testcontainers | tc-guide-working-with-jooq-flyway-using-testcontainers | java | jooq-flyway |
| 15 | Getting started for Node.js | tc-guide-getting-started-with-testcontainers-for-nodejs | nodejs | getting-started |
| 16 | REST API integrations with WireMock | tc-guide-testing-rest-api-integrations-using-wiremock | java | wiremock |
| 17 | Local dev with Testcontainers Desktop | tc-guide-simple-local-development-with-testcontainers-desktop | java | local-dev-desktop |
| 18 | Micronaut REST API with WireMock | tc-guide-testing-rest-api-integrations-in-micronaut-apps-using-wiremock | java | micronaut-wiremock |
| 19 | Micronaut Kafka Listener | tc-guide-testing-micronaut-kafka-listener | java | micronaut-kafka |
| 20 | Getting started for Python | tc-guide-getting-started-with-testcontainers-for-python | python | getting-started |
| 21 | Keycloak with Spring Boot | tc-guide-securing-spring-boot-microservice-using-keycloak-and-testcontainers | java | keycloak-spring-boot |

Already migrated: **#2 (Java getting-started)**, **#13 (Go getting-started)**, **#20 (Python getting-started)**

## Step 0: Pre-flight

1. Confirm `testing-with-docker` tag exists in `data/tags.yaml`. If not, add:
   ```yaml
   testing-with-docker:
     title: Testing with Docker
   ```
2. Check if new terms need adding to `_vale/config/vocabularies/Docker/accept.txt`.
3. Read `STYLE.md` and `COMPONENTS.md` to refresh on Docker docs conventions.

## Step 1: Clone the guide repo

Clone the guide repo to a temporary directory. This gives you all source files locally — no HTTP calls needed.

```bash
git clone --depth 1 https://github.com/testcontainers/{REPO_NAME}.git <tmpdir>/{REPO_NAME}
```

Where `<tmpdir>` is a temporary directory on your system (e.g. the output of `mktemp -d`).

The repo structure is:
- `<tmpdir>/{REPO_NAME}/guide/{SLUG}/index.adoc` — the AsciiDoc guide source
- `<tmpdir>/{REPO_NAME}/src/` — application source code (referenced by `include::` directives)
- `<tmpdir>/{REPO_NAME}/testdata/` — test data files (SQL scripts, configs, etc.)
- `<tmpdir>/{REPO_NAME}/pom.xml` or `go.mod` — build config

1. Read `guide/{SLUG}/index.adoc` to get the guide content.
2. Find all `include::{codebase}/path/to/file[]` directives. The `{codebase}` attribute points to a remote URL, but since you have the repo cloned, read the files directly from disk instead (e.g. `include::{codebase}/src/main/java/Foo.java[]` → read `<tmpdir>/{REPO_NAME}/src/main/java/Foo.java`).
3. If includes have `[lines="X..Y"]`, extract only those lines from the local file.
4. Note the `[source,lang]` block preceding each include — that determines the code fence language.

This cloned repo also serves as the base for Step 6 (code verification) — you can run the tests directly in it to confirm they pass before updating the code to the latest API.

## Step 2: Convert AsciiDoc to Markdown

| AsciiDoc | Markdown |
|---|---|
| `== Heading` | `## Heading` |
| `=== Heading` | `### Heading` |
| `*bold*` (AsciiDoc bold) | `**bold**` |
| `https://url[Link text]` | `[Link text](url)` |
| `[source,lang]\n----\ncode\n----` | `` ```lang\ncode\n``` `` |
| `[source,shell]` with `$` prompts | `` ```console `` |
| `[NOTE]\ntext` or `====\n[NOTE]\n...\n====` | `> [!NOTE]\n> text` |
| `[TIP]\ntext` | `> [!TIP]\n> text` |
| `:toc:`, `:toclevels:`, `:codebase:` | Remove entirely |
| `include::{codebase}/path[]` | Replace with fetched code in a code fence |
| YAML front matter (date, draft, repo) | Remove; transform to Docker docs format |

## Step 3: Apply Docker docs style rules

These are mandatory (from STYLE.md and AGENTS.md):

- **No "we"**: "We are going to create" → "Create" or "Start by creating"
- **No "let us" / "let's"**: → imperative voice or "You can..."
- **No hedge words**: remove "simply", "easily", "just", "seamlessly"
- **No meta-commentary**: remove "it's worth noting", "it's important to understand"
- **No "allows you to" / "enables you to"**: → "lets you" or rephrase
- **No "click"**: → "select"
- **No bold for emphasis or product names**: only bold UI elements
- **No time-relative language**: remove "currently", "new", "recently", "now"
- **No exclamations**: remove "Voila!!!" etc.
- Use `console` language hint for interactive shell blocks with `$` prompts
- Use contractions: "it's", "you're", "don't"

## Step 4: Update code to latest Testcontainers API

Research the latest API version for the target language before writing code.

**Best practices reference**: The Testcontainers team maintains Claude skills with up-to-date API patterns and best practices for each language at https://github.com/testcontainers/claude-skills/ — check the relevant language skill (testcontainers-go, testcontainers-node, testcontainers-dotnet) for current API signatures, cleanup patterns, wait strategies, and anti-patterns to avoid.

For each language, check the cloned repo's existing code, then update to the latest API. Key patterns per language:

**Go** (testcontainers-go v0.41.0):
- `postgres.RunContainer(ctx, opts...)` → `postgres.Run(ctx, "image", opts...)`
- `testcontainers.WithImage(...)` → image is now the 2nd positional param to `Run()`
- Manual `WithWaitStrategy(wait.ForLog(...))` → `postgres.BasicWaitStrategies()`
- `t.Cleanup(func() { ctr.Terminate(ctx) })` → `testcontainers.CleanupContainer(t, ctr)`
- `if err != nil { log.Fatal(err) }` → `require.NoError(t, err)` (use testify require/assert)
- Helper functions should accept `t *testing.T` as first param, call `t.Helper()`
- No `TearDownSuite()` needed if `CleanupContainer` is registered in the helper
- Go version prerequisite: 1.25+

**Java** (testcontainers-java 2.0.4):
- Artifacts renamed in 2.x: `org.testcontainers:postgresql` → `org.testcontainers:testcontainers-postgresql`
- Check the latest version at https://java.testcontainers.org/
- Use `@Testcontainers` and `@Container` annotations for JUnit 5 lifecycle
- Prefer module-specific containers (e.g. `PostgreSQLContainer`) over `GenericContainer`
- Use `@DynamicPropertySource` for Spring Boot integration

**.NET** (testcontainers-dotnet):
- Check the latest NuGet package version
- Use `IAsyncLifetime` for container lifecycle in xUnit
- Use builder pattern: `new PostgreSqlBuilder().Build()`

**Node.js** (testcontainers-node):
- Check the latest npm version
- Use module-specific packages (e.g. `@testcontainers/postgresql`)
- Use `GenericContainer` for services without a dedicated module

**Python** (testcontainers-python):
- Check the latest PyPI version
- Use context managers (`with PostgresContainer() as postgres:`)
- Use module-specific containers when available

For all languages: consult the corresponding Testcontainers skill at https://github.com/testcontainers/claude-skills/ for current best practices and anti-patterns.

## Step 5: Create guide directory structure

Directory: `content/guides/testcontainers-{LANG}-{GUIDE_ID}/`

Each guide is its own top-level entry under `/guides/`. Do NOT nest guides inside a shared parent section — otherwise they won't appear individually in the tag/language filters on the guides listing page.

### _index.md (landing page)

```yaml
---
title: {Full guide title}
linkTitle: {Short title for guides listing}
description: {One-line description}
keywords: testcontainers, {lang}, testing, {technologies used}
summary: |
  {2-3 line summary for the guides listing card}
toc_min: 1
toc_max: 2
tags: [testing-with-docker]
languages: [{lang}]
params:
  time: {estimated} minutes
---

<!-- Source: https://github.com/testcontainers/{REPO_NAME} -->
```

Content: what you'll learn (bulleted list), prerequisites, and a NOTE linking to `https://testcontainers.com/getting-started/` for newcomers.

### Sub-pages (chapters)

Split the guide into logical chapters. Each sub-page:

```yaml
---
title: {Chapter title}
linkTitle: {Short title for stepper}
description: {One-line description}
weight: {10, 20, 30, ...}
---
```

**No `tags`, `languages`, or `params` on sub-pages** — only on `_index.md`.

Typical chapter breakdown:
| Weight | File | Content |
|--------|------|---------|
| 10 | `create-project.md` | Project setup, dependencies, business logic |
| 20 | `write-tests.md` | First test using testcontainers |
| 30 | `test-suites.md` | Reusing containers, test helpers, suites |
| 40 | `run-tests.md` | Running tests, summary, further reading |

Adapt the split to the guide's content — some guides may need fewer or more chapters.

## Step 6: Verify code compiles and tests pass

This is CRITICAL. The code in the guide MUST compile and all tests MUST pass. Do not skip this step.

### 6a: Use the cloned repo as the verification project

The repo you cloned in Step 1 (`<tmpdir>/{REPO_NAME}`) already contains a working project with all source files, build config, and tests. Use it as the starting point:

```bash
cd <tmpdir>/{REPO_NAME}
```

First, verify the **original** code compiles and tests pass before you change anything. This confirms a good baseline.

### 6b: Update the code in the cloned repo

After confirming the original works, apply the API updates (from Step 4) directly in the cloned repo's source files. This is the same code you're putting in the guide — keep them in sync.

### 6c: Update dependencies and compile

Run compilation inside a container for reproducibility — no need to install the language toolchain on the host. Use the appropriate language Docker image, mounting the cloned repo:

```bash
docker run --rm -v "<tmpdir>/{REPO_NAME}":/app -w /app <language-image> sh -c "<compile command>"
```

Pick the right image for the language (e.g. `golang:1.25-alpine`, `maven:3-eclipse-temurin-21`, `gradle:jdk21`, `mcr.microsoft.com/dotnet/sdk:9.0`, `node:22-alpine`, `python:3.13-alpine`). Update dependencies to the latest Testcontainers version and compile.

If compilation fails, fix the code and update the guide markdown to match.

### 6d: Run tests in a container with Docker socket mounted

Run tests in the same kind of container, but **mount the Docker socket** so Testcontainers can create sibling containers.

#### macOS Docker Desktop workarounds

When running on macOS with Docker Desktop, these environment variables and flags are **required**:

- **`TESTCONTAINERS_HOST_OVERRIDE=host.docker.internal`** — On macOS, containers can't reach sibling containers via the Docker bridge IP (`172.17.0.x`). This tells Testcontainers (including Ryuk) to connect via `host.docker.internal` instead. **Do NOT disable Ryuk** — it is a core Testcontainers feature and the guides must demonstrate proper usage.
- **`docker-java.properties`** with `api.version=1.47` — Docker Desktop's minimum API version is 1.44, but docker-java defaults to 1.24. Create this file in the project root and mount it to `/root/.docker-java.properties` inside Java containers.
- **`-Dspotless.check.skip=true`** — The Spotless Maven plugin in the source repos is incompatible with JDK 21. Skip it since it's a code formatter, not part of the test.
- **`-Dmicronaut.test.resources.enabled=false`** — Micronaut's Test Resources service starts a separate process that can't connect to Docker from inside a container. The guide tests use Testcontainers directly, not Test Resources. Only needed for Micronaut guides.
#### Java guide test command

```bash
# Create docker-java.properties in the project root
echo "api.version=1.47" > <tmpdir>/{REPO_NAME}/docker-java.properties

docker run --rm \
  -v "<tmpdir>/{REPO_NAME}":/app \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "<tmpdir>/{REPO_NAME}/docker-java.properties":/root/.docker-java.properties \
  -e DOCKER_HOST=unix:///var/run/docker.sock \
  -e TESTCONTAINERS_HOST_OVERRIDE=host.docker.internal \
  -w /app \
  maven:3.9-eclipse-temurin-21 \
  mvn -B test -Dspotless.check.skip=true -Dspotless.apply.skip=true
```

For Quarkus guides, use `maven:3.9-eclipse-temurin-17` instead (Quarkus 3.22.3 compiles for Java 17).

#### Go guide test command

```bash
docker run --rm \
  -v "<tmpdir>/{REPO_NAME}":/app \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e DOCKER_HOST=unix:///var/run/docker.sock \
  -e TESTCONTAINERS_HOST_OVERRIDE=host.docker.internal \
  -w /app \
  golang:1.25-alpine \
  sh -c "apk add --no-cache gcc musl-dev && go test -v -count=1 ./..."
```

#### Python guide test command

```bash
docker run --rm \
  -v "<tmpdir>/{REPO_NAME}":/app \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e DOCKER_HOST=unix:///var/run/docker.sock \
  -e TESTCONTAINERS_HOST_OVERRIDE=host.docker.internal \
  -w /app \
  python:3.13-slim \
  sh -c "pip install -r requirements.txt && python -m pytest"
```

#### .NET guide test command

```bash
docker run --rm \
  -v "<tmpdir>/{REPO_NAME}":/app \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e DOCKER_HOST=unix:///var/run/docker.sock \
  -e TESTCONTAINERS_HOST_OVERRIDE=host.docker.internal \
  -w /app \
  mcr.microsoft.com/dotnet/sdk:9.0 \
  dotnet test
```

#### Node.js guide test command

```bash
docker run --rm \
  -v "<tmpdir>/{REPO_NAME}":/app \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e DOCKER_HOST=unix:///var/run/docker.sock \
  -e TESTCONTAINERS_HOST_OVERRIDE=host.docker.internal \
  -w /app \
  node:22-alpine \
  sh -c "npm install && npm test"
```

#### Important: run tests sequentially

Run guide tests **one at a time**. Running multiple concurrent DinD or sibling-container tests can overwhelm Docker Desktop's containerd store and cause `meta.db: input/output error` corruption, requiring a Docker Desktop restart.

### 6e: Fix until green

If any test fails, debug and fix the code in both the temporary project AND the guide markdown. Re-run until all tests pass. Do not proceed until verified.

## Step 7: Update cross-references

1. **`content/manuals/testcontainers.md`**: Add a bullet under the `## Guides` section:
   ```markdown
   - [Guide title](/guides/testcontainers-{LANG}-{GUIDE_ID}/)
   ```
2. **Do NOT update** `content/guides/testcontainers-cloud/_index.md` — keep its external links.
3. Link to `https://testcontainers.com/getting-started/` for the Testcontainers overview.
4. Use internal paths for already-migrated guides; keep `testcontainers.com` links for unmigrated ones.

## Step 8: Validate

**IMPORTANT**: Run ALL validation locally before committing. Vale checks run on CI and will block the PR if they fail — fixing after push wastes CI cycles and review time.

1. `npx prettier --write content/guides/testcontainers-{LANG}-{GUIDE_ID}/`
2. `npx prettier --write content/manuals/testcontainers.md`
3. `docker buildx bake lint` — must pass with no errors
4. `docker buildx bake vale` — then check for errors in the new files:
   ```bash
   grep -A2 "testcontainers-{LANG}-{GUIDE_ID}" tmp/vale.out
   ```
   Fix ALL errors before proceeding. Common issues:
   - **Vale.Spelling**: tech terms (library names, tools) not in the dictionary → add to `_vale/config/vocabularies/Docker/accept.txt` (alphabetical order)
   - **Vale.Terms**: wrong casing (e.g. "python" → "Python") → fix in the markdown. Watch for package names like `testcontainers-python` triggering false positives — rephrase to "Testcontainers for Python" in prose.
   - **Docker.Avoid**: hedge words like "very", "simply" → reword
   - **Docker.We**: first-person plural → rewrite to "you" or imperative
   - Info-level suggestions (e.g. "VS Code" → "versus") are not blocking but review them

   Re-run `docker buildx bake vale` after fixes until no errors remain in the new files.
5. Verify in local dev server (`HUGO_PORT=1314 docker compose watch`):
   - Guide appears when filtering by its language
   - Guide appears when filtering by `Testing with Docker` tag
   - Stepper navigation works across chapters
   - All links resolve (no 404s)
6. Verify all external URLs return 200:
   ```bash
   curl -s -o /dev/null -w "%{http_code}" -L "{url}"
   ```

## Step 9: Commit

One commit per guide. Message format:
```
feat(guides): add testcontainers {lang} {guide-id} guide

Migrated from https://github.com/testcontainers/{REPO_NAME}
Updated to testcontainers-{lang} v{version} API.
```

## Special cases

- **introducing-testcontainers**: Language-agnostic, conceptual. May overlap with `content/manuals/testcontainers.md`. Review for deduplication before migrating.
- **local-dev-testcontainers-desktop**: About Testcontainers Desktop (now part of Docker Desktop). May need significant rewriting rather than mechanical migration.
- **Java guides**: Many share the same language. Each still gets its own `testcontainers-java-{GUIDE_ID}` directory.

## Reference: completed migration (Go getting-started)

Use `content/guides/testcontainers-go-getting-started/` as the reference implementation:
- `_index.md` — landing page with frontmatter, prerequisites, learning objectives
- `create-project.md` (weight: 10) — project setup and business logic
- `write-tests.md` (weight: 20) — first test with testcontainers-go
- `test-suites.md` (weight: 30) — container reuse with testify suites
- `run-tests.md` (weight: 40) — running tests, summary, further reading
