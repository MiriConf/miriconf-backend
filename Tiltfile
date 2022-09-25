load('ext://git_resource', 'git_checkout')
git_checkout('git@github.com:MiriConf/miriconf-frontend.git#main', '../miriconf-frontend')

include('../miriconf-frontend/Tiltfile')

docker_build('miriconf-backend', '.', dockerfile='dev-deploy/Dockerfile')
k8s_yaml('dev-deploy/deployment.yaml')
k8s_resource('miriconf-backend', port_forwards=8081)

load('ext://helm_remote', 'helm_remote')
helm_remote('mongodb',
    repo_name='bitnami',
    repo_url='https://charts.bitnami.com/bitnami',
    values=["./dev-deploy/mongo-values.yaml"]
)