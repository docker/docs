# Cloud sandbox workflow examples

This directory is the prepared content for a companion repository. Publish its
source files, README files, `.gitignore`, and `LICENSE`, excluding generated
output and caches. No remote repository has been created. The documentation
guides remain Hugo drafts until the live workflows have been verified.

- `api-change-reporter/` contains the optional reference implementation for
  the build guide and the independent starting point for the transfer guide
- `express-upgrade/` contains the historical Express 4.21.2 baseline with a
  lockfile and tests for the upgrade investigation
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
the Express baseline, the upgrade failure, the incomplete wildcard fix, and
the complete fix. A deliberately failing baseline produces a blocked report
without changing the dependency or lockfile.

Local output is evidence of application behavior. It does not verify Console
authentication, public ports, GitHub delivery, email delivery, or sandbox moves.

The recorded local run used Python 3.14.8, PyYAML 6.0.3, and Node.js 24.21.0.
The Express 5.1.0 installation resolved router 2.2.0 and path-to-regexp 8.4.2.
The HTML report was rendered in Chromium. Its ZIP was downloaded through a local
HTTP server and verified after that server stopped.

## Complete live verification

Before removing `draft: true` from the guides:

- Publish this directory as a companion repository and record its baseline
  commit. Replace the companion setup notes in the guides with that verified
  URL and revision. Retain placeholders for the reader's fork and issue.
- Run each journey using Claude Code and an Anthropic API key. Record the kit
  reference and immutable version, agent version, CLI version where relevant,
  platform, selected policies, required hosts, prompts, and actual outputs.
- In the build journey, capture the Console terminal and published report.
  Download the ZIP through the public URL, close the port, remove the sandbox,
  and verify that the downloaded source and tests remain usable.
- In the upgrade journey, verify successful and blocked reports reach the
  selected GitHub issue. Capture the issue comment and its URL. If including
  the email extension, send a test message to the Resend account owner's
  address and verify receipt before adding delivery steps.
- In the transfer journey, verify the project, virtual environment, and
  `HANDOFF.md` arrive at the destination. Disconnect during cloud work, return,
  and copy and test the result bundle. Capture the destination identity,
  expiration, and retrieval. Verify cleanup of both sandboxes separately.

Store guide screenshots in `content/guides/images/`, prefixed with the guide
slug. The report screenshot included in the build draft is from the local
reference implementation. Replace it with the verified cloud run if its output
differs. Do not present local output as a Console session or an issue comment.
