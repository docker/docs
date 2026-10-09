# Cloud sandbox workflow examples

This directory holds local examples for reproducing the guides' outputs and
preparing the transfer workflow. The build guide starts from an idea and doesn't
require a companion repository. The Express example is also available in a
private repository for author testing. The guides render in previews while the
pull request remains a draft pending live verification.

- `api-change-reporter/` reproduces the illustrated report and provides the
  starting point for the transfer guide
- `express-upgrade/` contains a historical Express 4.21.2 support portal with
  a lockfile and eight tests for the migration assessment
- `verify.py` reproduces application behavior locally without sandbox access

## Reproduce local evidence

Use Python 3.11 or later with virtual environment support and Node.js 22 or
later. From this directory, run:

```console
$ python3 verify.py
```

The verifier installs dependencies in temporary directories and writes evidence
to the ignored `output/` directory. It checks eight reporter tests, generates
an HTML report and ZIP, extracts the ZIP, and reruns the tests. It also checks
the Express baseline, startup failure, root-route behavior, nested filters, and
empty request bodies. A deliberately failing baseline produces a blocked report
without changing the dependency or lockfile.

Local output is evidence of application behavior. It does not verify Console
authentication, public ports, GitHub issue creation, or sandbox moves.

The recorded local run used Python 3.14.8, PyYAML 6.0.3, and Node.js 24.21.0.
The Express 5.1.0 installation resolved router 2.2.0 and path-to-regexp 8.4.2.
The HTML report was rendered in Chromium. Its ZIP was downloaded through a local
HTTP server and verified after that server stopped.

## Complete live verification

Before publishing the guides:

- Publish the transfer guide's starting project and record its baseline commit.
  Replace its setup notes with a verified URL and revision. The build guide
  doesn't need a published project.
- Run each journey using Claude Code and an Anthropic API key. Record the kit
  reference and immutable version, agent version, CLI version where relevant,
  platform, selected policies, required hosts, prompts, and actual outputs.
- In the build journey, capture the Console terminal and published report.
  Download the ZIP through the public URL, close the port, remove the sandbox,
  and verify that the downloaded source and tests remain usable.
- In the migration assessment, capture creation of a concise tracking issue,
  including its recommendation, compatibility evidence, and remaining checks.
  Verify that a failing baseline produces a blocked assessment. Test reuse
  of an existing tracking issue before documenting it as observed behavior.
- In the transfer journey, verify the project, virtual environment, and
  `HANDOFF.md` arrive at the destination. Disconnect during cloud work, return,
  and copy and test the result bundle. Capture the destination identity,
  expiration, and retrieval. Verify cleanup of both sandboxes separately.

Store guide screenshots in `content/guides/images/`, prefixed with the guide
slug. The report screenshots in the build guide were rendered locally, including
the dark color scheme used for the refinement. Replace them with the verified
cloud run if its output differs. Do not present local output as a Console session
or a GitHub issue.
