wrappedNode(label: 'linux && x86_64') {
  deleteDir()
  stage "checkout"
  checkout scm
  stage "test"
  sh "hooks/pre_build"
}