# Roadmap process

The roadmap is tracked with GitHub issues with the "kind/roadmap" label. These things are large pieces of work that typically need to be delivered by a certain date, and may have complex dependencies that need tracking.

Using the "kind/roadmap" label allows us to filter out the major items of work from all the other smaller stuff on GitHub (bugs, feature requests, sub-tasks, etc).

The text of roadmap issues can be used as a wiki to give detail on that piece of work, keep track of its status, and link off to other relevant stuff. This can also link off to other sub-tasks if more granularity is necessary.

Issues should be given:

 - **A milestone** to define its aimed delivery date
 - **An owner** to define its DRI
 - **A priority** as a label to define its priority relative to other items

## Reports

GitHub is the source of truth for the roadmap, but GitHub also has a great API that allows us to produce reports to give visibility into what is going on both within the team and to people outside the team.

 - [Charing Cross](https://charingcross.i.dckr.io/repos/docker/pinata), a tool for viewing milestones and the status of roadmap issues within those milestones.

The issues can also be browsed using GitHub's interface. For example:

 - [All roadmap issues](https://github.com/docker/pinata/issues?q=is%3Aopen+is%3Aissue+label%3Akind%2Froadmap)
 - [All roadmap issues for "GA" milestone](https://github.com/docker/pinata/issues?q=is%3Aopen+is%3Aissue+label%3Akind%2Froadmap+milestone%3AGA) (modify milestone filter to view others)
