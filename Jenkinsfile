def reg = [credentialsId: 'csebuildbot', url: 'https://index.docker.io/v1/']

pipeline {
  agent {
    label 'ubuntu-1604-aufs-stable'
  }
  environment {
    DTR_URL               = credentials('dtr-url')
    DOCKER_HOST_STRING    = credentials('docker-host')
    UCP_BUNDLE            = credentials('ucp-bundle')
    SUCCESS_BOT_TOKEN     = credentials('success-bot-token')
    SLACK                 = credentials('slack-docs-webhook')
  }
  options {
    timeout(time: 1, unit: 'HOURS') 
  }
  stages {
    stage( 'docker.github.io' ) {
      when {
        expression { env.GIT_URL == 'https://github.com/Docker/docker.github.io.git' }
      }
      stages {
        stage( 'build + push stage image, update stage swarm' ) {
          when {
            branch 'master'
          }
          steps {
            sh """
              cat $SUCCESS_BOT_TOKEN | docker login  $DTR_URL --username 'success_bot' --password-stdin
              docker build -t $DTR_URL/docker/docker.github.io:stage-${env.BUILD_NUMBER} .
              docker push $DTR_URL/docker/docker.github.io:stage-${env.BUILD_NUMBER}
              unzip -o $UCP_BUNDLE
              export DOCKER_TLS_VERIFY=1
              export COMPOSE_TLS_VERSION=TLSv1_2
              export DOCKER_CERT_PATH=${WORKSPACE}/ucp-bundle-success_bot
              export DOCKER_HOST=$DOCKER_HOST_STRING
              docker service update --detach=false --force --image $DTR_URL/docker/docker.github.io:stage-${env.BUILD_NUMBER} docs-stage-docker-com_docs --with-registry-auth
            """
          }
        }
        stage( 'build + push prod image, update prod swarm' ) {
          when {
            branch 'published'
          }
          steps {
            withDockerRegistry(reg) {
              sh """
                docker build -t docs/docker.github.io:prod-${env.BUILD_NUMBER} .
                docker tag docs/docker.github.io:prod-${env.BUILD_NUMBER} docs/docker.github.io:latest
                docker push docs/docker.github.io:prod-${env.BUILD_NUMBER}
                docker push docs/docker.github.io:latest
                unzip -o $UCP_BUNDLE
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
    stage( 'docs-private' ) {
      when {
        expression { env.GIT_URL == "https://github.com/docker/docs-private.git" }
      } 
      stages {
        stage( 'build + push beta-stage image, update beta-stage swarm' ) {
          when {
            branch 'amberjack'
          }
          steps {
            sh """
              cat $SUCCESS_BOT_TOKEN | docker login  $DTR_URL --username 'success_bot' --password-stdin
              docker build -t $DTR_URL/docker/docs-private:beta-stage-${env.BUILD_NUMBER} .
              docker push $DTR_URL/docker/docs-private:beta-stage-${env.BUILD_NUMBER}
              unzip -o $UCP_BUNDLE
              export DOCKER_TLS_VERIFY=1
              export COMPOSE_TLS_VERSION=TLSv1_2
              export DOCKER_CERT_PATH=${WORKSPACE}/ucp-bundle-success_bot
              export DOCKER_HOST=$DOCKER_HOST_STRING
              docker service update --detach=false --force --image $DTR_URL/docker/docs-private:beta-stage-${env.BUILD_NUMBER} docs-beta-stage-docker-com_docs --with-registry-auth
            """
          }
        }
        stage( 'build + push beta image, update beta swarm' ) {
          when {
            branch 'published'
          }
          steps {
            sh """
              cat $SUCCESS_BOT_TOKEN | docker login  $DTR_URL --username 'success_bot' --password-stdin
              docker build -t $DTR_URL/docker/docs-private:beta-${env.BUILD_NUMBER} .
              docker push $DTR_URL/docker/docs-private:beta-${env.BUILD_NUMBER}
              unzip -o $UCP_BUNDLE
              export DOCKER_TLS_VERIFY=1
              export COMPOSE_TLS_VERSION=TLSv1_2
              export DOCKER_CERT_PATH=${WORKSPACE}/ucp-bundle-success_bot
              export DOCKER_HOST=$DOCKER_HOST_STRING
              docker service update --detach=false --force --image $DTR_URL/docker/docs-private:beta-${env.BUILD_NUMBER} docs-beta-docker-com_docs --with-registry-auth
            """
          }
        }
      }
    }
  }
  post {
    unsuccessful {
      sh """
        curl -X POST -H 'Content-type: application/json' --data '{"text":"Error in Jenkins build. Please contact the Customer Success Engineering team for help."}' $SLACK
      """
    }
  }
}
