FROM golang:1.16-alpine as build-env

#PRE-STEP set variable
ARG project_name

#STEP 1 build the image
#RUN apk --no-cache add git
ENV GO111MODULE=on
RUN mkdir /app
WORKDIR /app
COPY go.mod .
COPY go.sum .

# Get dependenciesgo
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN GO111MODULE=on GOOS=linux CGO_ENABLED=0 go build /app/${project_name}/main.go

FROM alpine

ARG project_name
ENV project_name=${project_name}
#RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN rm -rf /var/cache/apk/*

RUN mkdir /app
WORKDIR /

COPY --from=build-env /app/main /app/
RUN mv /app/main /app/${project_name}

ENV TMPDIR=/tmp
ENV GOMAXPROCS=8

EXPOSE 5000
ENTRYPOINT [ "sh", "-c", "/app/${project_name}" ]
