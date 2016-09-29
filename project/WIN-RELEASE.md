# Preparing a Windows release

This document describes the Windows-specific part of the release process. See
[RELEASE.md](RELEASE.md) for the overall process.

On windows, pushing a release on own s3 buckets is an automatic process.

## Tags and Channels

All of the process is made on AppVeyor. Each time it builds either on `master`
branch or a tag, AppVeyor will build an installer an upload it to s3 using the `docker-release` command.

Before tagging a release, please make sure the `CHANGELOG` file is up to date.

**WARNING**: To make sure we don't pollute OSX releases, every Windows release tag name must
start with `win-` prefix.

The build script analyze the tag name, and determine on which channel it uploads the artifact.
the tag format is the following `win-[channel]-v[version]-[milestone]` for example `win-beta-v1.12.0-beta18`.
`[channel]` can be either `beta`, `test`. Releasing on `stable` is made by omitting the channel altogether for
instance `win-v1.12.0-beta19`.

The channel it uploads to is chosen by those rules:
- If the commit has a tag that starts with `win-v`, then the installer is pushed to `stable` channel.
- If the commit has a tag that starts with `win-beta-v`, then the installer is pushed to `beta` channel.
- If the commit has a tag that starts with `win-test-v`, then the installer is pushed to `test` channel.
- Otherwise (either on master or malformed tag), the installer is pushed to `master` channel.

Tagging example
```
git tag win-test-v1.11-beta11
git push origin win-test-v1.11-beta11
```

Release to the `beta` or `stable` channels are not published right away to users, they are uploaded
and available to download but it will not trigger an automated users. You must use the `docker-release`
command to achieve this.

## Preparing the next iteration

After a release, the version needs to be bumped for the next release:

 - Update `src/SolutionInfo.cs` file
 - Set the date on the `CHANGELOG` for the current release, replacing `unreleased` by yyyy-mm-dd
 - Create a section for this release in the `CHANGELOG` file
 - Commit theses changes with a commit message similar to `Bump to version X.Y.Z-dev"
