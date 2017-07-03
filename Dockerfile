FROM starefossen/github-pages:137

COPY . /usr/src/app

CMD bundle exec jekyll serve -d /_site --watch -H 0.0.0.0 -P 4000