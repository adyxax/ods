##### Build ####################################################################
FROM docker.io/library/golang:1.23-alpine

WORKDIR /usr/src/app

COPY checksums /
RUN CGO_ENABLED=0 go build -o /usr/local/bin/ods ./

##### Run ######################################################################
FROM alpine:3.20

ENV LANG en_US.utf8

RUN apk upgrade --no-cache

COPY --from=0 /usr/local/bin/ods /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/ods"]
