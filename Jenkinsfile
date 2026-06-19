pipeline {
    agent any

    stages {
        stage('Checkout Code') {
            steps {
                checkout scm
            }
        }

        stage('Deploy to Docker') {
            steps {
                // THE HYPHENS ARE GONE
                sh 'docker compose down'
                sh 'docker compose up -d --build'
            }
        }
    }
}