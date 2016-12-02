wrappedNode(label: 'linux && x86_64') {
  deleteDir()
  stage "checkout"
  checkout scm
  stage "test"
  sh "docker build -t tests `pwd`/tests"
  sh "docker run --rm -v `pwd`:/docs tests"
  sh "docker rmi tests"
}