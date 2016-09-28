# GC

This is the GC binary which we run through the garbage collector job defined in
`garbagecollector/`.

# How it works

A casual mark-and-sweep GC for blobs in the datastore. It's a bit different from
the one in `docker/distribution`; we have a metadata store listing all manifests
and tags. Luckily, we also list which layers a manifest uses within the manifest
table.

This means we can mark by querying rethinkDB, then enumerate over all blobs and
delete if necessary.

Le steps:

1. Load all manifests from rethink
2. Create a map of all layers referenced from each manifest
3. Iterate over each blob in the blobstore
4. Delete the blob if it's not marked in our map

Le super simple.

In the future maybe we'll take a look at some concurrent shizzle. That would be
godly... no read-only mode! Oh my.

---

We be cleaning like

![](http://i.imgur.com/bFlHsIg.gifv)
