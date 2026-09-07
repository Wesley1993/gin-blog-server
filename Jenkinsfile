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
        DEPLOY_USER = 'root'
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

        stage('Deploy') {
            steps {
                sh '''
                    # Ship binary
                    rsync -az --progress bin/blog-server \
                      $DEPLOY_USER@$DEPLOY_HOST:/opt/zblog/blog-server.new

                    # Ship migrations
                    rsync -az --delete migrations/ \
                      $DEPLOY_USER@$DEPLOY_HOST:/opt/zblog/migrations/

                    # Sync frontend assets
                    rsync -az --delete blog/dist/ \
                      $DEPLOY_USER@$DEPLOY_HOST:/opt/zblog/www/blog/
                    rsync -az --delete web/dist/ \
                      $DEPLOY_USER@$DEPLOY_HOST:/opt/zblog/www/admin/

                    # Sync infra files
                    rsync -az docker-compose.infra.yml \
                      $DEPLOY_USER@$DEPLOY_HOST:/opt/zblog/docker-compose.infra.yml
                    rsync -az nginx/conf.d/ \
                      $DEPLOY_USER@$DEPLOY_HOST:/opt/zblog/nginx/conf.d/

                    # Activate
                    ssh $DEPLOY_USER@$DEPLOY_HOST "
                      set -e
                      cd /opt/zblog
                      # Pre-deploy backup
                      pg_dump -h 127.0.0.1 -U \\${POSTGRES_USER:-appuser} -Fc \\${POSTGRES_DB:-blog} > backup/pre_deploy.dump 2>/dev/null || true
                      # Swap binary
                      [ -f blog-server ] && cp blog-server blog-server.old
                      mv blog-server.new blog-server
                      chmod +x blog-server
                      # Restart
                      sudo systemctl restart zblog
                      # Reload nginx
                      docker compose -f docker-compose.infra.yml exec nginx nginx -s reload || true
                      # Health check
                      sleep 5
                      curl -fsS http://127.0.0.1:3000/health || exit 1
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
                  if [ -f blog-server.old ]; then
                    mv blog-server.old blog-server
                    sudo systemctl restart zblog
                    echo 'Rollback completed'
                  else
                    echo 'No backup binary found, manual intervention required'
                  fi
                " || echo "Rollback SSH failed - manual intervention required"
            '''
        }
        always {
            sh 'rm -rf bin/'
            cleanWs(deleteDirs: true, notFailBuild: true)
        }
    }
}
