---
advisory: Custom advisory text
description: test index page
draft: true
keywords:
- docker, documentation, about, technology, understanding,  release
menu:
  main:
    parent: mn_test
title: Test index page
---

# test index page

> **Note**: To enable autoredeploy on an image stored in a third party registry,
> you will need to use [redeploy triggers](triggers.md) instead.

<p><a class="button darkblue-btn" href="https://dyhfha9j6srsj.cloudfront.net/Docker.dmg">Get Docker for Mac</a></p>


1. There are
2. several
3. things
   that are good, but `docker run --rm -it` is lots to type all the time.
   
   - to
   - know
   
   ```Dockerfile
   FROM alpine
   MAINTAINER SvenDowideit <SvenDowideit@docker.com>
   ```

Some of the daemon's options are:

| Flag                  | Description                                               |
|-----------------------|-----------------------------------------------------------|
| `-D`, `--debug=false` | Enable or disable debug mode. By default, this is false. |
| `-H`,`--host=[]`      | Daemon socket(s) to connect to.                           |
| `--tls=false`         | Enable or disable TLS. By default, this is false.         |


   
## When writing markdown

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent a lorem et urna dictum faucibus. Morbi a nunc consequat neque feugiat commodo. Donec placerat sit amet eros at luctus. Sed non elementum lectus. Nunc lacus neque, dignissim in nunc eu, sollicitudin congue leo. Proin tempus metus ac lacinia viverra. Aenean tortor neque, egestas a nisl tempus, tempus gravida eros. Vivamus efficitur ipsum id est congue, sed rhoncus neque consectetur. Sed enim ante, mollis ut erat sit amet, semper euismod erat. Donec tincidunt consequat varius. Suspendisse potenti. Nullam blandit pellentesque consectetur. Nulla bibendum turpis tellus, sit amet tristique quam ultrices quis. Nam commodo eleifend lorem, vel suscipit nisl viverra et. Maecenas viverra lobortis nibh ac aliquet.

Lorem ipsum dolor sit amet, **consectetur adipiscing** elit. Praesent a lorem et urna dictum faucibus. Morbi a nunc consequat neque feugiat commodo. Donec placerat sit amet eros at luctus. Sed non elementum lectus. Nunc lacus neque, dignissim in nunc eu, sollicitudin congue leo. Proin tempus metus ac lacinia viverra. Aenean tortor neque, egestas a nisl tempus, tempus gravida eros. Vivamus efficitur ipsum id est congue, sed rhoncus neque consectetur. Sed enim ante, mollis ut erat sit amet, semper euismod erat. Donec tincidunt consequat varius. Suspendisse potenti. Nullam blandit pellentesque consectetur. Nulla bibendum turpis tellus, sit amet tristique quam ultrices quis. Nam commodo eleifend lorem, vel suscipit nisl viverra et. Maecenas viverra lobortis nibh ac aliquet.

Lorem ipsum dolor sit amet, consectetur adipiscing elit. **Praesent a lorem et** urna dictum faucibus. Morbi a nunc consequat neque feugiat commodo. Donec placerat sit amet eros at luctus. Sed non elementum lectus. Nunc lacus neque, dignissim in nunc eu, sollicitudin congue leo. Proin tempus metus ac lacinia viverra. Aenean tortor neque, egestas a nisl tempus, tempus gravida eros. Vivamus efficitur ipsum id est congue, sed rhoncus neque consectetur. Sed enim ante, mollis ut erat sit amet, semper euismod erat. Donec tincidunt consequat varius. Suspendisse potenti. Nullam blandit pellentesque consectetur. Nulla bibendum turpis tellus, sit amet tristique quam ultrices quis. Nam commodo eleifend lorem, vel suscipit nisl viverra et. Maecenas viverra lobortis nibh ac aliquet.

