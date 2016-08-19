---
{}
---

FROM docs/base:oss
MAINTAINER Docker Docs <docs@docker.com>

ENV PROJECT=docker-cloud
# To get the git info for this repo
COPY . /src
RUN rm -rf /docs/content/$PROJECT/
COPY . /docs/content/$PROJECT/

# This kludge only exists when run from the docker-cloud repo (useful for testing)
RUN rm -rf /docs/content/apidocs \
	&& mv -f /docs/content/$PROJECT/apidocs /docs/content/ || true
RUN mkdir -p /docs/layouts/cloud-api-docs/
RUN mv /docs/content/apidocs/layouts/* /docs/layouts/cloud-api-docs/ || true
