Triaging of issues
------------------

Triage provides an important way to contribute to an open source
project. Triage helps ensure issues resolve quickly by:

- Describing the issue's intent and purpose is conveyed
  precisely. This is necessary because it can be difficult for an
  issue to explain how an end user experiences a problem and what
  actions they took.
- Giving a contributor the information they need before they commit to
  resolving an issue.
- Lowering the issue count by preventing duplicate issues.
- Streamlining the development process by preventing duplicate discussions.

If you don't have time to code, consider helping with triage. The
community will thank you for saving them time by spending some of
yours.

### 1. Ensure the issue contains basic information

Before triaging an issue very far, make sure that the issue's author
provided the standard issue information. This will help you make an
educated recommendation on how this to categorize the issue. Standard
information that *must* be included in most issues are things such as:

-   the output of:
   -  `pinata diagnose -u` on OSX
   -  `DockerDebugInfo.ps1` using Powershell on Windows
-   a reproducible case if this is a bug, Dockerfiles FTW
-   page URL if this is a docs issue or the name of a man page

Depending on the issue, you might not feel all this information is
needed. Use your best judgement. If you cannot triage an issue using
what its author provided, explain kindly to the author that they must
provide the above information to clarify the problem.

If the author provides the standard information but you are still
unable to triage the issue, request additional information. Do this
kindly and politely because you are asking for more of the author's
time.

If the author does not respond requested information within the
timespan of a week, close the issue with a kind note stating that the
author can request for the issue to be reopened when the necessary
information is provided.

### 2. Classify the Issue

An issue can have multiple of the following labels.

#### Issue kind

| Kind             | Description                                                                                                                     |
|------------------|---------------------------------------------------------------------------------------------------------------------------------|
| kind/bug         | Bugs are bugs. The cause may or may not be known at triage time so debugging should be taken account into the time estimate.    |
| kind/docs        | Writing documentation, man pages, articles, blogs, or other significant word-driven task.                                       |
| kind/enhancement | Enhancement are not bugs or new features but can drastically improve usability or performance of a project component.           |
| kind/feature     | Functionality or other elements that the project does not currently support.  Features are new and shiny.                      |
| kind/performance | Performance-related issues                                                                                                      |
| kind/question    | Contains a user or contributor question requiring a response.                                                                   |
| kind/roadmap | Issue to track large items that need to be delivered. See [roadmap docs](ROADMAP.md) for more details. |
| kind/tracking    | High-level issue to keep track of other issues. Tracking issues can be closed when all child issues are resolved.               |
| kind/tracking/users | High-level issue to keep track of user feedback.                                                                             |

#### Functional area

| Area Common               | Description                                                                                                            |
|---------------------------|------------------------------------------------------------------------------------------------------------------------|
| area/qa                   | Issue related to code quality and infrasture-related tasks to deal with automated QA                                   |
| area/logging              | Ensure consistent logging accross all the components                                                                   |
| area/support              | The built-in features for supporting users and reporting defects to the engineering team                               |
| area/tokens               | Token mechanism to control the growth of beta users                                                                    |
| area/sdk                  |                                                                                                                        |
| area/agent                |                                                                                                                        |
| area/proxy                |                                                                                                                        |
| area/upstream             | Upstream bugs which needs to be reported                                                                               |
| area/tools                | Issues about the integration of open-source Docker tools (machine, toolbox, notary, etc.) in Pinata                    |

| Area Windows              | Description                                                                                                            |
|---------------------------|------------------------------------------------------------------------------------------------------------------------|
| area/windows              |                                                                                                                        |
| area/windows/build        |                                                                                                                        |
| area/windows/installer    |                                                                                                                        |
| area/windows/menubar      |                                                                                                                        |
| area/windows/cli          |                                                                                                                        |

