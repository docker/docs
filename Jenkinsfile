wrappedNode(label: 'linux && x86_64') {
  deleteDir()
  stage "checkout"
  checkout scm
  stage "test"
  sh "ls -al"
  bash hooks/pre_build
}