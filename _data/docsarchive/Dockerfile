FROM starefossen/github-pages
# The branch name to check out, also used in some find/replace
# This is just the default, specify at build-time using --build-arg VER=<foo>
ARG VER="v1.11"

# Get the branch contents without history
RUN echo "Started $VER" \
    && git clone --depth 1 --branch $VER https://github.com/docker/docker.github.io docs_src \
    && jekyll build -s docs_src -d docs \
    && rm -rf docs_src \
    && find docs -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/'"$VER"'/#g'


CMD jekyll serve -b $VER -s docs --no-watch -H 0.0.0.0 -P 4000
