def dtrVpnAddress = "vpn.corp-us-east-1.aws.dckr.io"
def ucpBundle = [file(credentialsId: "ucp-bundle", variable: 'UCP')]
def slackString = [string(credentialsId: 'slack-docs-webhook', variable: 'slack')]
def reg = [credentialsId: 'csebuildbot', url: 'https://index.docker.io/v1/']

pipeline {
  agent none
  options {
    timeout(time: 1, unit: 'HOURS') 
  }
  stages {
    stage( 'docs-private' ) {
      agent {
        label 'ubuntu-1604-aufs-stable'
      }
      when {
        expression { env.GIT_URL == 'https://github.com/docker/docs-private.git' }
      }
      stages {
        stage( 'build and push new beta-stage image' ) {
          when {
            branch 'jenkins-test'
          }
          steps {
            sh 'echo "Build and push new beta-stage image here"'
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
            sh 'echo "Build and push new beta image here"'
            withDockerRegistry(reg) {
              sh """
                docker image build --tag docs/docs-private:beta-${env.BUILD_NUMBER} . && \
                docker image push docs/docs-private:beta-${env.BUILD_NUMBER}
              """
            }
          }
        }
        stage( 'update beta-stage service' ) {
          when {
            branch 'jenkins-test'
          }
          steps {
            sh 'echo "Update beta-stage service here"'
            withVpn(dtrVpnAddress) {
              withCredentials(ucpBundle) {
                sh 'unzip -o $UCP' 
              }
              withDockerRegistry(reg) {
                sh """
                  export DOCKER_TLS_VERIFY=1
                  export COMPOSE_TLS_VERSION=TLSv1_2
                  export DOCKER_CERT_PATH=${WORKSPACE}/ucp-bundle-success_bot
                  export DOCKER_HOST=tcp://ucp.corp-us-east-1.aws.dckr.io:443
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
            sh 'echo "Update beta service here"'
            withVpn(dtrVpnAddress) {
              withCredentials(ucpBundle) {
                sh 'unzip -o $UCP' 
              }
              withDockerRegistry(reg) {
                sh """
                  export DOCKER_TLS_VERIFY=1
                  export COMPOSE_TLS_VERSION=TLSv1_2
                  export DOCKER_CERT_PATH=${WORKSPACE}/ucp-bundle-success_bot
                  export DOCKER_HOST=tcp://ucp.corp-us-east-1.aws.dckr.io:443
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
