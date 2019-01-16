def dtrVpnAddress = "vpn.corp-us-east-1.aws.dckr.io"
def ucpBundle = [file(credentialsId: "ucp-bundle", variable: 'UCP')]

pipeline {
  agent none
  stages {
    stage( 'docker.github.io' ) {
      agent {
        label 'ubuntu-1604-aufs-stable'
      }
      stages {
        stage( 'build and push stage image' ) {
          when {
            branch 'master'
          }
          steps {
            withCredentials([usernamePassword(credentialsId: 'csebuildbot', passwordVariable: 'PWD', usernameVariable: 'USR')]) {
              sh """
                docker login -u ${USR} -p ${PWD} && \
                docker image build --tag docs/docker.github.io:stage-${env.BUILD_NUMBER} -f Dockerfile . && \
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
            withCredentials([usernamePassword(credentialsId: 'csebuildbot', passwordVariable: 'PWD', usernameVariable: 'USR')]) {
              sh """
                docker login -u ${USR} -p ${PWD} && \
                docker image build --tag docs/docker.github.io:prod-${env.BUILD_NUMBER} -f Dockerfile . && \
                docker image push docs/docker.github.io:prod-${env.BUILD_NUMBER}
              """
            }
          }
        }
        stage( 'update docs-stage' ) {
          when {
            branch 'master'
          }
          steps {
            withVpn(dtrVpnAddress) {
              withCredentials(ucpBundle) {
                sh 'unzip -o $UCP' 
              }
              withCredentials([usernamePassword(credentialsId: 'csebuildbot', passwordVariable: 'PWD', usernameVariable: 'USR')]) {
                sh """
                  cd ucp-bundle-success_bot
                  export DOCKER_TLS_VERIFY=1
                  export COMPOSE_TLS_VERSION=TLSv1_2
                  export DOCKER_CERT_PATH=${WORKSPACE}/ucp-bundle-success_bot
                  export DOCKER_HOST=tcp://ucp.corp-us-east-1.aws.dckr.io:443
                  docker login -u ${USR} -p ${PWD}
                  docker service update --detach=false --force --image docs/docker.github.io:stage-${env.BUILD_NUMBER} docs-stage-docker-com_docs --with-registry-auth
                """
              }
            }
          }
        }
        stage( 'update docs-prod' ) {
          when {
            branch 'published'
          }
          steps {
            withVpn(dtrVpnAddress) {
              withCredentials(ucpBundle) {
                sh 'unzip -o $UCP' 
              }
              withCredentials([usernamePassword(credentialsId: 'csebuildbot', passwordVariable: 'PWD', usernameVariable: 'USR')]) {
                sh """
                  cd ucp-bundle-success_bot
                  export DOCKER_TLS_VERIFY=1
                  export COMPOSE_TLS_VERSION=TLSv1_2
                  export DOCKER_CERT_PATH=${WORKSPACE}/ucp-bundle-success_bot
                  export DOCKER_HOST=tcp://ucp.corp-us-east-1.aws.dckr.io:443
                  docker login -u ${USR} -p ${PWD}
                  docker service update --detach=false --force --image docs/docker.github.io:prod-${env.BUILD_NUMBER} docs-docker-com_docs --with-registry-auth
                """
              }
            }
          }
        }
      }
    }
  }
}
