// Only run on Linux atm
wrappedNode(label: 'linux') {
  deleteDir()
  stage "checkout"
  checkout scm

  documentationChecker("docs")
}
