FROM alpine:3.4

RUN addgroup -S tile38 && adduser -S -G tile38 tile38

RUN apk update \
    && apk add ca-certificates \
    && update-ca-certificates \
    && apk add openssl 

ADD tile38-server /usr/local/bin 
ADD tile38-cli /usr/local/bin 

RUN mkdir /data && chown tile38:tile38 /data

VOLUME /data
WORKDIR /data

EXPOSE 9851
CMD ["tile38-server"]
