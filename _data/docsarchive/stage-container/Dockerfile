FROM starefossen/ruby-node

RUN mkdir -p /docs
VOLUME /docs

EXPOSE 4000

WORKDIR /docs

RUN gem install github-pages

CMD jekyll clean && jekyll serve -H 0.0.0.0 -P 4000

# To run:
#
# git clone https://github.com/docker/docker.github.io.git docs
# cd docs
# docker run -ti --rm -v "$PWD":/docs -p 4000:4000 docs/docstage
