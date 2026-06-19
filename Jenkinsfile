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
                sh 'docker compose down'
                sh 'docker compose up -d --build'
            }
        }
    }
    
    post {
        success {
            withCredentials([string(credentialsId: 'slack-webhook', variable: 'SLACK_URL')]) {
                sh "curl -X POST -H 'Content-type: application/json' --data '{\"text\":\"✅ *SUCCESS:* Announcement Board deployed perfectly to production!\"}' \$SLACK_URL"
            }
        }
        failure {
            withCredentials([string(credentialsId: 'slack-webhook', variable: 'SLACK_URL')]) {
                sh "curl -X POST -H 'Content-type: application/json' --data '{\"text\":\"❌ *FAILED:* The pipeline crashed. Check Jenkins logs.\"}' \$SLACK_URL"
            }
        }
    }
}