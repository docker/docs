FROM starefossen/github-pages:onbuild

CMD bundle update && bundle exec jekyll serve -d /_site --watch -H 0.0.0.0 -P 4000
