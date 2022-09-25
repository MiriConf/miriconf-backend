# miriconf-backend

Backend api for miriconf installations, it is using the [Golang Gin Web Framework](https://gin-gonic.com) to assemble the api. This repo will automatically build and test the go binary, package it in a container, and push the container to [Github Packages](https://github.com/orgs/MiriConf/packages?repo_name=miriconf-backend). For our API documentation we are using the [swaggo](https://github.com/swaggo) Go libraries to generate and host our API documentation. When running this container go to [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) to view our API docs.

## Running a local dev environment with Tilt on Windows

Clone this repository as well as the [MiriConf Frontend](https://github.com/MiriConf/miriconf-frontend) repository, make sure they are both in the same folder like below:

```
Documents
|
--> miriconf
    |
    --> miriconf-frontend
    |
    --> miriconf-backend
```

Install dependencies:

- Install scoop by running the below commands in PowerShell:

```
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

```
irm get.scoop.sh | iex
```

- Install Tilt and Helm by running the below command:

```
scoop install tilt helm
```

- Install Docker - https://www.docker.com/products/docker-desktop/
- Make sure docker is started and enable kubernetes in docker settings and run the below in PowerShell:

```
kubectl config use-context docker-desktop
```

Once dependencies are met, open a Powershell terminal in the miriconf-backend folder and run `tilt up -f tiltfile-backend` then press space to open tilt in your browser. You should now have a local instance of miriconf running at http://localhost:8080 as well as a local instance of the miriconf backend api running at http://localhost:8081.
