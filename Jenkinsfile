wrappedNode(label: 'linux && x86_64') {
  deleteDir()
  stage "checkout"
  checkout scm
  sh "git submodule update --init --recursive"
  stage "test"

  def tag = "${JOB_BASE_NAME}-${BUILD_NUMBER}"

  sh "docker build -t docs-with-redirects:${tag} `pwd`"

  /* Jekyll creates html files to implement client side redirects.
    There are absolute links to docs.docker.com in these htmls
    we don't want them to be parsed by the tests for now.
    Removing jekyll-redirect-from option will make sure these pages
    are not generated when building with Jekyll. */
  sh "awk '/jekyll-redirect-from/{n=1}; n {n--; next}; 1' < _config.yml > _config.yml.tmp"
  sh "mv _config.yml.tmp _config.yml"
  sh "docker build -t docs-without-redirects:${tag} `pwd`"

  /* run both containers (with and without redirects), 
    exposing /usr/src/app/allvbuild as a volume */
  sh "docker run -v /usr/src/app/allvbuild --name docs-with-redirects-${tag} docs-with-redirects:${tag} /bin/true"
  sh "docker run -v /usr/src/app/allvbuild --name docs-without-redirects-${tag} docs-without-redirects:${tag} /bin/true"

  def filesWithRedirects = sh (
    script: "docker inspect docs-with-redirects-${tag} | grep \"Source\" | awk 'match(\$2,/\"([^\"]+)\"/) {print substr(\$2,RSTART+1,RLENGTH-2)}'",
    returnStdout: true
  ).trim()

  def filesWithoutRedirects = sh (
    script: "docker inspect docs-without-redirects-${tag} | grep \"Source\" | awk 'match(\$2,/\"([^\"]+)\"/) {print substr(\$2,RSTART+1,RLENGTH-2)}'",
    returnStdout: true
  ).trim()

  sh "docker build -t tests:${tag} `pwd`/tests"
    
  sh "docker run --rm -v ${filesWithRedirects}:/docs-with-redirects -v ${filesWithoutRedirects}:/docs-without-redirects -v `pwd`:/docs-source tests:${tag}"

  sh "docker rm -fv docs-with-redirects-${tag} docs-without-redirects-${tag}"
  sh "docker rmi docs-with-redirects:${tag} docs-without-redirects:${tag} tests:${tag}"
}