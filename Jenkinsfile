wrappedNode(label: 'ubuntu-1604 && x86_64') {
  timeout(time: 60, unit: 'MINUTES') {
    deleteDir()
    stage "checkout"
    checkout scm
    sh "git submodule update --init --recursive"
    stage "test"

    /* Jekyll creates html files to implement client side redirects.
      There are absolute links to docs.docker.com in these htmls
      we don't want them to be parsed by the tests for now.
      Removing jekyll-redirect-from option will make sure these pages
      are not generated when building with Jekyll. */
    sh "awk '/jekyll-redirect-from/{n=1}; n {n--; next}; 1' < _config.yml > _config.yml.tmp"
    sh "mv _config.yml.tmp _config.yml"

    sh "docker build -t docs:${JOB_BASE_NAME}-${BUILD_NUMBER} `pwd`"
    sh "docker build -t tests:${JOB_BASE_NAME}-${BUILD_NUMBER} `pwd`/tests"
    sh "docker run -v /usr/src/app/allvbuild --name docs-${JOB_BASE_NAME}-${BUILD_NUMBER} docs:${JOB_BASE_NAME}-${BUILD_NUMBER} /bin/true"
    sh "docker run --rm --volumes-from docs-${JOB_BASE_NAME}-${BUILD_NUMBER} -v `pwd`:/docs tests:${JOB_BASE_NAME}-${BUILD_NUMBER}"
    sh "docker rm -fv docs-${JOB_BASE_NAME}-${BUILD_NUMBER}"
    sh "docker rmi docs:${JOB_BASE_NAME}-${BUILD_NUMBER} tests:${JOB_BASE_NAME}-${BUILD_NUMBER}"
    deleteDir()
  }
}
