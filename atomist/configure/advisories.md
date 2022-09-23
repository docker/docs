---
title: Advisory sources
description: Add and remove vulnerability advisories
keywords: >
  advisories, vulnerabilities, databases, open source, configure, security,
  atomist
---

{% include atomist/disclaimer.md %}

With no configuration required, Atomist already draws vulnerability data from
several public advisories. You can extend this by adding your own, custom
advisories if you wish.

## Adding and updating advisories

To add your own advisories:

1. Create a repository called `atomist-advisories` in the GitHub account where
   you've installed the Atomist GitHub app.

2. In the default branch of the repository, add a new JSON file called
   `<source>/<source id>.json`, where:

   - `source` should be the name of your company
   - `source-id` has to be a unique id for the advisory within `source`.

3. The JSON file must follow the schema defined in
   [Open Source Vulnerability format](https://ossf.github.io/osv-schema/){:
   target="blank" rel="noopener" class=""}.

   Refer to the
   [GitHub Advisory Database](https://github.com/github/advisory-database/tree/main/advisories/github-reviewed){:
   target="blank" rel="noopener" class=""} for examples of advisories.

## Deleting advisories

Delete an advisory from the database by removing the corresponding JSON advisory
file from the `atomist-advisories` repository.

> **Note**
>
> Atomist only considers additions, changes and removals of JSON advisory files
> in the repository's **default branch**.
