Piñata Support
==============

This document outlines the roles and responsibilities of the support function
of Piñata.

# Roles and responsibilities

## L1 Support

The L1 Support staff are responsible for:

- Categorising and labelling issues in `docker/for-mac` and `docker/for-win`
- De-duplicating issues in the aforementioned repositories
- Ensuring that the issues are of high quality
  - They contain a `docker-diagnose` output
  - The logs for the report are available and uploaded to S3
    - This can be verified with `Dockerfile.fetch` in `support/s3`
  - There is sufficient information about the users environment
    - Is this a new/failed installation
    - An upgrade
    - A channel switch
  - Reproduction instructions are present. The following are acceptable:
    - A sequence of UI interactions that will cause the issue
    - A shell script that will trigger the issue
    - A `Dockerfile` that will trigger the issue
    - A `docker-compose.yml` that will trigger the issue
  - The reproduction instructions are valid, and reproduce the issue
  - **IF** the cause is obvious, open an issue on `docker/pinata` and ping
    the right person
    - This issue should contain all of the necessary information,
      not just a link to the `docker/for-X` bug
  - Sending a report at the end of a rotation to `pinata@docker.com` or
    `@pinata-team` on Slack that includes:
    - Major occurrences
    - Any patterns or trends identified
    - Statistics:
      - Issues Opened vs Closed this week
      - Issues Opened vs Close all time
      - Average time to resolution
      - Tickets over N days old, where N is the targeted resolution time

## L2 Support

The L2 Support staff are responsible for:

- Reviewing categorised issues with reproduction cases to identify the
  root cause
- Once that cause has been found, to create an issue in `docker/pinata` and
  assign the right person
- They should track the open issue to make sure it is resolved within a
  reasonable timeframe (SLA TBD)
- Once the issue has been fixed (or not), they should update the issue in the
  upstream tracker with the appropriate label
- Once the fix has been applied to a publicly available release, the
  appropriate label should be applied

## L3 Support

In the proposed model L3 support are engineering.
L3 support are responsible for:

- Reviewing issues where they have been assigned by L2 support
- Either fixing the issue, or marking as `0-wont-fix`

# Status Labels (to be documented on `docker/for-X`)

The following labels track the progression through the support process for
ease of triage.

| Label                     | Description                                                            |
|---------------------------|------------------------------------------------------------------------|
| status/0-triage           | The issue needs triaging                                               |
| status/0-wont-fix         | This issue will not be fixed and therefore can be closed                |
| status/0-more-info-needed | The issue needs more information before it can be triaged              |
| status/1-acknowledged     | The issue has been triaged and is ready for L2                         |
| status/2-in-progress      | The issue has been assigned to a engineer and is waiting a fix         |
| status/3-fixed            | The issue has been fixed in `master`                                   |
| status/4-fix-released-beta  | The fix has been released!                                           |
| status/4-fix-released-stable | The fix has been released!                                          |
