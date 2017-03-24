wrappedNode(label: 'linux && x86_64') {
  timeout(time: 60, unit: 'MINUTES') {
    /* Start with a clean workspace */
    deleteDir()

    /* Check out the branch or PR */
    stage "checkout"
    checkout scm

    /* Start the testing stage */
    stage "test"

    /* Build a container from the current checkout */
    sh "docker build -t docs:${JOB_BASE_NAME}-${BUILD_NUMBER} `pwd`"

    /* Start the container you just built, with two volumes: one for the source and the other for the HTML output */
    sh "docker run -v /usr/src/app/allvbuild -v /usr/src/app/target --name docs-${JOB_BASE_NAME}-${BUILD_NUMBER} docs:${JOB_BASE_NAME}-${BUILD_NUMBER} /bin/true"

    /* Start a linkchecker container using the volumes from the docs container and set it to check the HTML output for broken links */
    sh "docker run --rm --volumes-from docs-${JOB_BASE_NAME}-${BUILD_NUMBER} --env site=/usr/src/app/target --name linkchecker-${JOB_BASE_NAME}-${BUILD_NUMBER} docs/tools:linkchecker"

    /* Remove the containers, including volumes */
    sh "docker rm -fv docs-${JOB_BASE_NAME}-${BUILD_NUMBER} linkchecker-${JOB_BASE_NAME}-${BUILD_NUMBER}"

    /* Remove the image */
    sh "docker rmi docs:${JOB_BASE_NAME}-${BUILD_NUMBER}"

    /* Clean up the workspace */
    deleteDir()
  }
}
