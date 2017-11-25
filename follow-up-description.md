## Follow up improvements

Most pre-processing steps (if not all - if we're okay with creating archive images that use `/vX.Y/` as base-url) for the archive images can be moved to the branches for each archive.

Doing so, would make building `docs-base` just a matter of (pulling and) copying the HTML from each archive's image.

I added some scripts to test this theory.

Assuming you have build the other PR's, and tagged the archives as `docs/docker.github.io:<version>`:

First, generate the archive (this mimics the situation where all cleanup and minification steps were moved to the archive branches):

```bash
$ ./build-archive.sh
Building utilities image
Building docs-archive:v1.4
Building docs-archive:v1.5
Building docs-archive:v1.6
Building docs-archive:v1.7
Building docs-archive:v1.8
Building docs-archive:v1.9
Building docs-archive:v1.10
Building docs-archive:v1.11
Building docs-archive:v1.12
Building docs-archive:v1.13
Building docs-archive:v17.03
Building docs-archive:v17.06
```

Then, build the `docs-base` image:

```bash
$ ./build-docs-base.sh
```

```
Sending build context to Docker daemon  473.9kB
Step 1/56 : FROM docs-archive:v1.4 AS archive_v1.4
 ---> 6133277afda5
Step 2/56 : COPY --from=docs/utilities /* /usr/bin/
 ---> 946a42c67224
Step 3/56 : RUN apk add -q --no-cache gzip      && create_permalinks.sh ${TARGET} ${VER}      && compress_assets.sh ${TARGET}
 ---> Running in 83ddf79fae18
Creating permalinks for v1.4.done
Compressing assets in /usr/share/nginx/html......done.
Removing intermediate container 83ddf79fae18
 ---> 733e605f139e
Step 4/56 : FROM docs-archive:v1.5 AS archive_v1.5
 ---> 9c8c24b2ce60
Step 5/56 : COPY --from=docs/utilities /* /usr/bin/
 ---> f7e80d8fae88

...


Successfully built c644f0095bc0
Successfully tagged docs-base:latest

real	2m16.139s
user	0m0.166s
sys	0m0.076s
```


Running again (no cache bust):

```bash
$ ./build-docs-base.sh

...

Successfully built c644f0095bc0
Successfully tagged docs-base:latest

real	0m1.692s
user	0m0.096s
sys	0m0.108s
```

Now, assume an archive was updated; mimicking this by re-tagging the `v17.06` archive as `v17.03`

```bash
$ docker image tag docs-archive:v17.06 docs-archive:v17.03
```

And rebuild the docs-base image:

```bash
$ ./build-docs-base.sh
```

Everything up to `v1.13` is still cached :+1: :+1::

```
Step 48/56 : COPY --from=archive_v1.13  ${TARGET} ${TARGET}/v1.13
 ---> Using cache
 ---> 65003dcc3ecb
Step 49/56 : COPY --from=archive_v17.03 ${TARGET} ${TARGET}/v17.03
 ---> b2e3adf04a5b
Step 50/56 : COPY --from=archive_v17.06 ${TARGET} ${TARGET}/v17.06
 ---> d9155d86c6c4
Step 51/56 : ARG REPORT_SIZE
 ---> Running in 8e7ce8dc4e4c
Removing intermediate container 8e7ce8dc4e4c
 ---> b7fe2ba1b15f
Step 52/56 : COPY ./scripts/size_report.sh /usr/bin/
 ---> 7ef051513727
Step 53/56 : RUN       if [ -n "$REPORT_SIZE" ]; then           apk add -q --no-cache coreutils;           size_report.sh ${TARGET};       fi;
 ---> Running in 51e43dbde1b6
Removing intermediate container 51e43dbde1b6
 ---> 7cb0325d5232
Step 54/56 : COPY index.html ${TARGET}/
 ---> fcc88b99afeb
Step 55/56 : COPY ./default.conf  /etc/nginx/conf.d/default.conf
 ---> 3358db2c1b8a
Step 56/56 : CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000"; exec nginx -g 'daemon off;'
 ---> Running in 81cf9a3753d9
Removing intermediate container 81cf9a3753d9
 ---> 5063982f3a03
Successfully built 5063982f3a03
Successfully tagged docs-base:latest

real	0m11.981s
user	0m0.100s
sys	0m0.042s
```

