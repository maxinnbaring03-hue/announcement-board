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
                sh '''
                echo '{"text":"✅ *SUCCESS:* Announcement Board deployed perfectly to production!"}' > payload.json
                curl -X POST -H "Content-Type: application/json" -d @payload.json $SLACK_URL
                '''
            }
        }
        failure {
            withCredentials([string(credentialsId: 'slack-webhook', variable: 'SLACK_URL')]) {
                sh '''
                echo '{"text":"❌ *FAILED:* The pipeline crashed. Check Jenkins logs."}' > payload.json
                curl -X POST -H "Content-Type: application/json" -d @payload.json $SLACK_URL
                '''
            }
        }
    }
}