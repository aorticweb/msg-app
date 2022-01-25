FROM golang:latest

# ENV GO_ARCH linux-amd64

# PG Migration
RUN go get github.com/amacneil/dbmate

COPY ./containers/migration/entrypoint.sh /bin/entrypoint.sh
RUN chmod 757 /bin/entrypoint.sh

COPY /db/migrations /db/migrations

CMD ["/bin/entrypoint.sh"]