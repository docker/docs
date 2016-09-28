# Triaging Incoming Issues

Every week one engineer on the team will be responsible for triaging
incoming issues.  The following doc gives general guidelines for this
process, but you should always use your best judgment.  Our goal is to
have ~zero issues without asignee, milestone, and kind+area labels.


## Rotation

The rotation runs Wednesday through Tuesday.  Please remind the next
person in the rotation on the last day of your turn.

* Adrian
* Alex
* Arunan
* Daniel
* Evan
* Josh
* Tom
* Sean


## Labels

Make sure every issue has the following labels attached:

* [label] `kind/*` (one or more)
* [label] `area/*` (one or more)
* [label] `status/*` - Use as needed


## Assignee

Make sure every issue has an assigned owner.  If the issue is a "trivial"
fix, consider picking it up yourself even if you don't normally work
on the relevant section of code to get some more experience across
the codebase.

If the issue is a `kind/feature-request` and is more than a trivial
enhancement, assign to `vsaraswat` so it can be prioritized
by the PM team. (Also set the milestone to `future` and set
`status/needs-attention`)

If there's an obvious engineer to assign it to, do so.  Failing that,
you can use the following list to pick an asignee.

| Area | Default Owner | Description |
|------|---------------|-------------|
| area/analytics    | Alex      | Mixpanel integration |
| area/authN        | Josh      | Authentication within UCP |
| area/authZ        | Alex      | Authorization within UCP |
| area/bootstrapper | Adrian    | The install/upgrade/etc. bootstrap tool |
| area/certificates | Josh      | Certificate authority, CSRs, bundles, etc. |
| area/clustermgmt  | Alex      | Node management within the cluster |
| area/config       | Adrian    | UCP configuration |
| area/core         | Daniel    | Integration with core team components (engine, engine-api, etc.) |
| area/dtr          | Alex      | DTR integration |
| area/enzi         | Josh      | Enzi authentication service |
| area/installation | Adrian    | Installation flows |
| area/licensing    | Daniel    | License management |
| area/misc         | Daniel    | Grab bag when it doesn't match something else |
| area/network   | Daniel    | Overlay networking, and other networking issues in the cluster |
| area/stability    | Alex      | Issues related to stability of the infrastructure |
| area/swarm        | Alex      | Issues related to swarm or the integration with swarm |
| area/ui           | Tom       | GUI |
| area/upgrade      | Adrian    | Upgrade specific issues |
| area/visibilty    | Tom       | Cluster wide visibility of health, status, etc. |


## Milestone

Every triaged issue should have a milestone set.  When in doubt, assign
to `future` for the backlog.  This helps submitters understand that the
issue isn't currently scoped for a specific release and gives them the
opportunity to speak up if they feel the issue is urgent.

If you feel you have enough information, please try to assign to the
appropriate release.  You can consult the asignee and slot into the
current release if they feel they have the bandwidth to pick it up.

Generally try to avoid assigning "random" issues to the N+1 release,
unless they specifically apply to the themes or existing planned
features/tasks for that release.  Better to assign these to the "future"
milestone so they can be prioritized with other backlog items once the
N+1 release enters the planning phase.

## Documentation

Ideally, we would ship docs when we ship a feature. But right now we
have a big docs debt, and few resources to create docs, so we need to
be careful with the way we prioritize docs issues.

All documentation issues should be labeled with `kind/documentation`
and assigned to `Joao`.

Documentation issues that needs to go out on the day of a the release
need to be set with a milestone. Some examples of issues that need to
go out on the day of the release:

* Install docs
* Upgrade / migrate docs
* Release notes / breaking changes
* Reference docs (API, CLI, ...)

All other (maintenance) issues should be assigned to the `future`
milestone. These docs are released continuously on a release cycle that
is detached from the product cycle.  Docs maintenance issues are added to
a priority queue that is shared by all CAAS products. Joao prioritizes
these issues with feedback from devs, PMs, support, and field. Then the
issues with highest priority are solved.
