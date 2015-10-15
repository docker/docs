FROM docs/base:hugo-github-linking
MAINTAINER Mary Anthony <mary@docker.com> (@moxiegirl)

RUN svn checkout https://github.com/docker/tutorials/trunk/docs/mac /docs/content/mac \
&& svn checkout https://github.com/docker/tutorials/trunk/docs/windows /docs/content/windows \
&& svn checkout https://github.com/docker/tutorials/trunk/docs/linux /docs/content/linux

RUN svn checkout https://github.com/docker/docker/trunk/docs /docs/content/ 

RUN svn checkout https://github.com/docker/compose/trunk/docs /docs/content/compose \
&& svn checkout https://github.com/docker/swarm/trunk/docs /docs/content/swarm  \
&& svn checkout https://github.com/docker/machine/trunk/docs /docs/content/machine \
&& svn checkout https://github.com/docker/distribution/trunk/docs /docs/content/registry \
&& svn checkout https://github.com/kitematic/kitematic/trunk/docs /docs/content/kitematic \
&& svn checkout https://github.com/docker/opensource/trunk/docs /docs/content/opensource 

COPY . /src

COPY . /docs/content/docker-trusted-registry/
 

RUN wget https://raw.githubusercontent.com/docker/docs.docker.com/master/touch-up.sh 
RUN chmod 777 touch-up.sh && mv touch-up.sh /src && cat /dev/null > /docs/build.json

# To get the git info for this repo

RUN /src/touch-up.sh /docs
