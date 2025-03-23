FROM golang:1.24 AS build
LABEL org.opencontainers.image.authors="Gregory Man <man.gregory@gmail.com>"

WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /tcp-paste


FROM alpine
COPY --from=build /tcp-paste /tcp-paste
RUN apk add --update ca-certificates && rm -rf /var/cache/apk/*
RUN mkdir /data
VOLUME /data

ENV HOSTNAME=localhost:8080
ENV SLACK_CHANNEL=test

EXPOSE 8080 4343 9393

CMD ["sh", "-c", "exec /tcp-paste -storage=/data -hostname=${HOSTNAME} -slack-token=${SLACK_TOKEN} -slack-chanel=${SLACK_CHANNEL}"]
