pipeline {
    agent any

    stages {
        stage('Checkout Code') {
            steps {
                // Pulls the latest code from the main branch on GitHub
                checkout scm
            }
        }

        stage('Deploy to Docker') {
            steps {
                // Tears down the old container and rebuilds the monolith image cleanly
                sh 'docker-compose down'
                sh 'docker-compose up -d --build'
            }
        }
    }
}