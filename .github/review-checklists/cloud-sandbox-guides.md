# Verify the cloud sandbox guides

Use this checklist to review the three Docker Agentic Platform guides before
publication. All live checks are pending. Drafting used documentation references;
commands, Console labels, editor behavior, and version requirements still need
live verification. Source presence doesn't establish account availability.

Run each journey from its prerequisites through cleanup. For each row, replace
Pending with Pass, Fail, or Blocked, and record the date, observed result, and
supporting output or screenshot in the observations cell. Remove credentials,
tokens, and private data from evidence. Correct the guide after any discrepancy
and repeat the affected checks. A blocked or failed required check prevents
publication.

## Record the test environment

Complete a separate record for each guide and repeat it if the environment changes.

| Field | Share an app | Remote development | Local-to-cloud handoff |
| --- | --- | --- | --- |
| Reviewer and test date | Pending | Pending | Pending |
| Docs commit tested | Pending | Pending | Pending |
| Account capabilities: subscription, cloud entitlement, available architectures and compute sizes | Pending | Pending | Pending |
| Account network restrictions and credential scopes, without secret values | Pending | Pending | Pending |
| CLI version and installation method, or not used | Pending | Pending | Pending |
| Kit name and version or image digest | Pending | Pending | Pending |
| Local OS, architecture, and runtime version, or not used | Pending | Pending | Pending |
| Cloud architecture and compute size | Pending | Pending | Pending |
| Browser version; SSH client and editor/extension versions where used | Pending | Pending | Pending |
| Sample repository full commit ID | Pending | Pending | Pending |
| Start/end times and measured active and elapsed duration | Pending | Pending | Pending |
| Corrections made and retest evidence | Pending | Pending | Pending |

Don't infer minimum supported versions from one successful run. Check the
0.45.0-or-later requirement in the remote development and handoff guides against
the supported CLI release information and record that evidence. Record any
missing access as Blocked rather than treating an untested claim as confirmed.

## Share an app

Guide: [Share an application from a cloud sandbox](../../content/manuals/agentic-platform/guides/share-app.md).

| Guide section and claim to test | Expected result | Status | Date, observations, and evidence |
| --- | --- | --- | --- |
| Introduction: provisional 20-minute duration | Measure the entire journey and replace the provisional estimate with a measured duration in the guide. | Pending | Pending |
| Before you begin: entitlement, credentials, pricing, and timers | Account can launch the Shell kit without GitHub/model credentials or MCP tools. Console pricing and timeout choices support the guidance. Record compute charging rules and distinguish them from observed usage. | Pending | Pending |
| Step 1: Console launch controls and network policy | The documented labels and Balanced/Open choices exist. Shell starts with the chosen size, duration, and timeout action; no saved credentials or MCP tools are selected. Local files aren't mounted. | Pending | Pending |
| Step 2: clone, build, run, and inspect commands | All commands succeed from the stated directory. Record the sample revision. Container runs with port 3000 mapped inside the sandbox; the mapping alone doesn't publish a cloud URL. | Pending | Pending |
| Step 3: Connect, Connect via a Port, and Open Port | The documented UI publishes port 3000 and returns a working HTTPS URL. A separate browser session signed out of Docker displays the sample. Record that the app has no application authentication and contains no sensitive data. | Pending | Pending |
| Step 3: disconnect behavior | Closing the original browser leaves the sample reachable until an explicit stop or timeout action. | Pending | Pending |
| Step 4: withdraw public access | Closing port 3000 makes the application unreachable on reload in the separate session. Record the observed response and any delay. Confirm the documented container stop command succeeds. | Pending | Pending |
| Step 4: save files and credential scope | Project remains at the documented path. Copy it outside the sandbox using the linked transfer instructions and inspect the copy. If testing the Git alternative, use only the required repository permissions and record them without secrets. | Pending | Pending |
| Step 5: pause, resume, timeout, charges, and deletion | Pause retains files and stops compute charging under the documented billing rules. Resume starts a fresh timer with the original duration. Exercise Stop and Delete timeout actions on disposable sandboxes; confirm retention versus removal. Save work before deletion and confirm the sandbox is absent. | Pending | Pending |

## Remote development

Guide: [Develop remotely in a cloud sandbox](../../content/manuals/agentic-platform/guides/remote-development.md).

Run the editor-specific checks with both VS Code and Cursor. Record versions and
observations separately for each editor, including its Remote - SSH extension.

