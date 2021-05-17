FROM alpine

#RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN rm -rf /var/cache/apk/*

RUN mkdir /app
WORKDIR /

COPY parmenides /app/demokritos
COPY sullego /sullego

ENV TMPDIR=/tmp
ENV GOMAXPROCS=8

EXPOSE 5000
ENTRYPOINT ["/app/demokritos"]
