FROM alpine:3.13

RUN apk add --no-cache bash curl git

ENTRYPOINT ["ghput"]
CMD [ "-h" ]

COPY ghput_*.apk /tmp/
RUN apk add --allow-untrusted /tmp/ghput_*.apk
