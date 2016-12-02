wrappedNode(label: 'linux && x86_64') {
  deleteDir()
  stage "checkout"
  checkout scm
  stage "test"
  sh "chmod +x hooks/pre_build"
  sh "hooks/pre_build"
}