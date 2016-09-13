FROM starefossen/ruby-node:2-5

RUN gem install github-pages

VOLUME /docs

EXPOSE 4000

WORKDIR /docs

CMD jekyll clean && jekyll serve -H 0.0.0.0 -P 4000

WORKDIR /data

RUN git clone https://github.com/docker/docker.github.io.git docs
RUN git checkout v1.4

WORKDIR /data/docs

ENTRYPOINT ["jekyll", "serve", "--host=0.0.0.0"]
