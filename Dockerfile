FROM golang:latest

RUN mkdir gonet
COPY . /gonet/

EXPOSE 8081 31100 31200 31300 31700

WORKDIR /gonet/src/server

RUN go build

CMD  ["/gonet/server netgate"]
#USER root
#FROM centos:latest
#COPY ./bin /usr/local/bin
#ENV GATEWAY_LOG_LEVEL=info
#EXPOSE 8081 31100 31200 31300 31700
#WORKDIR /usr/local/bin
#RUN chmod u+x server
#RUN chmod u+x start.sh
#ENTRYPOINT  ["/bin/sh", "./start.sh"]
#USER root
