---
title: Run a coding agent safely on a Java project with Docker Sandboxes
linkTitle: Coding agent on Java project
description: Use a Docker Sandbox and a reusable sbx kit to run a coding agent on a Spring Boot project, without giving it unrestricted access to your machine.
keywords: ai, sbx, docker sandboxes, java, spring boot, testcontainers, maven, sdkman, coding agent, kit
weight: 3
summary: |
  Build a reusable sbx kit so a coding agent can build and test a Spring Boot
  project inside an isolated microVM, with a full Java toolchain and a minimal,
  reviewable network allowlist.
params:
  tags: [ai]
  featured: true
---

Coding agents are most useful when they can build and test your project:
compile the code, run the tests, pull the containers an integration test needs.
On a Java project that means a JDK, Maven, and a Docker daemon, plus network
access to Maven Central and a few image registries. The straightforward way to
give an agent all of that is to run it directly on your machine with your
credentials and your full network, which is also the part worth being careful
about.

The [Docker Sandboxes](../manuals/ai/sandboxes/_index.md) `sbx` CLI runs an
agent inside an isolated microVM with its own filesystem, its own Docker
daemon, and a network that denies everything by default. That isolation is
useful on its own, but a bare sandbox has no JDK and no Maven, and the
default-deny network blocks Maven Central. This guide closes that gap with a
[kit](../manuals/ai/sandboxes/customize/kits.md): a small, committed spec that
installs the toolchain and declares exactly which domains the project is
allowed to reach, so every teammate gets the same productive, isolated setup
from one command.

In this guide, you'll learn how to:

- Run a Java project inside a Docker Sandbox and see what's missing out of the
  box.
- Write a kit that installs SDKMAN, a JDK, and Maven, and puts them on `PATH`
  for both `sbx exec` and the agent's own shell.
- Discover the minimal network allowlist a Spring Boot and Maven workflow
  needs, instead of guessing.
- Run Testcontainers integration tests against the sandbox's built-in Docker
  daemon.
- Commit the kit so your whole team gets the same sandbox with one command.

## Assumptions

