pipeline {
    agent any

    stages {
        stage('Build') {
            node {
            // Install the desired Go version
                def root = tool name: 'Go 1.8', type: 'go'

            // Export environment variables pointing to the directory where Go was installed
                withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
                    sh 'go version'
                    sh 'go env'
                    sh 'docker -v'
                }
            }
        }
        stage('Test'){
            node {
            // Install the desired Go version
                def root = tool name: 'Go 1.8', type: 'go'

            // Export environment variables pointing to the directory where Go was installed
                withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
                    sh 'go version'
                    sh 'docker -v'
                }
            }
            steps {
                sh 'go version'
                sh 'go env'
            }
        }
        stage('Deploy') {
            steps {
                sh 'docker -v'
            }
        }
    }
}
