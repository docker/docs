FROM starefossen/ruby-node:2-5

RUN gem install github-pages

WORKDIR /data

RUN git clone https://github.com/docker/docker.github.io docs

WORKDIR /data/docs

RUN git checkout v1.11

RUN jekyll build

ENTRYPOINT ["jekyll", "serve", "--host=0.0.0.0"]
