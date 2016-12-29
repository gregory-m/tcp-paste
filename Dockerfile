FROM alpine:3.4
MAINTAINER Gregory Man <man.gregory@gmail.com>

COPY tcp-paste-linux-amd64 /tcp-paste
EXPOSE 8080 4343 9393

RUN mkdir /data
VOLUME /data

ENV HOSTNAME localhost:8080

CMD ["sh", "-c", "exec /tcp-paste -storage=/data -hostname=${HOSTNAME}"]
