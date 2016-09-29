# Reclaiming disk space on the Mac

See [docker/pinata#2080]

Engine 1.13 will include a feature allowing disk space to be reclaimed more
easily by deleting unused images. Unfortunately deleting files within the Moby VM
will not delete disk blocks stored on the host unless the block device supports
ATA TRIM / SCSI DISCARD. Note that Hyper-V supports AAT TRIM / SCSI DISCARD so this
should "just work" on Windows.

## The datapath

```
  docker
    |
    |  rm unused images
    V
  ext4
    |
    |  [1] TRIM unused blocks
    V
  ahci-hd [2] (virtio-blk doesn't support TRIM)
    |
    |   TRIM unused blocks
    V
  hyperkit [3] (needs to expose the capability, currently commented out)
    |
    |   TRIM unused blocks
    V
  ocaml-qcow
    |
    | GC [4]
    V
   HFS+
```
where

- [1] we need to investigate/decide whether TRIM should be run continuously
      as a side-effect of an `rm` (enabled via mounting with the `discard` flag)
      or whether we should have a `fstrim` running later, possibly on request?
      This will likely depend on the performance of the rest of the system.
      See [comments and notes on the Arch Wiki](https://wiki.archlinux.org/index.php/Solid_State_Drives).
- [2] unfortunately the virtio-blk protocol doesn't support TRIM so we have to
      switch to the only other option in hyperkit: ahci-hd. Note that qemu/KVM
      users would probably use virtio-scsi. We will have to investigate the
      performance implications.
- [3] hyperkit's bhyve heritage includes support for exposing TRIM, but
      the [capability is not exposed](https://github.com/docker/hyperkit/blob/547caeb5facb248067c529dd8c80931dbc1c56c6/src/block_if.c#L620)
      and the operation [returns ENOTSUPP](https://github.com/docker/hyperkit/blob/547caeb5facb248067c529dd8c80931dbc1c56c6/src/block_if.c#L416).
      This is clearly because it's not implemented by existing filesystems on
      the Mac.
- [4] since the Mac doesn't support TRIM on HFS+, we will need to implement a
      block GC to shrink the file.

## The qcow block GC

The general approach is to shrink the qcow2 file with `ftruncate` after copying
the trailing data blocks back to "holes" created by the TRIM. A shuffle operation
will be approximately:

1. Find the pointer which is pointing to the last physical cluster. This will
   require a reverse index to be created and maintained, since all the references
   within the physical file point in the other direction. We can make the
   simplifying assumption that the refcount is 1 since we don't support internal
   snapshots.
2. Lock the last physical cluster and a TRIMed cluster, to prevent concurrent
   modification from another concurrent thread.
3. Copy the last block into a TRIMed cluster. Note that reads from the TRIMed
   cluster will still return zeroes because the TRIM set will be consulted before
   reading the physical disk.
4. Issue a `flush`
5. Rewrite the pointer to point to the TRIMed cluster.
6. Issue a `flush`
7. Issue a `ftruncate`

Note the `flush` calls are necessary as write barriers to prevent the pointer
update from being persisted without the data copy over a crash.

Note the `flush` is implemented by `fcntl(F_FULLFSYNC)` on OSX, as described
by the `fsync` man page or [this email](http://lists.apple.com/archives/darwin-dev/2005/Feb/msg00072.html).
On a modern MBP, a `flush` takes about 10ms.

### The reverse cluster index

A qcow2 file is a tree of 64 KiB clusters where the root of the tree is
in the header at the beginning of the file. All clusters are linked either
from the header, or from other clusters. In order to shuffle a block within
the file we need to find the reference to it, so we can update the pointer
with the new location.
Unfortunately the qcow2 format only stores pointers to child blocks in the tree,
not pointers to parents, so we need to create and maintain our own reverse
cluster index.

A table mapping clusters to their parents would need to have one entry per
cluster i.e. 3 bytes per 64 KiB of data i.e. an extra 0.005% space overhead.

The largest possible qcow2 image is [2**63-513 bytes long](https://rwmj.wordpress.com/2011/10/03/maximum-qcow2-disk-size/). Rounding
the size up to `2**63`, the required overhead would be about 300 TiB.

The top-of-the-range Mac Pro can be configured with 1 TiB of storage. Assuming
hard drive capacity doubles every 2 years then by 2020 we could have a 4 TiB
qcow2. The required overhead would be about 192 MiB.

Therefore, although it's possible to construct a .qcow2 whose reverse cluster
index would need to be stored on disk, for the next couple of years at least
for our use-case we could use an in-memory table.

### Making the TRIM asynchronous

The block shuffling described above will be slow since it calls `flush` twice.
The `TRIM` command is supposed to be fast, and tools like `mkfs` will issue
thousands of them back-to-back. If we delete a lot of files all at once then
`ext4` could generate a lot of small `TRIM`s which could be coalesced and
executed all at once.

We could make the TRIM asynchronous by first marking the clusters as empty
(such that reads will return zeroes) in a persistent datastructure cached
in memory. The ideal
structure would have fast lookup (since it's on the read path) and would
automatically coalesce adjacent TRIMed regions together.
Martin Erwig's paper [Diets for Fat Sets](http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.26.4659)
describes a structure with fast lookup and automatic coalesce, and there is an OCaml implementation
inside the [batteries library](https://ocaml-batteries-team.github.io/batteries-included/hdoc2/BatISet.html).
To make the changes persistent we could use a simple append-only log to store
updates, and occasionally flush the whole structure to disk. The
[mirage/shared-block-ring](https://github.com/mirage/shared-block-ring) library
is designed for this use-case and was used to speed up LVM metadata operations
in the xapi-project.

## Making the shuffling faster

After we have coalesced a large number of TRIMs together, we can shuffle the
blocks more efficiently by batching. In the description above, one block was
moved at a time. We can extend this to an arbitrary number of blocks, with
the caveat that they must be locked against writes for the duration of the
shuffle. We can therefore amortise the cost of the `flush` over many blocks.

## Possible development phases

### Phase 1

If we expose the TRIM capability and track the TRIMed blocks, we could create
an offline tool for geneating a small qcow2 (rather than shuffling blocks).
This would allow motivated users to get their space back.

### Phase 2

We could run the offline tool ourselves between reboots.

### Phase 3

In this phase we could shuffle the blocks around live.
