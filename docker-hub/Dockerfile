---
{}
---

FROM docs/base:oss
MAINTAINER Docker Docs <docs@docker.com>

ENV PROJECT=docker-hub

COPY . /src
RUN rm -rf /docs/content/$PROJECT/
COPY . /docs/content/$PROJECT/
