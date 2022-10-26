ARG TARGETARCH

FROM --platform=$TARGETPLATFORM mongo:latest

COPY ./dev-deploy/dummy-data /home/data

CMD ["mongorestore", "mongodb://admin:localdev@mongo-svc.default.svc.cluster.local", "/home/data"]