This guide assumes you are comfortable with Maven and Spring Boot, and that you
have used a coding agent before (this guide uses Claude Code). You don't need
prior experience with Docker Sandboxes or kits. Knowing roughly what
[Testcontainers](https://testcontainers.com/) does will help, since the payoff
is a passing integration test, but it isn't required to follow the steps.

## Prerequisites

- The `sbx` CLI installed and authenticated. See
  [Get started with Docker Sandboxes](../manuals/ai/sandboxes/get-started.md#install-and-sign-in)
  to install it, then run `sbx login` once and follow the browser prompt.
- The sample project. This guide uses the official Testcontainers sample,
  [`tc-guide-testing-spring-boot-kafka-listener`](https://github.com/testcontainers/tc-guide-testing-spring-boot-kafka-listener),
  updated to Java 25 and Spring Boot 4. A copy with the finished kit is at
  [shelajev/tc-guide-java-sbx-kits](https://github.com/shelajev/tc-guide-java-sbx-kits);
  clone that to follow along, or apply the same steps to your own project.

The sample is a Spring Boot service that listens for product price changes on a
Kafka topic and writes the new price to MySQL through Spring Data JPA. Its
integration test publishes an event with `KafkaTemplate` and uses Awaitility to
assert the row landed in the database, so it starts two containers, Kafka and
MySQL, plus the Testcontainers Ryuk resource reaper. Between them, these
exercise a complete, representative Java development setup — the JDK and Maven
toolchain, network access to Maven Central and image registries, and a working
Docker daemon for the containers the tests spin up — which is exactly what the
kit needs to provide.

## Run your project in a sandbox

Start with no kit, so you can see what a bare sandbox gives you. From the
project directory, launch an agent in a fresh sandbox:

```console
$ cd tc-guide-java-sbx-kits
$ sbx run claude
```

`sbx run` creates a sandbox, mounts the current directory into it, and starts
the `claude` agent attached to your terminal. The agent can read and edit your
files, and it has its own Docker daemon, but the microVM contains only a
minimal set of tools by default. Ask it to check the toolchain, or run the
check yourself through the agent's `!` shell (the `!` prefix runs a command in
the sandbox and prints the result):

```text
!java -version
bash: java: command not found
!mvn -v
bash: mvn: command not found
```

That's the gap. The agent can see the code but can't build it, and if you tried
`./mvnw test` it would also hit the default-deny network the moment Maven
reached out to Maven Central. The rest of this guide fills both gaps with a
single kit. Exit this sandbox before continuing:

```console
$ sbx rm <sandbox-name>
```

## Build the toolchain kit

A kit is a directory containing a `spec.yaml` file. It declares install
commands that run when the sandbox is built, plus a network allowlist that
applies while it runs. This guide uses a `mixin` kit, which layers extra
capabilities on top of an existing agent.

You can keep a kit anywhere. This guide puts it in a `.sbx/kits/<name>/`
directory in the project as a convention, but the location is up to you.
Checking the kit into your repository is optional — a local path works fine —
but committing it is a good way to share one setup across the team, so
everyone's sandbox is configured the same way.

Create `.sbx/kits/java-toolchain/spec.yaml`. The install section has three
steps, shown here in full and explained below. The network block comes later,
once you've discovered which domains the project needs:

```yaml
schemaVersion: "1"
kind: mixin
name: java-toolchain
displayName: Java toolchain (SDKMAN + JDK + Maven)
description: Installs SDKMAN, Temurin JDK 25, and Maven, and puts java and mvn on PATH for sbx exec and agent shell calls.

commands:
  install:
    - command: |
        set -eu
        apt-get update
        apt-get install -y zip unzip
      user: "root"
      description: Install SDKMAN prerequisites

    - command: |
        set -eu
        exec bash <<'BASH'
        set -eo pipefail

        export SDKMAN_DIR="$HOME/.sdkman"

        if [ ! -s "$SDKMAN_DIR/bin/sdkman-init.sh" ]; then
          curl -fsSL "https://get.sdkman.io?rcupdate=false" | bash
        fi

        mkdir -p "$SDKMAN_DIR/etc"
        {
          echo "sdkman_auto_answer=true"
          echo "sdkman_selfupdate_feature=false"
        } >> "$SDKMAN_DIR/etc/config"

        # shellcheck disable=SC1091
        . "$SDKMAN_DIR/bin/sdkman-init.sh"

        sdk install java 25.0.3-tem
        sdk install maven 3.9.16

        java -version
        mvn -v
        BASH
      user: "1000"
      description: Install SDKMAN, Temurin JDK 25, and Maven

    - command: |
        set -eu

        agent_home="$(getent passwd 1000 | cut -d: -f6)"
        if [ -z "$agent_home" ]; then
          agent_home="/home/agent"
        fi

        touch /etc/sandbox-persistent.sh
        if ! grep -q 'java-toolchain-kit-path' /etc/sandbox-persistent.sh; then
          cat >> /etc/sandbox-persistent.sh <<EOF

        # java-toolchain-kit-path
        export SDKMAN_DIR="$agent_home/.sdkman"
        case ":\$PATH:" in
          *:"$agent_home/.sdkman/candidates/java/current/bin":*) ;;
          *) export PATH="$agent_home/.sdkman/candidates/java/current/bin:\$PATH" ;;
        esac
        case ":\$PATH:" in
          *:"$agent_home/.sdkman/candidates/maven/current/bin":*) ;;
          *) export PATH="$agent_home/.sdkman/candidates/maven/current/bin:\$PATH" ;;
        esac
        EOF
        fi
        chmod 0644 /etc/sandbox-persistent.sh
      user: "root"
      description: Put java and mvn on PATH for every shell
```

Each entry in `commands.install` runs in order after the sandbox is created.
The `user` field selects who runs it: `root` for system changes like installing
packages, and the agent's own user (UID `1000`) for anything that should land in
the agent's home directory.

### Install the prerequisites and the toolchain

The first step installs SDKMAN's prerequisites as root. `curl` is already in the
base image, but SDKMAN also needs `zip` and `unzip`.

The second step installs the toolchain as the agent user, so SDKMAN lands in
`/home/agent/.sdkman`. A few details here aren't optional, and skipping them is
what usually breaks a headless SDKMAN install:

- `?rcupdate=false` tells the installer not to edit shell startup files. The kit
  puts the tools on `PATH` in the third step, deterministically.
- `sdkman_auto_answer=true` makes `sdk install` answer its own prompts, and
  `sdkman_selfupdate_feature=false` stops it trying to self-update mid-build.
  Without these, the install hangs waiting for input that never comes.
- The versions are pinned: `25.0.3-tem` is the Temurin JDK 25 build, matching
  the sample's `<java.version>25</java.version>`, and `3.9.16` is Maven. Pinning
  means every teammate's sandbox has the same toolchain.

The closing `java -version` and `mvn -v` make the build fail loudly if either
tool didn't install, so a broken kit never reaches a green checkmark.

### Put the toolchain on PATH

Installing the tools isn't enough. A sandbox runs commands two ways: the agent's
own `!` shell, and `sbx exec` from your host. Both start fresh shells that don't
read SDKMAN's rc-file hooks, so without help neither finds `java`. The sandbox
sources `/etc/sandbox-persistent.sh` before every command in both paths, so
that's where the `PATH` entries go, and that's what the third step writes. This
is the standard pattern for version managers like SDKMAN and nvm; see
[Customize the shell environment](../manuals/ai/sandboxes/customize/kit-examples.md#customize-the-shell-environment)
for the same technique applied to nvm.

The third step appends SDKMAN's `current/bin` directories for the JDK and Maven
straight to `PATH`. The `grep` guard with the `java-toolchain-kit-path` marker
keeps the block from being added twice, and the `case` statements keep each
directory from being added to `PATH` more than once — that file is sourced
before every single command, and a growing `PATH` adds up fast.

> [!WARNING]
> Add only the `PATH` exports, never SDKMAN's bash-completion script. Because
> `/etc/sandbox-persistent.sh` is sourced before every command and not only at
> shell startup, a completion script (which expects `COMP_WORDS` and related
> variables to exist) breaks every subsequent command, and the shell goes
> silent. Sourcing `sdkman-init.sh` here would also give you the `sdk` command
> in every shell, at the cost of running SDKMAN's init on every command. Point
> `PATH` at the binaries directly instead: `java` and `mvn` are what the agent
> needs for normal work, and if you want `sdk` later you can run
> `source "$SDKMAN_DIR/bin/sdkman-init.sh"` by hand.

### Validate the kit

Before running anything, check the spec:

```console
$ sbx kit validate .sbx/kits/java-toolchain
Kit "java-toolchain" is valid.
```

With a valid kit, you add the network block next.

## Open the network for a Spring Boot and Maven workflow

By default, a sandbox denies every outbound network request that isn't on an
allowlist. Rather than guessing the full list, let the workflow tell you what it
needs. Build the sandbox with the kit so far (no network block yet) and try the
install or a build inside it. The first outbound request fails with a structured
block:

```text
!curl -sS https://repo.maven.apache.org
Blocked by network policy: domain repo.maven.apache.org
  detail: no matching allow rule — blocked by default deny policy
```

Inspect what's been blocked from your host:

```console
$ sbx policy log <sandbox-name>
```

`sbx policy log` shows each connection with its host, the rule that matched, and
the reason. Working through the install and a test run, the blocks appear in the
order the tools reach out: SDKMAN's API, then the JDK download broker, then
Maven Central, then Docker Hub when Testcontainers pulls images. Allow a domain
for the running sandbox with `sbx policy allow network <domain>`, or approve it
in the `sbx` TUI, then re-run and watch for the next block. Repeat until the
build and test pass.

<!-- TODO: add screenshot of the sbx TUI network-approval view, showing a blocked domain and the allow action. -->

Each domain you confirm becomes a candidate for your kit spec. Once you've
identified all the required domains, add a `network` block to `spec.yaml`:

```yaml
network:
  allow:
    # SDKMAN + JDK/Maven distribution
    - "get.sdkman.io:443"
    - "api.sdkman.io:443"
    - "broker.sdkman.io:443"
    - "github.com:443"
    - "release-assets.githubusercontent.com:443"

    # Maven Central
    - "repo.maven.apache.org:443"

    # Testcontainers images from Docker Hub: cp-kafka, mysql, and Ryuk
    - "auth.docker.io:443"
    - "registry-1.docker.io:443"
    - "production.cloudfront.docker.com:443"
```

This is the whole list for this project: SDKMAN and its JDK broker, GitHub and
its release-asset host (where SDKMAN fetches the Temurin build), Maven Central,
and the three Docker Hub endpoints an image pull touches. It's short enough to
review in a code review, which is the point. Anything not on it stays blocked,
and you can see exactly what the agent is allowed to reach.

> [!NOTE]
> Your host might already have broader policy rules (for GitHub or Docker Hub,
> say) from other work. A domain those rules already cover won't show up as a
> block. The kit still declares it explicitly so the allowlist is complete on a
> teammate's machine that doesn't have those host rules.

## Run the Testcontainers integration tests

Recreate the sandbox with the finished kit and start the agent:

```console
$ sbx run claude --kit .sbx/kits/java-toolchain
```

The toolchain is on `PATH` for both entry points. Verify from your host with
`sbx exec`:

```console
$ sbx exec <sandbox-name> bash -lc 'java -version && mvn -v'
openjdk version "25.0.3" 2026-xx-xx
OpenJDK Runtime Environment Temurin-25.0.3+x (build 25.0.3+x)
OpenJDK 64-Bit Server VM Temurin-25.0.3+x (build 25.0.3+x, mixed mode, sharing)
Apache Maven 3.9.16
```

and from inside the agent with `!java -version`. Both find the toolchain, which
is the proof the `PATH` step worked across both shells.

Run the integration test. Testcontainers talks to the sandbox's built-in Docker
daemon, so there's nothing extra to configure:

```text
!./mvnw test
```

On the first run, Maven downloads dependencies from Maven Central, and
Testcontainers pulls `confluentinc/cp-kafka:7.6.1`, `mysql:8.0.32`, and the Ryuk
reaper image from Docker Hub, all over the domains the kit allows. The
containers start and the test runs:

```text
[INFO] Creating container for image: testcontainers/ryuk:...
[INFO] Container testcontainers/ryuk:... started
[INFO] Creating container for image: confluentinc/cp-kafka:7.6.1
[INFO] Container confluentinc/cp-kafka:7.6.1 started
[INFO] Creating container for image: mysql:8.0.32
[INFO] Container mysql:8.0.32 started
[INFO] Tests run: 1, Failures: 0, Errors: 0, Skipped: 0
[INFO] BUILD SUCCESS
```

The agent built and tested a real Spring Boot service, with live Kafka and MySQL
containers, entirely inside the sandbox.

To see the other half — that the sandbox is still closed — try a domain the kit
doesn't allow:

```text
!curl -sS https://example.com
Blocked by network policy: domain example.com
  detail: no matching allow rule — blocked by default deny policy
```

The agent has exactly what the project needs and nothing more.

## Make it reproducible for the team

The kit lives in the repository, so committing it turns this from a one-off
setup into something the whole team shares. A teammate clones the repository,
launches the sandbox with the kit, and gets the same JDK, the same Maven, and
the same network allowlist, in their own isolated microVM. The setup is reviewed
like any other code, because it is code.

### Give the agent context about the kit

The agent starts fresh in each sandbox and doesn't automatically know what the
kit installed. Add an `agentContext` block to `spec.yaml`; the sandbox appends
it to the agent's memory file at creation, so the agent starts already knowing
the toolchain is there and how to use it:

```yaml
agentContext: |
  This sandbox has a Java toolchain installed and on PATH: Temurin JDK 25 and
  Maven, managed by SDKMAN. Build and test the project with `./mvnw`.
  Testcontainers uses the sandbox's built-in Docker daemon, so integration
  tests run with `./mvnw test`.
```

This is the place for kit-specific notes the agent should act on, and it travels
with the kit — so you don't have to repeat it in the project's `CLAUDE.md`.

### Document how to start

The one thing a person needs is the launch command, so put it in the project's
`README` where a teammate will look for it:

````md
## Run this project in a Docker Sandbox

Start a coding agent in a sandbox with the committed kit:

```console
$ sbx run claude --kit .sbx/kits/java-toolchain
```
````

A teammate clones the repository, runs that one command, and gets the same JDK,
the same Maven, and the same network allowlist, in their own isolated microVM.

> [!TIP]
> To steer teammates away from running a coding agent directly on their host,
> add an instruction to the project's agent instructions (for example,
> `CLAUDE.md`) that tells the agent to stop and relaunch in a sandbox:
>
> ```md
> If you are not running inside a Docker Sandbox, stop and tell the user to
> relaunch you with `sbx run claude --kit .sbx/kits/java-toolchain`.
> ```
>
> It's opinionated, but it turns "use the sandbox" from a suggestion into
> something the agent helps enforce.

## Conclusion

You have a coding agent that is productive on a Java project, with a full JDK
and Maven toolchain and a working Docker daemon for Testcontainers, running in
an isolated microVM whose network only reaches the handful of domains the
project actually uses. The kit makes that combination reproducible: it's
committed to the repository, readable in a pull request, and applied with one
command, so the safe and productive setup is the default for everyone on the
team rather than something each person reconstructs by hand.

This kit is deliberately composite. It bundles SDKMAN, the JDK, and Maven
together because they're always needed as a unit for this project. If your team
has kits for other stacks, it might make sense to factor them into smaller,
composable kits instead.

The finished sample, including the kit and the `README` section, is at
[shelajev/tc-guide-java-sbx-kits](https://github.com/shelajev/tc-guide-java-sbx-kits).

## Further reading

- [Docker Sandboxes overview](../manuals/ai/sandboxes/_index.md) for what the
  isolation gives you and how `sbx` works.
- [Customize sandboxes with kits](../manuals/ai/sandboxes/customize/kits.md) for
  what kits can do and how to build one.
- [Kit spec reference](../manuals/ai/sandboxes/customize/kit-reference.md) for
  every field in `spec.yaml`, including `network`, `agentContext`, and command
  users.
- [Testing Spring Boot Kafka listeners with Testcontainers](https://testcontainers.com/guides/testing-spring-boot-kafka-listener-using-testcontainers/),
  the original guide behind the sample app.
- [sbx-moderne-kit](https://github.com/shelajev/sbx-moderne-kit) for a more
  involved kit that installs a different toolchain using the same
  PATH-persistence pattern.
