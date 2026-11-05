---
title: Create and manage Docker Support portal tickets
linkTitle: Create support ticket
description: >-
  Create a Docker support case, follow its status, add watchers, upload
  files, and resolve cases in the Docker Support portal.
keywords: support portal, support cases, create case, case history, resolve
  case, watchers, attachments
weight: 10
---

The Docker Support portal is where paid subscribers create a case, follow
its status, add comments and watchers, and upload files.

## Open the portal

> [!TIP]
>
> Before reaching out for support, review the troubleshooting documentation
> for your product.

Open **Help** in [Docker Hub](https://hub.docker.com/),
[Docker Home](https://app.docker.com/), or
[Docker Scout](https://scout.docker.com/):

- Select **Support portal** to open the portal for the account you're using.
- Select **Contact support** to open the
  [support request form](https://hub.docker.com/support/contact/). **Help**
  shows this item for accounts that don't have the portal.

On Docker Home, **Support portal** is also in the account sidebar.

To open the portal:

1. Sign in, and select the organization or personal account the case belongs
   to.
1. Open **Help** and select **Support portal**.

## Portal home

The portal home lists **Active cases** and includes **View cases** and
**Create case**.

- Select a case or **View details** to open it.
- Open the **⋮** menu on a case and select **Details** or **Resolve case**.
- Select **View cases** to open the full case history.
- If there are no open cases, the page shows **No active cases**. Select
  **See case history** to open past cases.

## Create a case

1. Select **Create case**.
1. Enter a **Subject** and **Details** that describe the issue.
1. Select **Get instant help**. Gordon, Docker's AI assistant, reads your
   subject and details and posts an answer below the form. If the answer
   solves the issue, you can stop there.
1. To send the case to the Support team, complete the remaining fields:
   - Confirm **Contact name**. **Email** is filled in from your account
     and can't be edited.
   - Select a **Namespace**. This is the Docker organization or personal
     account the case is for.
   - Select a **Topic**. If the topic has subtopics, select a **Subtopic**.
   - Fill in any other fields the topic shows, such as **Version number**,
     **Operating system**, **Diagnostics ID**, or **Affected users**. For
     Docker Desktop, the form asks for all of these.
   - If **Severity** is shown, select the level that matches the issue. See
     [Support severity levels](/manuals/support/_index.md#support-severity-levels).
1. To add files, select **Select files to attach** or drop files on the
   upload area. See [Attach files](#attach-files) for file rules.
1. Select **Submit**.

The portal returns to the home page and lists the case under **Active
cases**. The case number shows as **PENDING CASE NUMBER…** until Docker
Support assigns one.

## Attach files

You can attach files when you create a case or when you comment on an
open case. Files are listed under **Attachments**. A file shows as
**pending** while it is checked. After the check:

- The filename becomes a download link if you can comment on the case.
- **security risk detected** means the file was rejected. The name is
  crossed out and can't be downloaded. Comment on the case if you think
  the file was rejected by mistake.
- **expired** means the file is no longer available to download.

Every uploaded file is scanned for malware and checked to confirm its
contents match the file type. A file that fails either check shows
**security risk detected** without identifying which check failed.

File rules:

- Each file can be up to 50 MB. Empty files are rejected.
- Accepted types are `tar`, `gz`, `tgz`, `tar.gz`, `zip`, `rar`, `jpeg`,
  `jpg`, `png`, `json`, `txt`, `log`, `err`, `xml`, `yaml`, `yml`,
  `plist`, `evtx`, and `evt`.
- Filenames can't contain path characters, end with a period or space,
  or exceed 255 bytes.
- If any file you select is invalid, an error explains why and none of
  the files are uploaded.

Upload limits:

| Limit | Value |
| :--- | :--- |
| Maximum size per file | 50 MB |
| Uploads per hour, per user | 10 files |
| Total upload size per day, per user | 5 GiB |
| Total attachment size per case | 2 GiB |

Hourly and daily limits reset at the start of each hour and each day
(UTC). The per-case limit applies to all attachments on the case
combined. If you reach a limit, wait until it resets and try again.

You can add files only while the case is open. Files on a resolved or
closed case remain available to download.

## Manage case

After a case is open, you can review its history, comment, add watchers,
attach files, and resolve it.

### Case history

**Case history** lists cases for the account. Each row shows the status,
case number, subject, topic, subtopic, who created the case, and when it
was last modified. A company account also shows the organization.

From **Case history** you can:

- Select a case to open it.
- Select **Create case** to start another case.
- Use the **⋮** menu to open **Details** or **Resolve case**.
- Show 10, 25, or 50 cases on a page.

### Case detail view

A case shows the case number, subject, topic, description, and an
**Activity** panel with comments. On an open case you can comment, add or
remove watchers, attach files, and resolve the case. What you can do
depends on your access. See [Who can work on a case](#who-can-work-on-a-case).

### Comment on a case

1. Open the case.
1. Under **Activity**, do one or more:
   - Enter a comment.
   - To attach files, select **Select files to attach** or drop files on
     the upload area. See [Attach files](#attach-files) for file rules.
1. Select **Send**.

You can comment when you created the case or when someone added you as a
watcher. Otherwise the field tells you to contact the case creator.

### Add or remove watchers

A watcher follows an organization or company case and can comment, attach
files, and add or remove other watchers. The watcher field appears on
organization and company cases.

To add or remove a watcher:

1. Open the case.
1. In **Enter additional watchers**, do one of the following:
   - To add a watcher, select an organization member or enter the email
     address of a Docker user. You can add Docker users who aren't
     members of your organization.
   - To remove a watcher, select the **x** on that person. Select **Undo**
     to restore them before you save.
1. Select **Save**.

The page confirms that the watchers were updated, and Docker emails the
people you added or removed. The case creator can't be added as a watcher,
and each person can be added only once. A removed watcher can no longer
comment on the case or upload or download attachments.

Watchers also get an email when the case has new activity.

### Resolve a case

If the issue is solved, resolve the case yourself.

1. Start resolution from either of these places:
   - Select **Resolve case** on the open case.
   - Select **Resolve case** from the **⋮** menu on the portal home or on
     **Case history**.
1. Select a reason:
   - Issue resolved
   - No longer relevant
   - Duplicate case
   - Resolved by another team
   - Other
1. Add comments if you have them. Comments can be up to 1000 characters.
1. Select **Resolve case**.

The status changes to **RESOLVED**. A resolved case can't be reopened.
When Docker Support closes the case, the status changes to **CLOSED** and
the person who created the case receives a short satisfaction survey by
email.

On a resolved or closed case, you can still read **Activity**. The
comment field, watcher field, and file upload are hidden, and the page
says the case is closed. Select **Create case** to start another case.

## Who can work on a case

Case access follows organization membership. Access flows downward only,
from company to organization to user; nothing flows upward.
[Teams](/manuals/accounts/organization/manage/manage-a-team/) don't
change who can see or update a case.

- Organization members can see that organization's cases.
- Organization owners can see cases in their organization.
- Company owners can see cases for the organizations in the company.
- The case creator and any watchers can see the case.

| Action | Who can do it |
| :--- | :--- |
| See a case | Organization members, organization owners, company owners, the case creator, and watchers |
| Comment, add or remove watchers, upload files, download files | Case creator and watchers |
| Resolve a case | Case creator and watchers |

When the comment field tells you to contact the case creator, ask them to
add you as a watcher.