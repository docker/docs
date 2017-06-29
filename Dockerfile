FROM starefossen/github-pages:137

COPY . /usr/src/app/

CMD jekyll serve -H 0.0.0.0