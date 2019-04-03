def reg = [credentialsId: 'csebuildbot', url: 'https://index.docker.io/v1/']

pipeline {
  agent none
  options {
    timeout(time: 1, unit: 'HOURS') 
  }
  stages {
    stage( 'docker.github.io' ) {
      agent {
        label 'ubuntu-1604-aufs-stable'
      }
      environment {
        DTR_VPN_ADDRESS       = credentials('dtr-vpn-address')
        DOCKER_HOST_STRING    = credentials('docker-host')
        UCP_BUNDLE            = credentials('ucp-bundle')
        SLACK                 = credentials('slack-docs-webhook')
      }
      when {
        expression { env.GIT_URL == 'https://github.com/Docker/docker.github.io.git' }
      }
      stages {
        stage( 'build and push stage image' ) {
          when {
            branch 'master'
          }
          steps {
            withDockerRegistry(reg) {
              sh """
                docker image build --tag docs/docker.github.io:stage-${env.BUILD_NUMBER} . && \
                docker image push docs/docker.github.io:stage-${env.BUILD_NUMBER}
              """
            }
          }
        }
        stage( 'build and push prod image' ) {
          when {
            branch 'published'
          }
          steps {
            withDockerRegistry(reg) {
              sh """
                docker image build --tag docs/docker.github.io:prod-${env.BUILD_NUMBER} . && \
                docker image push docs/docker.github.io:prod-${env.BUILD_NUMBER}
              """
            }
          }
        }
        stage( 'update docs stage' ) {
          when {
            branch 'master'
          }
          steps {
            withVpn("$DTR_VPN_ADDRESS") {
              sh "unzip -o $UCP_BUNDLE"
              withDockerRegistry(reg) {
                sh """
                  export DOCKER_TLS_VERIFY=1
                  export COMPOSE_TLS_VERSION=TLSv1_2
                  export DOCKER_CERT_PATH=${WORKSPACE}/ucp-bundle-success_bot
                  export DOCKER_HOST=$DOCKER_HOST_STRING
                  docker service update --detach=false --force --image docs/docker.github.io:stage-${env.BUILD_NUMBER} docs-stage-docker-com_docs --with-registry-auth
                """
              }
            }
          }
        }
        stage( 'update docs prod' ) {
          when {
            branch 'published'
          }
          steps {
            withVpn("$DTR_VPN_ADDRESS") {
              sh "unzip -o $UCP_BUNDLE"
              withDockerRegistry(reg) {
                sh """
                  cd ucp-bundle-success_bot
                  export DOCKER_TLS_VERIFY=1
                  export COMPOSE_TLS_VERSION=TLSv1_2
                  export DOCKER_CERT_PATH=${WORKSPACE}/ucp-bundle-success_bot
                  export DOCKER_HOST=$DOCKER_HOST_STRING
                  docker service update --detach=false --force --image docs/docker.github.io:prod-${env.BUILD_NUMBER} docs-docker-com_docs --with-registry-auth
                  curl -X POST -H 'Content-type: application/json' --data '{"text":"Successfully published docs. https://docs.docker.com/"}' $SLACK
                """
              }
            }
          }
        }
      }
    }
    stage( 'docs-private' ) {
      agent {
        label 'ubuntu-1604-aufs-stable'
      }
      environment {
        DTR_VPN_ADDRESS       = credentials('dtr-vpn-address')
        DOCKER_HOST_STRING    = credentials('docker-host')
        UCP_BUNDLE            = credentials('ucp-bundle')
      }
      when {
        expression { env.GIT_URL == "https://github.com/docker/docs-private.git" }
      } 
      stages {
        stage( 'build and push new beta stage image' ) {
          when {
            branch 'amberjack'
          }
          steps {
            withDockerRegistry(reg) {
              sh """
                docker image build --tag docs/docs-private:beta-stage-${env.BUILD_NUMBER} . && \
                docker image push docs/docs-private:beta-stage-${env.BUILD_NUMBER}
              """
            }
          }
        }
        stage( 'build and push new beta image' ) {
          when {
            branch 'published'
          }
          steps {
            withDockerRegistry(reg) {
              sh """
                docker image build --tag docs/docs-private:beta-${env.BUILD_NUMBER} . && \
                docker image push docs/docs-private:beta-${env.BUILD_NUMBER}
              """
            }
          }
        }
        stage( 'update beta stage service' ) {
          when {
            branch 'amberjack'
          }
          steps {
            withVpn("$DTR_VPN_ADDRESS") {
              sh "unzip -o $UCP_BUNDLE"
              withDockerRegistry(reg) {
                sh """
                  export DOCKER_TLS_VERIFY=1
                  export COMPOSE_TLS_VERSION=TLSv1_2
                  export DOCKER_CERT_PATH=${WORKSPACE}/ucp-bundle-success_bot
                  export DOCKER_HOST=$DOCKER_HOST_STRING
                  docker service update --detach=false --force --image docs/docs-private:beta-stage-${env.BUILD_NUMBER} docs-beta-stage-docker-com_docs --with-registry-auth
                """
              }
            }
          }
        }
        stage( 'update beta service' ) {
          when {
            branch 'published'
          }
          steps {
            withVpn("$DTR_VPN_ADDRESS") {
              sh "unzip -o $UCP_BUNDLE"
              withDockerRegistry(reg) {
                sh """
                  export DOCKER_TLS_VERIFY=1
                  export COMPOSE_TLS_VERSION=TLSv1_2
                  export DOCKER_CERT_PATH=${WORKSPACE}/ucp-bundle-success_bot
                  export DOCKER_HOST=$DOCKER_HOST_STRING
                  docker service update --detach=false --force --image docs/docs-private:beta-${env.BUILD_NUMBER} docs-beta-docker-com_docs --with-registry-auth
                """
              }
            }
          }
        }
      }
    }
  }
}
