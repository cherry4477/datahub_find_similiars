FROM golang:1.5.1

ENV SERVICE_NAME datahub-find-similiars
RUN go get github.com/asiainfoLDP/$SERVICE_NAME

ENV SRC_DIR /go/src/github.com/asiainfoLDP/$SERVICE_NAME
WORKDIR $SRC_DIR
RUN go build

ENV SERVICE_PORT 9999
EXPOSE $SERVICE_PORT

ENV START_SCRIPT start.sh
RUN chmod +x $START_SCRIPT

# docker doesn't expend env in most instructions
CMD ["sh", "-c", "./$START_SCRIPT $SERVICE_NAME $SERVICE_PORT"]

