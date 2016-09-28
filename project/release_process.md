# Releasing images for UCP

In general, our approach is to codify most of the smarts in our build
files within the project tree (Makefiles, scripts, etc.) and keep the CI
system's jobs code fairly "simple."  When the build is run locally without
specfying overrides to various environment variables, the assumption
is you are in developer mode.  We try to make developer mode look as
close as reasonably possible to what customers see, so the default "org"
for the images you build as a developer are "docker"

We use two hub orgs: **dockerorcadev** and **docker**

The general flow looks like the following

1. Developer submits PR
    * CI builds it
        * https://jenkins-orca.dockerproject.com/job/orca-pr/
    * To test a PR:

            for i in $(docker run --rm dockerorcadev/ucp:${VERSION}-pr${PR_NUMBER} images --list dev: ) ; do docker pull $i; done
            docker run --rm --name ucp \
                -v /var/run/docker.sock:/var/run/docker.sock \
                dockerorcadev/ucp:${VERSION}-pr${PR_NUMBER} \
                install --image-version dev: [arguments...]

2. PR is merged
3. CI Builds
    * https://jenkins-orca.dockerproject.com/job/orca-master/ triggers https://jenkins-orca.dockerproject.com/job/orca-build/
    * These builds push to `dockerorcadev/ucp*:${VERSION}-ob${BUILD_NUMBER}`
        * `${VERSION}` comes from `./version/version.go:Version`
        * `${BUILD_NUMBER}` is the jenkins build number
4. Release captain announces the intention to release and tells people to start their final testing before promoting the images to the public repositories
    * Identify the version tag to test.  E.g., `1.1.0-ob498` and report the Git SHA

            docker run --rm dockerorcadev/ucp:${TAG} --version

5. Testing occurs using the `dockerorcadev` builds identified by the release captain above.

        for i in $(docker run --rm dockerorcadev/ucp:${TAG} images --list --image-version dev:${TAG} ) ; do docker pull $i; done
        docker run --rm -it \
            --name ucp \
            -v /var/run/docker.sock:/var/run/docker.sock \
            dockerorcadev/ucp:${TAG} \
            install --image-version dev:${TAG} [arguments...]

6. Promote as an externally visible RC
    * Requires the exact SHA identified above, and the target RC number as the `PRE_RELEASE` (e.g., "rc1", "rc2", etc.)
    * Upon completion of this step, customers will be able to access the new images.  E.g., "1.1.0-rc1"
    * **Note: With our current CI system, the images will be re-built from source, but based on the same git commit sha**
        * This may be an area we want to change in the future to reduce the small chance of drift in underlying content (OS patches, etc.) but in the time gap should typically be very short.
    * Use a specific Git commit **NOT HEAD** to ensure the version you want actually gets published.
        * If no changes have been merged to master, this will be the same, but it might not always be.
    * https://jenkins-orca.dockerproject.com/job/orca-promote
6. Promote as an externally visible GA (non-RC build)
    * Same steps as above, using the same git commit SHA, however no `PRE_RELEASE` string is specified (left blank)
7. At this point, the GA image is live for customers, but is not identified as "latest"
    * Unless customers explicitly target this new version, they'll continue to get the prior release.
8. Perform final sanity test using the explicit version tag with the "docker" org to ensure nothing looks out of place
    * Verify the sha is as expected
    * Sanity test product - nothing catastrophic broke
9. Run the final CI job to tag the bootstrapper image as latest to customers get it by default
    * This build just updates the "latest" tag of the existing hub images - you must specify the version tag to mark latest
    * https://jenkins-orca.dockerproject.com/job/orca-bootstrap-update-latest/
10. Generate combined tar file and upload to S3

        UCP_VERSION=1.1.2
        DTR_VERSION=2.0.2
        for i in $(docker run docker/dtr:${DTR_VERSION} images) ; do docker pull $i; done
        for i in $(docker run docker/ucp:${UCP_VERSION} images --list) ; do docker pull $i; done
        docker save \
            $(docker run docker/dtr:${DTR_VERSION} images) \
            docker/dtr:${DTR_VERSION} \
            $(docker run docker/ucp:${UCP_VERSION} images --list) \
            docker/ucp:${UCP_VERSION} |gzip -c9  > ucp-${UCP_VERSION}_dtr-${DTR_VERSION}.tar.gz

