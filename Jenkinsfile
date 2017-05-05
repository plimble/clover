pipeline {
    agent any
    node {
    // Install the desired Go version
        def root = tool name: 'Go 1.8', type: 'go'

    // Export environment variables pointing to the directory where Go was installed
        withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
            sh 'go version'
        }
    }

    stages {
        stage('Build') {
            steps {
                sh 'docker -v'
            }
        }
        stage('Test'){
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
