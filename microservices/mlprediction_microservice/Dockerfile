# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /


COPY go.mod ./
RUN go mod download


COPY *.go ./

RUN go build -o /mlprediction

EXPOSE 9001

CMD [ "/mlprediction" ]

