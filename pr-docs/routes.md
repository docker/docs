# Routes

Two items affect this proposal.

1. [Distribution's work](https://github.com/docker/distribution),
   specifically relating to defining Repositories, Images, Manifests,
   Digests and Tags.
2. The Current Hub's Routing Issues

## Distributions Work (partially summarized)

### Repository

* A set of blobs
* Subsets of these blobs make up Images

### Image

* A set of blobs
  - Layers
  - Tag
  - Signatures
  - Manifest
* A Tag (potentially containing signatures) points to a Manifest
* A Manifest points to multiple layers.

### Manifest

As defined in the [distribution][manifest-pr] Manifest PR:

> A [Content Manifest][manifest] is a simple JSON file which contains
> general fields that are typical to any package management
> system. The goal is for these manifests to describe an application
> and its dependencies in a content-addressable and verifiable way.

### Tag

As defined in the [distribution][d-tag-pr] PR:

> A [tag][tag] is simply a named pointer to content. The content can
> be any blob but should mostly be a manifest. One can sign tags to
> later verify that they were created by a trusted party.

### Additional Content

Image names will be allowed to have many slashes in the future.

## Current Hub Issues

### Collisions

The URLs for user and repo collide:

A user's Starred Repos:

```
/u/biscarch/starred/
```

A user's repository, named Starred.

```
/u/biscarch/starred/
```

## Future Problems

An image for the user `biscarch`, named `my/repo`:

```
/u/biscarch/my/repo/
```

An image for the user `biscarch`, named `my`, tagged `repo`:

```
/u/biscarch/my/repo/
```

## Solutions

Namespace `Users`, `Repos` and `Images` as such (with the user
`biscarch`)

```
/u/:user
/r/:user/:repo
/i/:user/:repo/:tag
```

### Solving Starred Repos

Prefix defines whether we are referring to a repo or attribute of a
user:

```
/u/biscarch/starred
/r/biscarch/starred
```

### Solving Repo/Image Conflicts

Prefix determines whether we are referring to a Repository or Image:

```
/r/biscarch/my/repo/
/i/biscarch/my/repo/
```

## The new Spec

```
/u/
/u/:user/
/r/:user/:repo/
/i/:user/:repo/:tag/
```

### Full List

### Dashboard

`/`

### Official Repositories

"username" === library, which is represented as the root `_`.
All management of `library` namespaced repos is done from the usual
`/u/library/:repo/`

```
/_/:repo/
/_/:repo/dockerfile/
/_/:repo/dockerfile/raw
/_/:repo/tags/
```

### Single Endpoints

* Search
  - `/search/`
* Plans
  - `/plans/`

### Account

Mostly Settings; Add Repository Page;

`/account/` should redirect to `/account/settings/`

```
/account/accounts/
/account/authorized_services/
/account/change-password/
/account/confirm-email/<var>/
/account/emails/
/account/notifications/
/account/organizations/
/account/organizations/:org_name/
/account/organizations/:org_name/groups/:group_id/
/account/repositories/add/
/account/settings/
/account/subscriptions/
```

### Users

```
/u/
/u/:user/
/u/:user/activity/
/u/:user/contributed/
/u/:user/starred/
```

### Repos

```
/r/:user/:repo/
/r/:user/:repo/~/settings/
/r/:user/:repo/~/settings/collaborators/
/r/:user/:repo/~/settings/links/
/r/:user/:repo/~/settings/triggers/
/r/:user/:repo/~/settings/webhooks/
/r/:user/:repo/~/settings/tags/
```

Current build history urls:

```
/r/:user/:repo/~/builds_history/
```

### Images

We currently don't do a lot for Images. Repositories have been the
main focus.

```
/i/:user/:repo/:tag/
/i/:user/:repo/:tag/~/dockerfile/
/i/:user/:repo/:tag/~/dockerfile/raw/
```

### Automated Builds

```
/automated-builds/
/builds/
/builds/:user/:repo/
```

### Convenience Redirects

Also, potential pages to build out more agressively.

* `/official/`
  - redirects to `/search?q=library&f=official`
  - future: Potentially `Explore` type page for official repos
* `/most_stars/`, `/popular/`
  - redirects to `search?q=library&s=stars`
* `/recent_updated/`
  - `search?q=library&s=last_updated`

#### Help

* `/help`
  - `https://www.docker.com/resources/help/`
  - Can we rely on this url to stick around?
* `/help/docs`
  - `https://docs.docker.com/`

## Make Separate Sites for:

### Highland URLs

We need to pull out the APIs used on the current Hub for this.

```
/highland/
/highland/build-configs/
/highland/builds/
/highland/search/
/highland/stats/
```


## More Issues

* There are no links to comments
* `/opensearch.xml` times out on the current site
  - Should we re-implement?
* `/sitemap.xml`

# Concerns with this Proposal

* Automated Build urls need to be given more thought

[tag]: https://github.com/stevvooe/distribution/blob/a8d3f3474b7b60576dc64250d95db3717bf07c33/doc/spec/tags.md#tags
[d-tag-pr]: https://github.com/docker/distribution/pull/173/files
[d-manifest-pr]: https://github.com/docker/distribution/pull/62
[manifest]: https://github.com/jlhawn/distribution/blob/e8b5c8c32b565b9b643c3a0b0e87339bf40eb206/doc/spec/manifest.md
