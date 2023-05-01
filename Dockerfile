# syntax=docker/dockerfile:1

FROM golang:1.20.3-alpine3.17
LABEL description="Forum"
LABEL vendor="Created by: Olena Budarahina (Gitea username: obudarah), Kristina Volkova (Gitea username: Mustkass)"
WORKDIR /forum

RUN apk add build-base
RUN apk add sqlite
COPY go.mod .
RUN go mod download

COPY . .
RUN go build -C app -o ../forum

EXPOSE 8080

CMD [ "./forum" ]