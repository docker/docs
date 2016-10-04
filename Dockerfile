FROM starefossen/ruby-node

RUN mkdir -p /docs
VOLUME /docs

EXPOSE 4000

WORKDIR /docs

RUN gem install github-pages

CMD jekyll clean && jekyll serve -H 0.0.0.0 -P 4000
