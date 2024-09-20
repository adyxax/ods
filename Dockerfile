##### Build ####################################################################
FROM docker.io/library/golang:1.23-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build ./...

##### Run ######################################################################
FROM alpine:3.20

ENV LANG en_US.utf8

RUN apk upgrade --no-cache

COPY --from=0 /usr/src/app/ods /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/ods"]
