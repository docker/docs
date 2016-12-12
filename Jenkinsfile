wrappedNode(label: 'linux && x86_64') {
  deleteDir()
  stage "checkout"
  checkout scm
  stage "test"
  sh "docker build -t docs `pwd`"
  sh "docker build -t tests `pwd`/tests"
  sh "docker run -v /usr/src/app/allvbuild --name docs docs /bin/true"
  sh "docker run --rm --volumes-from docs -v `pwd`:/docs tests"
  sh "docker rm -fv docs"
  sh "docker rmi docs tests"
}