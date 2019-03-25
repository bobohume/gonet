FROM golang:latest

ADD . /go/

WORKDIR /go/bin
RUN ["/bin/sh", "./build.sh"]
RUN chmod +x server
EXPOSE 8081 31100 31200 31300 31700
#ENTRYPOINT  ["/bin/sh", "./server"]
#ENTRYPOINT  ["/bin/sh", "./server"]
#ENTRYPOINT  ["/bin/sh", "./start.sh"]
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
