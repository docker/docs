FROM docs/base:oss
MAINTAINER Mary Anthony <mary@docker.com> (@moxiegirl)

ENV PROJECT=docker-trusted-registry
# To get the git info for this repo
COPY . /src
RUN rm -rf /docs/content/$PROJECT/
COPY . /docs/content/$PROJECT/

# This kludge only exists when run from the DTR repo (useful for testing)
RUN mv -f /docs/content/$PROJECT/apidocgen/ /docs/content/apidocs/ || true
