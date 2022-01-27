FROM golang:latest as dev

EXPOSE 8080

COPY . /msg-app
WORKDIR /msg-app

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
RUN cd app/cmd/ && CGO_ENABLED=0 GOOS=linux go build -v -o api

CMD ["air"]

FROM scratch as prod

# FROM ubuntu:bionic as prod

EXPOSE 8080

COPY --from=dev /msg-app/app/cmd/api /bin/api

CMD ["/bin/api"]
