# add back ARG TARGETARCH

# add back --platform=$TARGETPLATFORM
FROM busybox:latest AS repo-fetch

WORKDIR /home

RUN wget https://github.com/NixOS/nixpkgs/archive/refs/heads/master.zip

RUN unzip master.zip

# add back --platform=$TARGETPLATFORM
FROM golang:latest AS go-build

RUN mkdir /home/build

COPY ./pkgs-src /home/build

WORKDIR /home/build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -v ./

# add back --platform=$TARGETPLATFORM
FROM nixos/nix:master

COPY --from=repo-fetch /home/nixpkgs-master/pkgs /home/pkgs

COPY --from=go-build /home/build/pkgs-job /home/pkgs-job

COPY ./pkgs-src/applications.bson /home/data/applications.bson

COPY ./pkgs-src/applications.metadata.json /home/data/applications.metadata.json

WORKDIR /home/pkgs

ENTRYPOINT [ "/home/pkgs-job" ]
