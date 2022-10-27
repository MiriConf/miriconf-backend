# miriconf-backend

This is the backend api for miriconf installations, it is using the [Gorilla Mux](https://github.com/gorilla/mux) package to manage routing for the api. This repo will automatically build and test the go binary, package it in a container, and push the container to [Github Packages](https://github.com/orgs/MiriConf/packages?repo_name=miriconf-backend). For our API documentation we are using the [swaggo](https://github.com/swaggo) Go libraries to generate and host our API documentation. When running this container go to [http://localhost:8080/docs/](http://localhost:8080/docs/) to view our API docs.

The tiltfile in this repo is using a hardcoded JWT token secret