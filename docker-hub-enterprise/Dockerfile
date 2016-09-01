FROM docs/base:latest
MAINTAINER Sven Dowideit <SvenDowideit@docker.com> (@SvenDowideit)

# to get the git info for this repo
COPY . /src

# Reset the /docs dir so we can replace the theme meta with the new repo's git info
RUN git reset --hard

COPY version /docs/VERSION
COPY documentation/* /docs/sources/docker-hub-enterprise/
# sadly, COPY is only recursive if you don't also need to name the destination
COPY documentation/assets/ /docs/sources/docker-hub-enterprise/assets/
COPY documentation/mkdocs.yml /docs/mkdocs-dhe.yml

# Then build everything together, ready for mkdocs
RUN /docs/build.sh
