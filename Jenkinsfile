pipeline {
    agent any

    options {
        disableConcurrentBuilds()
        buildDiscarder(logRotator(numToKeepStr: '20'))
        timeout(time: 15, unit: 'MINUTES')
        timestamps()
    }

    environment {
        DEPLOY_HOST = credentials('zblog-deploy-host')
        DEPLOY_USER = 'deploy'
        TAG         = "${env.BUILD_NUMBER}-${env.GIT_COMMIT?.take(7) ?: 'unknown'}"
    }

    triggers {
        GenericTrigger(
            genericVariables: [[key: 'ref', value: '$.ref']],
            token: 'zblog-webhook',
            causeString: 'Git push: $ref',
            printContributedVariables: true,
            printPostContent: false
        )
    }

    stages {
        stage('Quality Gate') {
            parallel {
                stage('Go Vet') {
                    steps {
                        sh 'go vet ./...'
                    }
                }
                stage('Web Lint & TypeCheck') {
                    steps {
                        dir('web') {
                            sh 'npm ci --silent'
                            sh 'npx oxlint'
                            sh 'npx tsc --noEmit'
                        }
                    }
                }
                stage('Blog TypeCheck') {
                    steps {
                        dir('blog') {
                            sh 'npm ci --silent'
                            sh 'npx tsc --noEmit'
                        }
                    }
                }
            }
        }

        stage('Build') {
            parallel {
                stage('Backend') {
                    steps {
                        sh '''
                            CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
                            go build -trimpath -ldflags="-s -w" \
                            -o bin/blog-server ./cmd
                        '''
                    }
                }
                stage('Blog Frontend') {
                    steps {
                        dir('blog') {
                            sh 'npm run build'
                        }
                    }
                }
                stage('Admin Frontend') {
                    steps {
                        dir('web') {
                            sh 'npm run build'
                        }
                    }
                }
            }
        }

        stage('Package') {
            steps {
                sh '''
                    docker build --platform linux/amd64 \
                      -f deploy/Dockerfile.api \
                      -t zblog-api:$TAG .
                '''
            }
        }

        stage('Deploy') {
            steps {
                sh '''
                    # Ship image (no registry needed)
                    docker save zblog-api:$TAG | gzip | \
                      ssh -o StrictHostKeyChecking=accept-new \
                          -o ServerAliveInterval=15 \
                          $DEPLOY_USER@$DEPLOY_HOST "gunzip | docker load"

                    # Sync frontend assets
                    rsync -az --delete blog/dist/ \
                      $DEPLOY_USER@$DEPLOY_HOST:/opt/zblog/www/blog/
                    rsync -az --delete web/dist/ \
                      $DEPLOY_USER@$DEPLOY_HOST:/opt/zblog/www/admin/

                    # Activate
                    ssh $DEPLOY_USER@$DEPLOY_HOST "
                      set -e
                      cd /opt/zblog
                      # Pre-deploy DB backup
                      docker compose exec -T postgres pg_dump -U \${POSTGRES_USER:-appuser} -Fc \${POSTGRES_DB:-blog} > backup/pre_deploy_\\$(date +%F_%H%M).dump 2>/dev/null || true
                      # Switch tag
                      sed -i.bak \\"s/^TAG=.*/TAG=$TAG/\\" .env
                      # Restart API
                      docker compose up -d api
                      # Reload Nginx
                      docker compose exec nginx nginx -t 2>/dev/null
                      docker compose exec nginx nginx -s reload
                      # Health check
                      sleep 10
                      docker compose exec -T api wget -qO- http://127.0.0.1:3000/health || exit 1
                      curl -fsS http://127.0.0.1/health || exit 1
                    "
                '''
            }
        }
    }

    post {
        failure {
            sh '''
                echo "Deployment failed - attempting rollback..."
                ssh -o StrictHostKeyChecking=accept-new $DEPLOY_USER@$DEPLOY_HOST "
                  cd /opt/zblog
                  if [ -f .env.bak ]; then
                    cp .env.bak .env
                    docker compose up -d api
                    echo 'Rollback completed'
                  else
                    echo 'No backup found, manual intervention required'
                  fi
                " || echo "Rollback SSH failed - manual intervention required"
            '''
        }
        always {
            sh 'docker rmi zblog-api:$TAG 2>/dev/null || true'
            sh 'rm -rf bin/'
            cleanWs(deleteDirs: true, notFailBuild: true)
        }
    }
}
