
FROM starefossen/github-pages:onbuild

ONBUILD RUN git clone https://www.github.com/docker/docker.github.io docs

ONBUILD WORKDIR docs

ONBUILD RUN git checkout vnext-kitematic

ONBUILD COPY . /usr/src/app

CMD jekyll serve -d /_site -H 0.0.0.0 -P 4000
