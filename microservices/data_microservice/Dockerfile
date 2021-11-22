# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /


COPY go.mod ./
RUN go mod download


COPY *.go ./

RUN go build -o /datamicroservice

EXPOSE 9000

CMD [ "/datamicroservice" ]

