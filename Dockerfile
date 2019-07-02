FROM golang:latest

ADD . /go/

WORKDIR /go/bin
RUN ["/bin/sh", "./build.sh"]
RUN chmod +x server
EXPOSE 8081 31200 31300 31700
#ENTRYPOINT  ["./server", "netgate"]
CMD ["./server", "account"]
CMD ["./server", "world"]
CMD ["./server", "netgate"]

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

#docker文件运行
#docker build -t dockerfile .
#删除空的镜像
#docker rmi -f $(docker images | grep "none" | awk '{print $3}')
#后台运行容器
#docker run -d -t -p31700:31700 -p31100:31100 -p31700:31700 dockerfile
#docker run -i -t dockerfile