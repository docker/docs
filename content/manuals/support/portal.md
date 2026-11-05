---
title: Docker Support portal
linkTitle: Support portal
description: >-
  Open the Docker Support portal, create and manage support cases, and see
  who can view or update a case.
keywords: support portal, support cases, create case, case history, resolve
  case, watchers, attachments, severity levels, case status
weight: 10
---

The Docker Support portal lets you create and manage support cases, follow
their progress, and collaborate with Docker Support and your teammates in
one place. This page covers how to open the portal, how to work with
cases, and who can see or update a case.

## Who gets access

> [!NOTE]
>
> The Support portal is being rolled out in phases, so it may not yet be
> available to every organization. Organizations with a Premium Support
> subscription are enrolled automatically.

After your organization is enrolled, open the portal from the **Help**
menu.

## Access the Support portal

> [!TIP]
>
> Before reaching out for support, review the troubleshooting documentation
> for your product.

1. Sign in to [Docker Hub](https://hub.docker.com/) or
   [Docker Home](https://app.docker.com/). You can also open the portal
   from [Docker Scout](https://scout.docker.com/).
1. Select the organization from the organization switcher. When the case
   belongs to a personal account, select that account instead.
1. Open the **Help** (?) menu in the top navigation.
1. Select **Support portal**.

The Support portal opens in your current browser tab.

On Docker Home, **Support portal** is also in the account sidebar. **Help**
shows **Contact support** for accounts that do not have the portal. That
item opens the [support request form](https://hub.docker.com/support/contact/).

## Landing page

The landing page gives you a snapshot of your support activity. From here
you can:

- See your active cases
- Select **Create case** to create a new case
- Select a case, or **View details**, to open that case
- Select **View cases** to see your full case history
- Open the **⋮** menu on a case and select **Details** or **Resolve case**

The page lists these cases under **Active cases**. If there are no open
cases, the page shows **No active cases**. Select **See case history** to
open past cases.

## Working with cases

### Create a case

1. Select **Create case**.
1. Enter a **Subject** and **Details**.
1. Select **Get instant help**. Gordon, Docker's AI assistant, reads your
   subject and details and posts an answer below the form. If the answer
   solves the issue, you can stop there.
1. To send the case to Docker Support, complete the remaining fields:
   - Confirm **Contact name**. **Email** is filled in from your account
     and is read-only.
   - Select the **Namespace** the issue relates to. This field is
     required. A namespace is the Docker organization or personal account
     the case is for.
   - Select a **Topic**. If the topic has subtopics, select a
     **Subtopic**.
   - If a **Severity** field appears, select the severity that matches
     your issue, using the [severity levels](#severity-levels) below. The
     field does not appear for the Account, Billing, Legal, Sales, and
     Security topics, because a severity is set automatically for them.
   - Complete any additional fields the topic shows. These can include
     **Version number**, browser, **Operating system**, **Diagnostics ID**,
     and **Affected users**. Docker Desktop asks for all of these.
1. Optionally, attach files to the case. Select **Select files to attach**
   or drag and drop files into the upload area. See
   [Attach files](#attach-files) for file rules.
1. Select **Submit**.

The portal returns to the landing page and lists the case under **Active
cases**. The case number shows as **PENDING CASE NUMBER…** until Docker
Support assigns one.

#### Severity levels

| Severity | Definition |
| :--- | :--- |
| S1 (Critical) | Complete loss of core functionality; broad/systemic impact; no workaround. |
| S2 (High) | Significant subset of users loses key functionality; severe business impact; no workaround. |
| S3 (Medium) | Partial/non-critical loss for an individual or small group; workaround typically available. |
| S4 (Low) | Inquiry, minor issue, or enhancement request. |

### View cases

The **Case history** page lists every case associated with your
organization. Each row shows the status, case number, subject, topic,
subtopic, who created it, and when it was last modified. A company
account also shows the organization.

Select a case to open it. You can show 10, 25, or 50 cases on a page.
**Create case** and the **⋮** menu work the same way as on the landing
page.

#### Case statuses

| Status | Meaning |
| :--- | :--- |
| NEW | Your case was received and is waiting for Docker Support to respond. |
| PENDING SUPPORT | Docker Support is working on your case. |
| PENDING CUSTOMER | Docker Support is waiting for a response from you. |
| ON HOLD | Docker Support has replied, and the resolution depends on another team at Docker. |
| CLOSED | The case is complete. |
| RESOLVED | The case was marked as resolved, by you or by Docker Support. Docker Support will review and close it. |

### Case details

Each case shows its full progress and history: the case number, subject,
topic, summary, description, and the comments and responses in
**Activity**.

What you can do depends on your access. See
[Who can interact with a case](#who-can-interact-with-a-case). The steps
below cover comments, watchers, files, and resolving the case.

#### Comment on a case

1. Open the case.
1. Enter your message in the **Activity** editor.
1. Optionally, attach files. Select **Select files to attach** or drag
   and drop files into the upload area. See
   [Attach files](#attach-files) for file rules.
1. Select **Send**.

You can comment on a case only if you created it or were added as a
watcher. If you can see a case but are not allowed to comment, the comment
box is unavailable and explains why. The message tells you to contact the
case creator.

### Add watchers

A watcher is a Docker user who follows a case and can take part in it
without being the person who created it. Watchers are available on cases
created under an organization or a company.

#### Add a watcher

1. Open the case.
1. In **Enter additional watchers**, choose an organization member from
   the list, or type the email address of a Docker user.
1. Select **Save**.

A **Watchers updated successfully** message confirms the change, and new
watchers receive an email letting them know they were added. The email can
take a few minutes to arrive.

#### Remove a watcher

1. Select the **x** next to the watcher.
1. Select **Save**.

Before saving, select **Undo** to restore a watcher you removed.
Removed watchers are notified by email and can no longer comment on the
case, upload files, or download attachments.

#### Watcher rules

- Watchers must have a Docker account. You can add Docker users who are
  not members of your organization.
- Each watcher can be added only once.
- The case creator cannot be added as a watcher.
- Watchers can comment, upload and download attachments, and add or remove
  other watchers.
- Watchers receive email notifications about new comments and attachments
  on the case.
- The watcher field is not available on cases created under a personal
  account.

### Attach files

You can attach files when you create a case or when you
[comment on a case](#comment-on-a-case).

Uploaded files appear in the **Attachments** section at the top of the
case. A file becomes a download link after it passes the security check.
You can download the file if you can comment on the case.

#### File requirements

- Each file can be up to 50 MB. Empty files are not accepted.
- Supported file types: `tar`, `gz`, `tgz`, `tar.gz`, `zip`, `rar`,
  `jpeg`, `jpg`, `png`, `json`, `txt`, `log`, `err`, `xml`, `yaml`,
  `yml`, `plist`, `evtx`, `evt`
- File names cannot contain path characters, end with a dot or space, or
  exceed 255 bytes.
- If any file you select is not valid, an error explains why and none of
  the files are uploaded.

#### Security checks

Every file you upload is automatically scanned for malware, and its
contents are checked to confirm they match the file type. After you send
a file, it shows as **pending** while it is checked. If it fails either
check, the file is removed. It appears crossed out in **Attachments** with
a **security risk detected** message and cannot be downloaded. The message
does not say which check failed. If you believe a file was flagged by
mistake, add a comment to your case.

**Expired** means the upload did not finish in time and the file is not
available. Attach the file again.

#### Upload limits

Uploads are limited to keep the portal reliable.

| Limit | Value |
| :--- | :--- |
| Maximum size per file | 50 MB |
| Uploads per hour, per user | 10 files |
| Total upload size per day, per user | 5 GiB |
| Total attachment size per case | 2 GiB |

Hourly and daily limits reset at the start of each hour and each day, in
UTC. The per-case limit applies to all attachments on the case combined.
If you reach a limit, wait until it resets and try again.

Files cannot be added to closed or resolved cases. Files already attached
to a case stay available to download for 90 days after the case is closed.

### Resolve a case

If your issue is solved, you can resolve the case yourself without waiting
for Docker Support to close it.

1. Select **Resolve case** on the case, or **Resolve case** from the
   **⋮** menu on the landing page or **Case history**.
1. Select a reason:
   - **Issue resolved**
   - **No longer relevant**
   - **Duplicate case**
   - **Resolved by another team**
   - **Other**
1. Optionally, add comments of up to 1000 characters.
1. Select **Resolve case**.

The case status changes to **RESOLVED**, and Docker Support is notified.
When Docker Support closes the case, its status changes to **CLOSED** and
the person who created the case receives a short satisfaction survey by
email.

#### Resolved and closed cases

- Resolved cases cannot be reopened.
- You can still read the full activity history, but the comment box,
  watcher field, and file upload are hidden. The page says the case is
  closed.
- If you need more help, select **Create case** to start a new case.

## Permissions

### Who can see a case

Access to cases is based on organization membership.

- If you are part of an organization, you can see that organization's
  cases.
- Organization owners can see all cases in their organization.
- Company owners can see all cases across their organizations.
- The case creator can see the case.
- Watchers can see the cases they were added to.

Access flows downward only, from company to organization to user. Nothing
flows upward automatically.

### Who can interact with a case

Seeing a case and interacting with it are separate permissions. You can
interact with a case if you:

- Created the case
- Were added as a watcher on the case

If you are an organization member who can see a case but are not the
creator or a watcher, you can read the case but cannot comment on it. The
comment box is unavailable and explains why. Ask the case creator to add
you as a watcher.

### Teams and case access

[Teams](/manuals/accounts/organization/manage/manage-a-team.md) are groups
of Docker users within an organization, used to manage repository access.
Case visibility and permissions are controlled at the organization level.

### Summary of case permissions

| Action | Who can do it |
| :--- | :--- |
| See a case | Organization members, organization owners, company owners, case creator, watchers |
| Comment | Case creator, watchers |
| Add or remove watchers | Case creator, watchers |
| Upload attachments | Case creator, watchers |
| Download attachments | Case creator, watchers |
| Resolve a case | Case creator, watchers |
