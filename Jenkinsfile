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
                sh 'docker rm -f production-board || true'
                sh 'docker compose down'
                sh 'docker compose up -d --build'
            }
        }
    }
    
    post {
        success {
            withCredentials([string(credentialsId: 'slack-webhook', variable: 'SLACK_URL')]) {
                sh '''
                curl -g -X POST -H "Content-Type: application/json" -d '{"text":"✅ SUCCESS: Announcement Board deployed perfectly to production!"}' "$SLACK_URL"
                '''
            }
        }
        failure {
            withCredentials([string(credentialsId: 'slack-webhook', variable: 'SLACK_URL')]) {
                sh '''
                curl -g -X POST -H "Content-Type: application/json" -d '{"text":"❌ FAILED: The pipeline crashed. Check Jenkins logs."}' "$SLACK_URL"
                '''
            }
        }
    }
}