ARG TARGETARCH

FROM --platform=$TARGETPLATFORM mongo:latest

COPY ./dev-deploy/dummy-data /home/data

CMD ["mongorestore", "mongodb+srv://admin:localdev@mongodb-svc.default.svc.cluster.local/?replicaSet=mongodb&ssl=false", "--nsInclude=miriconf.teams", "/home/data"]
