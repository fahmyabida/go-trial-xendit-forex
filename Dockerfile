FROM golang:alpine as build-env
ARG SERVICE_NAME=forex

RUN mkdir /builder
ADD . /builder
 
WORKDIR /builder
RUN apk add --no-cache git
RUN go build -o ${SERVICE_NAME} .
 
FROM alpine
RUN mkdir /app
WORKDIR /app
COPY --from=build-env /builder/${SERVICE_NAME}          /app/${SERVICE_NAME}

RUN apk add --no-cache tzdata
ENV TZ Asia/Jakarta

ENTRYPOINT ["/app/forex"]
