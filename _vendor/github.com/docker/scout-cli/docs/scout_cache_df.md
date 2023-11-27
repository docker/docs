# docker scout cache df

<!---MARKER_GEN_START-->
Show Docker Scout disk usage


<!---MARKER_GEN_END-->

## Description

Docker Scout uses a temporary cache storage for generating image SBOMs.
The cache helps avoid regenerating or fetching resources unnecessarily.

This `docker scout cache df` command shows the cached data on the host.
Each cache entry is identified by the digest of the image.

You can use the `docker scout cache prune` command to delete cache data at any time.

## Examples

### List temporary and cache files

```console
$ docker scout cache df
Docker Scout temporary directory to generate SBOMs is located at:
   /var/folders/dw/d6h9w2sx6rv3lzwwgrnx7t5h0000gp/T/docker-scout
   this path can be configured using the DOCKER_SCOUT_CACHE_DIR environment variable

                               Image Digest                               │ Size
──────────────────────────────────────────────────────────────────────────┼────────
  sha256:c41ab5c992deb4fe7e5da09f67a8804a46bd0592bfdf0b1847dde0e0889d2bff │ 21 kB

Total: 21 kB


Docker Scout cached SBOMs are located at:
   /Users/user/.docker/scout/sbom

                               Image Digest                               │ Size of SBOM
──────────────────────────────────────────────────────────────────────────┼───────────────
  sha256:02bb6f428431fbc2809c5d1b41eab5a68350194fb508869a33cb1af4444c9b11 │ 42 kB
  sha256:03fc002fe4f370463a8f04d3a288cdffa861e462fc8b5be44ab62b296ad95183 │ 100 kB
  sha256:088134dd33e4a2997480a1488a41c11abebda465da5cf7f305a0ecf8ed494329 │ 194 kB
  sha256:0b80b2f17aff7ee5bfb135c69d0d6fe34070e89042b7aac73d1abcc79cfe6759 │ 852 kB
  sha256:0c9e8abe31a5f17d84d5c85d3853d2f948a4f126421e89e68753591f1b6fedc5 │ 930 kB
  sha256:0d49cae0723c8d310e413736b5e91e0c59b605ade2546f6e6ef8f1f3ddc76066 │ 510 kB
  sha256:0ef04748d071c2e631bb3edce8f805cb5512e746b682c83fdae6d8c0b243280b │ 1.0 MB
  sha256:13fd22925b638bb7d2131914bb8f8b0f5f582bee364aec682d9e7fe722bb486a │ 42 kB
  sha256:174c41d4fbc7f63e1f2bb7d2f7837318050406f2f27e5073a84a84f18b48b883 │ 115 kB

Total: 4 MB
```
