---
title: Prevent tags from being overwritten
description: Learn how to make your tags immutable and don't allow users to
  overwrite a tag if it already exists
keywords: registry, immutable
---

{% assign domain="dtr.example.org" %}
{% assign org="library" %}
{% assign repo="wordpress" %}
{% assign tag="latest" %}

By default, users with access to push to a repository, can push the same tag
multiple times to the same repository.
As an example, a user pushes an image to `{{ org }}/{{ repo }}:{{ tag }}`, and later another
user can push the image with exactly the same name but different functionality.
This might make it difficult to trace back the image to the build that generated
it.

To prevent this from happening, you can configure a repository to be immutable.
Once you push a tag, DTR won't allow anyone else to push another tag with the same
name.

## Make tags immutable

To make tags immutable, in the **DTR web UI**, navigate to the
**repository settings** page, and change **Immutability** to **On**.

![](../../images/immutable-repo-1.png){: .with-border}

From now on, users will get an error message when trying to push a tag
that already exists:

```bash
docker push {{ domain }}/{{ org }}/{{ repo }}:{{ tag }}
unknown: tag={{ tag }} cannot be overwritten because {{ domain }}/{{ org }}/{{ repo }} is an immutable repository
```

## Where to go next

- [Sign images](sign-images/index.md)