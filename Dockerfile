FROM starefossen/github-pages:onbuild

ONBUILD RUN git clone https://www.github.com/docker/docker.github.io docs

ONBUILD WORKDIR docs

ONBUILD COPY . /usr/src/app
