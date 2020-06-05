pipeline {

 agent {label 'dockerbuild'}

 options {
  buildDiscarder(logRotator(artifactDaysToKeepStr: '30', artifactNumToKeepStr: '15', daysToKeepStr: '150', numToKeepStr: '15'))
  disableConcurrentBuilds()
  ansiColor('xterm')
  timeout(time: 30, unit: 'MINUTES')
 }

 stages {

  stage('checkout') {
   steps {
    checkout scm
   }
  }

  stage("Build") {
   agent {
    docker {
     image 'golang:1.14.2'
     args '-u root'
    }
   }
   steps {
    withCredentials([usernamePassword(credentialsId: 'hpbp-robot', passwordVariable: 'PASS', usernameVariable: 'USER')]) {
     sh "echo machine github.azc.ext.hp.com login ${USER} password ${PASS} > /root/.netrc"
   }
    sh 'go mod tidy && go get github.com/onsi/ginkgo/ginkgo && make test'
   }
  }
 }

 post {
  always {
   script {
    // build status of null means successful
    currentBuild.result = currentBuild.result ?: 'SUCCESS'

    // If unstable or result changed, notify
    if (currentBuild.resultIsWorseOrEqualTo('UNSTABLE') || (currentBuild.previousBuild == null || currentBuild.previousBuild.result != currentBuild.result)) {
     emailext(
      subject: "${currentBuild.result}: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]'",
      body: "${env.BUILD_URL}",
      attachLog: true,
      recipientProviders: [
       [$class: 'RequesterRecipientProvider']
      ],
      to: 'yong.jiang1@hp.com'
     )
    }
   }
  }
 }
}
