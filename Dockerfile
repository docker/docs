---
{}
---

FROM docs/base:oss
MAINTAINER Docker Docs <docs@docker.com>

# because both the 2 dir's are going into the root
env PROJECT=

# To get the git info for this repo
COPY . /src
#RUN rm -rf /docs/content/$PROJECT/
COPY . /docs/content/$PROJECT/
