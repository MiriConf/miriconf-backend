ARG TARGETARCH

FROM --platform=$TARGETPLATFORM mongo:latest

COPY ./dev-deploy/dummy-data /home/data

CMD ["mongorestore", "--host=mongodb.default.svc.cluster.local.", "--port=27017", "--username=root", "--password=localdev", "/home/data"]