| Guide section and claim to test | Expected result | Status | Date, observations, and evidence |
| --- | --- | --- | --- |
| Introduction and prerequisites: provisional 25-minute duration and CLI requirement | Measure each editor journey and replace provisional timing. Record evidence for the minimum CLI version, cloud entitlement, and use of the same Docker account in Console and CLI. | Pending | Pending |
| Step 1: Console controls, credential scope, and clone | Shell launches with Open, without saved credentials or MCP tools. Confirm the editor server and dependencies download with this policy. Stop is available as the timeout action. Clone succeeds without a GitHub token and creates source and Dockerfile at the stated remote path. | Pending | Pending |
| Step 2: login, SSH setup, list, connect, and exit | Each command succeeds. The sbx_-prefixed ID works as the SSH user with sbx_cloud; project files are visible and exit returns to the local terminal. Record any SSH setup prompts or prerequisites missing from the guide. | Pending | Pending |
| Step 3: VS Code connection | The named extension and Command Palette action connect using the tested ID. Remote server installs inside the sandbox and the remote window opens the documented directory. Record exact UI labels and extension version. | Pending | Pending |
| Step 3: Cursor connection | Cursor's Remote - SSH support and Command Palette action connect using the tested ID. Remote server installs inside the sandbox and the remote window opens the documented directory. Record exact extension identity, UI labels, and version. | Pending | Pending |
| Step 4: edit, build, run, preview, and withdraw in each editor | Saved heading change is in the remote files. Build/run commands succeed; public HTTPS preview shows the edit. Closing the public port withdraws access on reload. No sensitive data is exposed. | Pending | Pending |
| Step 5: disconnect and reconnect in each editor | Closing the window leaves the sandbox running. Reconnection opens the same files and saved change. After a timer stop, resume permits reconnection with a fresh timer; restart the workload if required and record that behavior. | Pending | Pending |
| Step 5: copy and optional Git credential | The documented cp command creates a local project copy containing the edit. Files don't sync automatically. If testing the Git alternative, record required repository permissions without secrets. | Pending | Pending |
| Before you begin and step 6: file/process location, timers, charges, and cleanup | Editor window runs locally; project, terminal processes, and remote server live in the sandbox. Confirm billing guidance for running/disconnected and paused states. Pause retains files, resume refreshes the original timer, and deletion removes the sandbox after verifying the saved copy. | Pending | Pending |

## Local-to-cloud handoff

Guide: [Continue local work in the cloud](../../content/manuals/agentic-platform/guides/local-handoff.md).

| Guide section and claim to test | Expected result | Status | Date, observations, and evidence |
| --- | --- | --- | --- |
| Introduction and prerequisites: provisional 25-minute duration, CLI/runtime, and architecture | Measure the journey and replace provisional timing. Record evidence for minimum CLI version, supported local runtime, and support for both amd64 and arm64 in the cloud. The destination matches the local architecture. | Pending | Pending |
| Before you begin: login, diagnose, and separate secret commands | Commands succeed with the recorded account. Cloud Claude Code requires an API key, while local Claude Code also supports subscriptions. Verify this distinction against Anthropic's terms and supported authentication guidance. Local and cloud Anthropic credentials are configured separately; record required scope without keys. Confirm compute and inference charging guidance against billing information. | Pending | Pending |
| Step 1: create, policy, clone, and run commands | Sandbox is created without a host workspace argument; public sample is cloned inside its filesystem. Agent starts with the managed local credential and no GitHub token. | Pending | Pending |
| Step 1: bounded prompt, baseline, and credential removal | HANDOFF.md records the full matching Git revision, clean tree, project path, and remaining task without starting the review. No application files change or dependencies install. Review filesystem credential files before capture and remove them without recording secrets in evidence. | Pending | Pending |
| Step 2: move command and timeout support | Exact move command succeeds with a two-hour Stop action or fails explicitly if stopping is unsupported. Local source stops for capture; a distinct cloud destination and backing template ID are reported. The template contains the captured HANDOFF.md. Record command output and supported options. | Pending | Pending |
| Step 2: transfer boundaries | Destination contains filesystem files, but local processes, memory, sockets, host mounts, clone-mode volumes, managed secrets, and local network rules aren't transferred. Use disposable non-sensitive fixtures for applicable boundaries and record any boundary not exercised as pending. Verify documented TCP port republication/warning behavior separately on disposable resources. | Pending | Pending |
| Step 3: transferred handoff and destination settings | The exec cat command succeeds using the printed sbx_-prefixed ID. The handoff file contains the recorded baseline, project path, and remaining task. Platform, size, expiration, cloud credential, and destination policy are correct; no unintended ports remain. | Pending | Pending |
| Step 3: run and restart the agent | Run starts the cloud agent using destination credentials and policy. Agent confirms the baseline and completes only the saved task, without an application run. The local agent process hasn't transferred. | Pending | Pending |
| Step 4: detach, continue, and run after stop | Ctrl+\ detaches without ending cloud work. Agent continues while the client is disconnected until timeout. No work runs while stopped; run resumes the stopped destination and starts a fresh timer. | Pending | Pending |
| Step 4: retrieve and verify artifacts | Both cp commands succeed and local copies contain the matching baseline, build/startup explanation, listening port, one recommendation with source references, and limitations. Compare the report with the source; record corrections rather than accepting agent claims as test evidence. | Pending | Pending |
| Step 5: remove both sandboxes and the backing template, or retain a stopped destination | Sandbox removal succeeds before template removal; local/cloud sandbox lists and the cloud template list confirm absence after local copies are verified. Confirm the template persists if only its sandbox is removed. If retaining the destination, confirm stopped status, expiration after resume, file retention, and billing guidance. Remove all disposable test resources. | Pending | Pending |

## Publication decision

- [ ] Every required row passes with dated observations and evidence.
- [ ] Both editor journeys pass with separate version records.
- [ ] Commands, UI labels, availability, and minimum versions match the evidence.
- [ ] Provisional durations are replaced with measured durations.
- [ ] Discrepancies are corrected in the guides and affected checks repeated.
- [ ] Saved copies are verified and test sandboxes, backing templates, public ports, and temporary credentials are cleaned up.
- [ ] Reviewer records the tested docs commit and approves publication.

Decision: Pending. Reviewer, date, and remaining blockers: Pending.
