FROM starefossen/ruby-node:2-4

RUN git clone https://www.github.com/docker/docker.github.io docs

RUN gem install --no-document github-pages jekyll-github-metadata

EXPOSE 4000

EXPOSE 4000

CMD jekyll serve --source docs -d docs/_site -H 0.0.0.0 -P 4000
