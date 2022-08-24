FROM alpine:latest

COPY ./miriconf-backend .

EXPOSE 8080

CMD ["./miriconf-backend"]