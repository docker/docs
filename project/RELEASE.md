# Project Pinata release process

The Pinata release process currently takes two days to allow the teams in each time zone to prepare for the release. The first day is for planning and merging final PRs and freezing the UI. The second day is for creating and testing final binaries before the actual release in the afternoon/evening. If the release is to be synchronised with a documentation release, the documentation team also needs time to prepare final changes after the binaries have been built.

This document describes the release process as of July 5th 2016.

Jump to

 - [Summary: Release preparations day 1 ](#release-preparations-day-1)
  - [Issue triaging](#issue-triaging)
  - [Release meeting](#release-meeting)
  - [Github tracking issue](#github-tracking-issues)
  - [Updating CHANGELOGs](#changelogs)
 - [Summary: Release preparations day 2](#release-preparations-day-2)
  - [Tagging releases](#tagging-releases)
  - [Synchronised release](#synchronised-release)
 - [Troubleshooting](#troubleshooting)
  - [Web UI](#web-ui)
  - [Rollback Release](#rollback-release)
  - [Update release notes with docker-release](#update-release-notes-with-docker-release)
  - [Surgically updating the CHANGELOG](#surgically-updating-the-changelog)
 - [Hotfixes and updates on stable channel](#hotfixes-and-updates-on-stable-channel)
  - [Processes for PRs](#processes-for-prs)
  - [Releasing a beta hotfix](#releasing-a-beta-hotfix)
  - [Releasing a stable hotfix](#releasing-a-stable-hotfix)

## Release preparations day 1

Summary:

- Triage issues before release meeting ([see details below](#issue-triaging))
- Prepare tracking issues on Github ([see details below](#github-tracking-issues))
- Publish release plans in #pinata-dev which should contain:
    - When next release is scheduled
    - Link to agenda for the [release meeting](#release-meeting)
    - Link to Github tracking issue
    - Links to open Github issues for this milestone and remaining P0 issues. If there are many open issues it may also be useful with a reminder to triage them, as they will all be discussed in the release meeting.
- Prepare CHANGELOGs for Windows/Mac ([see details below](#changelogs))
- Release meeting @ 5pm UK time ([see details below](#release-meeting))
- When UI complete binaries are ready, send binary links and initial CHANGELOGs to documentation team

#### Github tracking issues

Each release is tracked by a tracking issue on Github. The tracking issue should include the tasks that should be completed per day and the person assigned to it, if relevant. The issue should be tagged `kind/tracking/release`, `area/osx` and `area/windows`.

[List of earlier tracking issues](https://github.com/docker/pinata/issues?utf8=✓&q=is%3Aissue%20label%3Akind%2Ftracking%2Frelease%20)

In addition, create a tracking issues for the test plans and link them from the main tracking issue.

The Windows testing issue template is in `WIN-TESTING.md`, and the new issue should be called: `win: Release testing checklist` and be tagged with `kind/tracking/release` and `area/windows`.

The Mac testing issue template is in `MAC-TESTING.md`, and the new issue should be called: `mac: Release testing checklist` and be tagged with `kind/tracking/release` and `area/osx`.

#### Issue triaging

Triage bugs on Github with the DRIs using the latest milestones defined on
https://github.com/docker/pinata/milestones

See ISSUE-TRIAGE.md for details about the triaging process.

#### CHANGELOGs

Update the CHANGELOG files for Mac and Windows, with user-facing updates since the last release.
  - add issue/PR number if available
  - do not put names
  - the "compare" feature of GitHub is helpful:
  `https://github.com/docker/pinata/compare/<last-VERSION>...master`

An initial CHANGELOG should be prepared prior to the release meeting so it can be discussed in the meeting and be ready for the documentation team. The final CHANGELOG will not be ready until actual release. Sync with docs team if there are large changes from the initial version.

Remember to include docs in the list of changes if the release is synchronised with the binaries.

The Windows and Mac CHANGELOGs are currently stored separately in [CHANGELOG](../CHANGELOG) and [win/CHANGELOG](../win/CHANGELOG).

#### Release meeting

The release meeting is held by the release manager on Mondays @ 5pm UK time. The agenda is posted in `#pinata-dev` Monday morning.

Bluejeans link: https://bluejeans.com/747188286

Agenda should include
  - Summary of previous release, discuss possible improvements
  - Plans for next release (expected release time, release plans for documentation, engine upgrades etc)
  - Discuss any open P0s on Github, plans for fixing them and/or if release should be delayed
  - Triage open issues for current release milestone, move issues that will not be fixed in time
  - Discuss UI changes for the release and if there are new features/UI that require documentation updates
  - Go through [CHANGELOGs](#changelogs) for Windows/Mac
  - If time: go through issues without an assigned milestones, in case there are important new issues

Some earlier agendas
- [Beta 17 agenda](https://docs.google.com/document/d/1KdyKZ4Jvl4e42QSa4h42SZZaKvYmFKhXGzt1ApvvojM/edit)
- [Beta 13 agenda with notes](https://docs.google.com/document/d/1YPgm-oH_E1A4DylIlBXPz49WQ1hLVJjMEj3LmJkmRlU/edit)

## Release preparations day 2

Summary:
- If all PRs are in (check with DRIs, @pinata-team on slack) tag a test release ([see details below](#release-candidate))
- Follow release process per platform: [Windows](WIN-RELEASE.md) / [Mac](MAC-RELEASE.md)
- Test according to test plans: [Windows](WIN-TESTING.md) / [Mac](MAC-TESTING.md)). As a minimum the release should pass all `rt-local` tests with the `release` label.
- Tag release binaries for beta or stable channel ([see details below](#tagging-releases))
- When release binaries are ready, run final tests on them and post binary link in `#pinata-dev`
- Send release binaries to @victoria with final release notes
- If synchronised release with docs, wait for final documentation updates
- Release binaries and announce the new release ([see below for details](#synchronise-release))

### Tagging releases

#### Release candidate

Creating a test release of Beta20 with Docker 1.12.0-rc4:

* for Mac the tag should be: `mac-rc-v1.12.0-beta20` (note that for engine versions with a suffix, such as `1.12.0-rc4`, the tag is still `mac-rc-v1.12.0-beta20` due to limitations with current regexp)
* for Windows the tag should be: `win-test-v1.12.0-rc4-beta15`

This will upload the artifact to the following endpoint, given a CI build of `8000`, and auto-publish it:

```
https://download-stage.docker.com/mac/rc/1.12.0.8000/
https://download-stage.docker.com/win/test/1.12.0.8000/
```

E.g.:
Make sure on right commit/branch, then:

```
VERSION="mac-rc-v1.12.0-beta20"
git tag $(VERSION) -a -m "Test release $(VERSION)"
git push upstream $(VERSION)
```

After the tag has been built the files on each endpoint should be:
```
NOTES
Docker.dmg.sha256sum (mac)
Docker.dmg (mac)
InstallDocker.msi (win)
InstallDocker.msi.sha256sum (win)
```

For more details about the uploaded artifacts, see the CI build output.

##### Beta Release

Creating a beta release of Beta20 with Docker 1.12.0-rc4:

* for Mac the tag should be: `mac-v1.12.0-beta20`
* for Windows the tag should be: `win-beta-v1.12.0-rc4-beta20`

(Note that only the engine version is used in the Mac tag with suffixes such as `-rcX` ignored. E.g. the tag for `v1.12.0-rc4` is `mac-v1.12.0-beta20`. This is due to a limitation of the regexp used in the Mac build scripts. This is not the case for the Windows tag.)

This will upload the artifact to the following endpoint, given a CI build of `8001`:
```
https://download.docker.com/mac/beta/1.12.0.8001/
https://download.docker.com/win/beta/1.12.0.8001/
```

The binaries are not released until the `appcast.xml` file is updated - how to release them is described [below](#release-artifacts)

##### Stable Release

Creating a stable release with Docker 1.12.0:

* for Mac the tag should be: `mac-v1.12.0`
* for Windows the tag should be: `win-v1.12.0`

This will upload the artifact to the following endpoint, given a CI build of `8002`, and hold-it to be published manually:
```
https://download.docker.com/mac/stable/1.12.0.8002/
https://download.docker.com/win/stable/1.12.0.8002/
```

The stable channel is not used until after 1.12 GA / Pinata GA.

### Synchronised release

When the builds have been tested and docs are ready for deployment, add Github release details, release the artifacts with `docker-release` and announce the release.

Verify that the release notes are correct prior to release, as these will be shown in the Auto update window.

Check that docs are up to date (verify release notes!), S3 downloads are working and correct and that
License.tar.gz has been updated.

#### Create GitHub releases

Create GitHub releases for Win and Mac and update the
release notes in https://github.com/docker/pinata/releases:

- create a new release for the value of the Mac release tag
- copy/paste the related lines from the the `CHANGELOG` file
- create a new release for the value of the Win release tag
- copy/paste the related lines from the the `win/CHANGELOG` file

#### Release artifacts

To release the artifacts to the public from the beta channel, given a version of `1.12.0.8002` run:

OSX:

`docker-release --aws-access-key xxx --aws-secret-key xxx --channel beta --build 1.12.0.8002 --human 1.12.0-beta20 --arch mac -prod publish`

Windows:

`docker-release --aws-access-key xxx --aws-secret-key xxx --channel beta --build 1.12.0.8002 --human 1.12.0-beta20 --arch win -prod publish`

Replace the AWS keys, the build number and the human string with the correct values.

All other channels than stable and beta, are published automatically upon upload of the artifacts. Publishing a certain version also promotes that version as the 'latest' downloadable artifact, which will create the following endpoint for DOCS user download purpose (new user trying to download the edition)

OSX:

`https://download.docker.com/mac/stable/Docker.dmg`
`https://download.docker.com/mac/beta/Docker.dmg`

Windows:

`https://download.docker.com/win/stable/InstallDocker.msi`
`https://download.docker.com/win/beta/InstallDocker.msi`

#### Announce the release

Announce the release in #pinata-dev with `@here` and post the final CHANGELOGs for Windows
and Mac.

Announce in #ship-it:
```
Mac and Windows Beta ... has been released!
https://download-stage.docker.com/docs/windows/release-notes/
https://download-stage.docker.com/docs/mac/release-notes/
```

Announce in #pinata:
```
Mac and Windows Beta ... has been released!
https://download-stage.docker.com/docs/windows/release-notes/
https://download-stage.docker.com/docs/mac/release-notes/
```

The release should now be complete!

## Troubleshooting

### Web UI

A basic web UI to list releases and files on the public channels is available:

http://omakase.omakase.e57b9b5b.svc.dockerapp.io/#

### Rollback Release

Rolling back to a previous release is as simple as publishing the old release.
Given a current release of: `1.12.0-beta18`, rollingback to the previous version/build would be:

Mac:

`docker-release --channel beta --build 1.12.0.9779 --human 1.12.0-rc2-beta17 --arch mac --prod publish`

Windows:

`docker-release --channel beta --build 1.12.0.5022 --human 1.12.0-rc2-beta17 --arch win --prod publish`

Note that you need to know the build number of the previous release to do this. This is usually recorded in the tracking issue, but list of releases in the [Web UI](#web ui) can also be helpful.

### Update release notes with `docker-release`
When the release is tagged the CHANGELOG is extracted and uploaded in the `NOTES` file. This file can be re-uploaded with `docker-release`.

To download `NOTES` for editing:

`curl -SLO https://download.docker.com/mac/beta/1.12.0.9996/NOTES`

To upload the updated `NOTES` file and overwrite the existing one:

`docker-release --arch mac --channel beta --build 1.12.0.9996 --prod upload NOTES`

### Surgically updating the CHANGELOG

Usually `docker-release` should provide everything you need to add/update files on S3. If you should need to manually modify the files, one way of doing it is to use the widely available `s3cmd` (`awscli` should also work in a similar way)

* configure s3cmd to point to `download.docker.com` buckets
* Note that we have two buckets one for beta/stable `s3://editions-us-east-1` another for all other channels (`s3://editions-beta-us-east-1-150610-005505`)
* The path to the proper download is : `s3:[bucket]/[os]/[channel]/` where `os` is either `mac` or `win`, and channel either `beta` or `stable`.
* list the documents in the proper channel by `s3cmd ls s3://editions-us-east-1/mac/beta/`
* find the proper version  and list files there `s3cmd ls s3://editions-us-east-1/mac/beta/1.12.0.9996/`
* the changelog is the `NOTES` files
* get the remote `NOTES` file `s3cmd get s3://editions-us-east-1/mac/beta/1.12.0.9996/NOTES`
* copy over the revelant part of `docker/pinata/CHANGELOG` for mac or `docker/pinata/win/CHANGELOG` for windows
* upload the new `NOTES` file : `s3cmd put NOTES s3://editions-us-east-1/mac/beta/1.12.0.9996/ --acl-public -m "text/plain; charset=utf-8"`
* check the public bucket URL is ok : `curl http://editions-us-east-1.s3.amazonaws.com/win/beta/1.12.0.5226/NOTES`
* if something is wrong at this point, this is probably missing metadata check them with `s3cmd info s3://editions-us-east-1/mac/beta/1.12.0.9996/NOTES`
* check the public endpoint : `curl https://download.docker.com/win/beta/1.12.0.5226/NOTES` ( you may have slight delay due to cloudfront caching)

## Hotfixes and updates on stable channel

### Processes for PRs

#### Non urgent, but necessary fixes in stable

Label the PR with `process/cherry-pick` and assign it to the relevant release milestone. This will mark the PR as a candidate to be picked into that stable release from master, when that release happens. Typically this is the next minor engine release. The PR in master should be included in one of the beta releases at least one week before the stable release.

#### Urgent stable fixes

If the issue needs to be urgently fixed, contact the release manager to make a plan for the release and to create a hotfix milestone. Typically this will require a beta hotfix first to test the fix, then a stable hotfix a week later.

#### Urgent beta-only fixes

If there is an issue in the beta that does not need to be fixed in stable, a hotfix can be released on the beta channel only. Contact the release manager to make a plan for the release. The PR does not need a label.

### Releasing a beta hotfix

A hotfix is usually a minimal update that we send out to fix an urgent issue between planned releases.

* Create a branch from the beta release tag called `bXX-release-fixes`, where XX is the beta number the hotfix is for. In some cases it may be necessary to create `bXX-release-fixes-mac` or `bXX-release-fixes-win` if there were important commits in master between the win/mac tags.
* Cherry-pick the relevant PRs from master into the new branch
* Update the CHANGELOG in the branch to include the list of hotfixes under a "* Hotfixes" title. This will be shown in the update dialog.
* Add the release to the CHANGELOG in master
* Create a release tag on the last commit on the branch (as on master) to build the release binary in CI. Add the suffix `.X` where `X` is the hotfix number for this release, e.g. `mac-v1.12.0-beta21.1`.
* Re-test the binary and release as normal. Additional testing may be required depending on the changes that were fixed.

### Releasing a stable hotfix

A stable hotfix is an urgent update on the stable channel between planned releases. The PRs included in the hotfix should previously be tested in a beta before the changes are added to the stable branch.

* Determine the minimal set of changes that needs to be included in the update, make sure all changes have been tested on beta channel first. Merge the PR in the branch for the current stable release.
* Add a letter to the version number, e.g. "1.12.0b" for the second update to `1.12.0`
* Update CHANGELOG with a friendly message describing the reason for the update (this will be displayed in the update dialog)
* Tag and release as normal - additional testing may be required depending on the changes that were fixed
