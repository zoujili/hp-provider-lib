pipeline {

    agent any

    options {
        buildDiscarder(logRotator(artifactDaysToKeepStr: '30', artifactNumToKeepStr: '15', daysToKeepStr: '150', numToKeepStr: '15'))
        disableConcurrentBuilds()
        ansiColor('xterm')
    }

    stages {

        stage('checkout') {
            steps {
                checkout scm
            }
        }

        stage("run code analysis") {
            agent {
                docker {
                    image 'quay.io/isaacfitstation/pims_microservice_code_analysis'
                    args '-e HOME=. -e TEST_JUNIT_REPORTER=true'
                    alwaysPull true
                    reuseNode true
                    registryUrl 'https://quay.io/user/isaacfitstation'
                    registryCredentialsId 'quay-isaacfitstation-jenkins'
                }
            }
            steps {
                sh 'rm -rf output'
                sh 'mkdir output'
                sh 'rm -rf .cache'
                sh 'go vet ./pkg/... > ./output/go_vet || true'
                sh 'rm -rf .cache'
                sh 'go test -v ./pkg/... || true'
                sh 'mv ./pkg/v1/test_*.xml ./output/ 2>/dev/null'
                sh 'golint ./pkg/... > ./output/go_lint'
                sh 'gocyclo -over 15 ./pkg > ./output/go_cyclo'
                sh 'ineffassign ./pkg > ./output/ineffassign || true'
            }
            post {
                always {
                    junit 'output/test_*.xml'
                    recordIssues(tools: [goLint(pattern: 'output/go_lint'), goVet(pattern: 'output/go_vet')])
                }
            }
        }
    }

    post {
        always {
            script {
                // build status of null means successful
                currentBuild.result = currentBuild.result ?: 'SUCCESS'

                // If unstable or result changed, notify the culprits and Slack
                if (currentBuild.resultIsWorseOrEqualTo('UNSTABLE') || (currentBuild.previousBuild == null || currentBuild.previousBuild.result != currentBuild.result)) {
                    emailext(
                            subject: "${currentBuild.result}: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]'",
                            body: '${JELLY_SCRIPT,template="html"}',
                            recipientProviders: [
                                    [$class: 'DevelopersRecipientProvider'],
                                    [$class: 'CulpritsRecipientProvider'],
                                    [$class: 'RequesterRecipientProvider']
                            ]
                    )
                }
            }
        }
    }
}
