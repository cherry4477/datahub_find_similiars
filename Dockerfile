FROM golang
 
COPY . /go/src/github.com/yaxinlx/datahub-find-similiars
RUN go install github.com/yaxinlx/datahub-find-similiars

ENV DB root:root@tcp(10.1.235.96:3306)/datahub?charset=utf8
ENV PORT 6666
ENTRYPOINT /go/bin/similiars -port=$PORT 
EXPOSE $PORT
