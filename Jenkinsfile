pipeline {
    agent any

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
