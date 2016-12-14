wrappedNode(label: 'linux && x86_64') {
  deleteDir()
  stage "checkout"
  checkout scm
  stage "test"

  # Jekyll creates html files to implement client side redirects.
  # There are absolute links to docs.docker.com in these htmls
  # we don't want them to be parsed by the tests for now.
  # Removing jekyll-redirect-from option will make sure these pages
  # are not generated when building with Jekyll.
  sh "awk '/jekyll-redirect-from/{n=1}; n {n--; next}; 1' < _config.yml > _config.yml.tmp"
  sh "mv _config.yml.tmp _config.yml"
  
  sh "docker build -t docs `pwd`"
  sh "docker build -t tests `pwd`/tests"
  sh "docker run -v /usr/src/app/allvbuild --name docs docs /bin/true"
  sh "docker run --rm --volumes-from docs -v `pwd`:/docs tests"
  sh "docker rm -fv docs"
  sh "docker rmi docs tests"
}