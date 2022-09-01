# miriconf-backend

Backend api for miriconf installations, it is using the [Golang Gin Web Framework](https://gin-gonic.com) to assemble the api. This repo will automatically build and test the go binary, package it in a container, and push the container to [Github Packages](https://github.com/orgs/MiriConf/packages?repo_name=miriconf-backend). For our API documentation we are using the [swaggo](https://github.com/swaggo) Go libraries to generate and host our API documentation. When running this container go to [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) to view our API docs.

## Running a local test instance

You can run a local instance to test things out by running the `bootstrap-local.sh` script in this repo. The pre-requisites for running this script are:

- docker - [https://docs.docker.com/engine/install/](https://docs.docker.com/engine/install/)
- docker-compose (included with docker if you are on Windows or Mac)
- mongodb tools - [https://www.mongodb.com/docs/database-tools/](https://www.mongodb.com/docs/database-tools/installation/installation/)