| Area OSX                  | Description                                                                                                            |
|---------------------------|------------------------------------------------------------------------------------------------------------------------|
| area/osx                  |                                                                                                                        |
| area/osx/build            |                                                                                                                        |
| area/osx/menubar          |                                                                                                                        |
| area/osx/prefpane         |                                                                                                                        |
| area/osx/installer        |                                                                                                                        |
| area/osx/watchdog         |                                                                                                                        |
| area/osx/cli              |                                                                                                                        |

| Area Backend                 | Description                                                                                                         |
|------------------------------|---------------------------------------------------------------------------------------------------------------------|
| area/backend                 | General issues about the backend. If possible backend issues should have a more specific labels than this one.      |
| area/backend/build           | Build issues related to the backend components (OPAM packages, etc.)                                                |
| area/backend/storage/volumes | Host/container file-sharing issues                                                                                  |
| area/backend/network/nat     | NAT networking issues                                                                                               |
| area/backend/network/hostnet | Hostnet networking issues                                                                                           |
| area/backend/network/mdns    | mDNS issues                                                                                                         |
| area/backend/run/qemu        | Issues related to the QEMU execution driver                                                                         |
| area/backend/run/xhyve       | Issues related with the Xhyve hypervisor                                                                            |
| area/backend/run/hyper-v     | Issues related with the Hype-V hypervisor                                                                           |
| area/backend/moby/linux      | Issues with the OS distribution for running Linux containers                                                        |
| area/backend/moby/windows    | Issues with the OS distribution for running Windows containers                                                      |

### 3. Prioritizing issues

Issues with a known resolution date should be given a priority and attached to a milestone.
Issues that do not have a known resolution date should have a priority attached as this will assist with planning.

The following labels are used to indicate the degree of priority (from more urgent to less urgent).

| Priority    | Description                                                                                                                       |
|-------------|-----------------------------------------------------------------------------------------------------------------------------------|
| priority/P0 | Urgent: Security, critical bugs, blocking issues. P0 basically means drop everything you are doing until this issue is addressed. |
| prority/unknown | The resolution date and/or the priority of that issue unclear. That issue should be discussed in a weekly meeting.            |
| priority/P1 | Important: P1 issues are a top priority and a must-have for the next release.                                                     |
| priority/P2 | Normal priority: default priority applied.                                                                                        |
| priority/P3 | Best effort: those are nice to have / minor issues.                                                                               |

The release manager (or the maintainer if they know what they are
doing) is responsible for prioritizing issues.

### 4. Handling User Reported Issues

Issues that have been reported by a user via email to beta-feedback@docker.com, on the Docker Forums or by other means MUST have the `kind/tracking/user` label attached. 

A user reported issue is prioritized in accordance with the following table:

| Priority | Descrtiption                                                                                                      |
|----------|-------------------------------------------------------------------------------------------------------------------|
| P0       | As P1 but effects a large portion of the user base. Please consult the release manager before applying this label |
| P1       | I can't install/use the app at all - it's totally broken!                                                         |
| P2       | The app doesn't work and it severely impacts my workflow                                                          |
| P3       | I have a problem that effects my workflow but it doesn't stop me getting work done                                |

If multiple reports of the same issue are provided, this will enable the release manager or PM to bump the priority of the issue.
Links to the source of the reports should be added to the GitHub issue.

### 4. Milestone

There should be a milestone per planned release. The release managers are
reponsible for assigning issues to milestones, in coordination with PM and
the rest of the engineering group.

### 5. Status

Status labels can be attached to issues to add more detailed status than just "open" and "closed". These are particularly useful to track the status of roadmap issues.

| Label                     | Description                                                                                                            |
|---------------------------|------------------------------------------------------------------------------------------------------------------------|
| status/at-risk            | Issue may not land in the milestone specified                                                                          |
| status/delayed            | Issue will not land in the milestone specified                                                                         |
| status/in-progress        | Issue is currently being worked on                                                                                     |
| status/postponed          | Issue is being pushed back and milestone is yet to be determined                                                       |
