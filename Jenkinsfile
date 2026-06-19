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
            // Natively writes the file so Linux quotes can't break it
            writeFile file: 'success.json', text: '{"text":"✅ *SUCCESS:* Announcement Board deployed perfectly to production!"}'
            withCredentials([string(credentialsId: 'slack-webhook', variable: 'SLACK_URL')]) {
                // The -g flag tells curl to ignore all brackets
                sh 'curl -g -X POST -H "Content-Type: application/json" -d @success.json "$SLACK_URL"'
            }
        }
        failure {
            writeFile file: 'fail.json', text: '{"text":"❌ *FAILED:* The pipeline crashed. Check Jenkins logs."}'
            withCredentials([string(credentialsId: 'slack-webhook', variable: 'SLACK_URL')]) {
                sh 'curl -g -X POST -H "Content-Type: application/json" -d @fail.json "$SLACK_URL"'
            }
        }
    }
}