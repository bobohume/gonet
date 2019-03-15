FROM golang:latest

RUN chmod +x /bin/start.sh

ENV GATEWAY_LOG_LEVEL=info

EXPOSE 8081 31100 31200 31300 31700

WORKDIR /bin
ENTRYPOINT ["/bin/sh", "./start.sh"]
