FROM golang:1.19.0-bullseye



RUN apt-get install git
COPY entrypoint.sh /bin
RUN mkdir -p /go/src/github.com/yihanzhen/fetcher/cmd
ADD cmd /go/src/github.com/yihanzhen/fetcher/cmd
COPY go.mod /go/src/github.com/yihanzhen/fetcher


ENTRYPOINT [ "/bin/entrypoint.sh" ]