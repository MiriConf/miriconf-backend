ARG TARGETARCH

FROM --platform=$TARGETPLATFORM golang:latest AS go-build

RUN mkdir /home/build

COPY ./src /home/build

WORKDIR /home/build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -v ./

FROM --platform=$TARGETPLATFORM alpine:latest

COPY --from=go-build /home/build/miriconf-backend /usr/local/bin

EXPOSE 8081

CMD ["miriconf-backend"]