Lorem ipsum dolor sit amet, consectetur adipiscing elit. __Praesent a lorem__ et urna dictum faucibus. Morbi a nunc consequat neque feugiat commodo. Donec placerat sit amet eros at luctus. Sed non elementum lectus. Nunc lacus neque, dignissim in nunc eu, sollicitudin congue leo. Proin tempus metus ac lacinia viverra. Aenean tortor neque, egestas a nisl tempus, tempus gravida eros. Vivamus efficitur ipsum id est congue, sed rhoncus neque consectetur. Sed enim ante, mollis ut erat sit amet, semper euismod erat. Donec tincidunt consequat varius. Suspendisse potenti. Nullam blandit pellentesque consectetur. Nulla bibendum turpis tellus, sit amet tristique quam ultrices quis. Nam commodo eleifend lorem, vel suscipit nisl viverra et. Maecenas viverra lobortis nibh ac aliquet.
### level 3
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent a lorem et urna dictum faucibus. Morbi a nunc consequat neque feugiat commodo. Donec placerat sit amet eros at luctus. Sed non elementum lectus. Nunc lacus neque, dignissim in nunc eu, sollicitudin congue leo. Proin tempus metus ac lacinia viverra. Aenean tortor neque, egestas a nisl tempus, tempus gravida eros. Vivamus efficitur ipsum id est congue, sed rhoncus neque consectetur. Sed enim ante, mollis ut erat sit amet, semper euismod erat. Donec tincidunt consequat varius. Suspendisse potenti. Nullam blandit pellentesque consectetur. Nulla bibendum turpis tellus, sit amet tristique quam ultrices quis. Nam commodo eleifend lorem, vel suscipit nisl viverra et. Maecenas viverra lobortis nibh ac aliquet.
#### level 4
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent a lorem et urna dictum faucibus. Morbi a nunc consequat neque feugiat commodo. Donec placerat sit amet eros at luctus. Sed non elementum lectus. Nunc lacus neque, dignissim in nunc eu, sollicitudin congue leo. Proin tempus metus ac lacinia viverra. Aenean tortor neque, egestas a nisl tempus, tempus gravida eros. Vivamus efficitur ipsum id est congue, sed rhoncus neque consectetur. Sed enim ante, mollis ut erat sit amet, semper euismod erat. Donec tincidunt consequat varius. Suspendisse potenti. Nullam blandit pellentesque consectetur. Nulla bibendum turpis tellus, sit amet tristique quam ultrices quis. Nam commodo eleifend lorem, vel suscipit nisl viverra et. Maecenas viverra lobortis nibh ac aliquet.
## level 2
```bash
root@7d7ed2558cd3:~/sven/src/docs-base-1.12-integration/content# make docs
make: *** No rule to make target 'docs'.  Stop.
root@7d7ed2558cd3:~/sven/src/docs-base-1.12-integration/content# cd ..
root@7d7ed2558cd3:~/sven/src/docs-base-1.12-integration# make docs
docker build -t "docs-base:master" .
Sending build context to Docker daemon 8.087 MB
Step 1 : FROM debian:jessie
 ---> 47af6ca8a14a              && apt-get install -y gettext git wget libssl-dev make python-dev python-pip python-setuptools subversion-tool
s vim-tiny ssed curl ary&& apt-get clean cker.co&& rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
 ---> Using cache
 ---> 79f09c378b62
Step 4 : RUN pip install awscli==1.4.4 pyopenssl==0.12
 ---> Using cache
 ---> 5a93a77583a9
Step 5 : ENV HUGO_VERSION 0.16-pre3
 ---> Using cache
 ---> c8925b6323c8
Step 6 : RUN curl -sSL -o /usr/local/bin/hugo https://github.com/docker/hugo/releases/download/${HUGO_VERSION}/hugo  && chmod 755 /usr/local/b
in/hugo  && /usr/local/bin/hugo version
 ---> Using cache
 ---> bc7cde498ed7
Step 7 : RUN curl -sSL -o /usr/local/bin/markdownlint https://github.com/docker/markdownlint/releases/download/v0.9.5/markdownlint  && chmod 7
55 /usr/local/bin/markdownlint
 ---> Using cache
 ---> c3101af070ff
Step 8 : RUN curl -sSL -o /usr/local/bin/linkcheck https://github.com/docker/linkcheck/releases/download/v0.3/linkcheck  && chmod 755 /usr/loc
al/bin/linkcheck
 ---> Using cache
 ---> 116e258eda46
Step 9 : WORKDIR /docs
 ---> Using cache
 ---> 0bae93a3fb0c
Step 10 : COPY . /docs
 ---> 72cc369c38e4
Removing intermediate container 845e09bb8df8
Step 11 : RUN chmod 755 /docs/validate.sh
 ---> Running in 4d4494c7b56d
 ---> eb09024649aa
Removing intermediate container 4d4494c7b56d
Step 12 : EXPOSE 8000
 ---> Running in ac42def06104
 ---> 0739c6007a34
Removing intermediate container ac42def06104
Step 13 : CMD /docs/validate.sh
 ---> Running in d595b2f0a12b
 ---> 8504027ae838
Removing intermediate container d595b2f0a12b
Successfully built 8504027ae838
docker run --rm -it  -e AWS_S3_BUCKET -e NOCACHE --name docker-docs-tools -p 8000:8000 -e DOCKERHOST "docs-base:master" hugo server --port=800
0 --baseUrl=localhost --bind=0.0.0.0 --config=./config.toml
INFO: 2016/06/06 00:30:31 hugo.go:455: Using config file: ./config.toml
WARN: 2016/06/06 00:30:31 hugo.go:549: Unable to find Static Directory: /docs/static/
INFO: 2016/06/06 00:30:31 hugo.go:558: /docs/themes/docker-2016/static is the only static directory available to sync from
INFO: 2016/06/06 00:30:31 hugo.go:601: syncing static files to /
Started building site
INFO: 2016/06/06 00:30:31 site.go:1248: found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
WARN: 2016/06/06 00:30:31 site.go:1998: Unable to locate layout for section test: [section/test.html _default/section.html _default/list.html
indexes/test.html _default/indexes.html theme/section/test.html theme/_default/section.html theme/_default/list.html theme/indexes/test.html t
heme/_default/indexes.html theme/section/test.html theme/_default/section.html theme/_default/list.html theme/indexes/test.html theme/_default
/indexes.html theme/theme/section/test.html theme/theme/_default/section.html theme/theme/_default/list.html theme/theme/indexes/test.html the
me/theme/_default/indexes.html]
WARN: 2016/06/06 00:30:31 site.go:1974: "test" is rendered empty
WARN: 2016/06/06 00:30:31 site.go:1998: Unable to locate layout for section : [section/.html _default/section.html _default/list.html indexes/
.html _default/indexes.html theme/section/.html theme/_default/section.html theme/_default/list.html theme/indexes/.html theme/_default/indexe
s.html theme/section/.html theme/_default/section.html theme/_default/list.html theme/indexes/.html theme/_default/indexes.html theme/theme/se
ction/.html theme/theme/_default/section.html theme/theme/_default/list.html theme/theme/indexes/.html theme/theme/_default/indexes.html]
WARN: 2016/06/06 00:30:31 site.go:1974: "" is rendered empty
0 of 3 drafts rendered
0 future content
4 pages created
18 non-page files copied
0 paginator pages created
0 tags created
0 categories created
in 36 ms
WARN: 2016/06/06 00:30:31 hugo.go:625: Skip LayoutDir: lstat /docs/layouts: no such file or directory
ERROR: 2016/06/06 00:30:31 hugo.go:629: Walker:  lstat /docs/static: no such file or directory
Watching for changes in /docs/{data,content,themes}
WARN: 2016/06/06 00:30:31 hugo.go:625: Skip LayoutDir: lstat /docs/layouts: no such file or directory
ERROR: 2016/06/06 00:30:31 hugo.go:629: Walker:  lstat /docs/static: no such file or directory
Serving pages from memory
Web Server is available at http://localhost:8000/ (bind address 0.0.0.0)
Press Ctrl+C to stop
```
