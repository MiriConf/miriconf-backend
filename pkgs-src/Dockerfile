ARG TARGETARCH

FROM --platform=$TARGETPLATFORM busybox:latest AS repo-fetch

WORKDIR /home

RUN wget https://github.com/NixOS/nixpkgs/archive/refs/heads/master.zip

RUN unzip master.zip

FROM --platform=$TARGETPLATFORM golang:latest AS go-build

RUN mkdir /home/build

COPY . /home/build

WORKDIR /home/build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -v ./

FROM --platform=$TARGETPLATFORM nixos/nix:master

COPY --from=repo-fetch /home/nixpkgs-master/pkgs /home/pkgs

COPY --from=go-build /home/build/pkgs-job /home/pkgs-job

COPY ./applications.bson /home/data/applications.bson

COPY ./applications.metadata.json /home/data/applications.metadata.json

WORKDIR /home/pkgs

ENTRYPOINT [ "/home/pkgs-job" ]
